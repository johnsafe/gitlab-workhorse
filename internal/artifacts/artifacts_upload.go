package artifacts

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/filestore"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/helper"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/upload"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/zipartifacts"
)

var zipSubcommandsErrorsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gitlab_workhorse_zip_subcommand_errors_total",
		Help: "Errors comming from subcommands used for processing ZIP archives",
	}, []string{"error"})

func init() {
	prometheus.MustRegister(zipSubcommandsErrorsCounter)
}

type artifactsUploadProcessor struct {
	opts *filestore.SaveFileOpts

	upload.SavedFileTracker
}

func (a *artifactsUploadProcessor) generateMetadataFromZip(ctx context.Context, file *filestore.FileHandler) (*filestore.FileHandler, error) {
	metaReader, metaWriter := io.Pipe()
	defer metaWriter.Close()

	metaOpts := &filestore.SaveFileOpts{
		LocalTempPath:  a.opts.LocalTempPath,
		TempFilePrefix: "metadata.gz",
	}
	if metaOpts.LocalTempPath == "" {
		metaOpts.LocalTempPath = os.TempDir()
	}

	fileName := file.LocalPath
	if fileName == "" {
		fileName = file.RemoteURL
	}

	zipMd := exec.CommandContext(ctx, "gitlab-zip-metadata", fileName)
	zipMd.Stderr = os.Stderr
	zipMd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	zipMd.Stdout = metaWriter

	if err := zipMd.Start(); err != nil {
		return nil, err
	}
	defer helper.CleanUpProcessGroup(zipMd)

	type saveResult struct {
		error
		*filestore.FileHandler
	}
	done := make(chan saveResult)
	go func() {
		var result saveResult
		result.FileHandler, result.error = filestore.SaveFileFromReader(ctx, metaReader, -1, metaOpts)

		done <- result
	}()

	if err := zipMd.Wait(); err != nil {
		st, ok := helper.ExitStatus(err)

		if !ok {
			return nil, err
		}

		zipSubcommandsErrorsCounter.WithLabelValues(zipartifacts.ErrorLabelByCode(st)).Inc()

		if st == zipartifacts.CodeNotZip {
			return nil, nil
		}

		if st == zipartifacts.CodeLimitsReached {
			return nil, zipartifacts.ErrBadMetadata
		}
	}

	metaWriter.Close()
	result := <-done
	return result.FileHandler, result.error
}

func (a *artifactsUploadProcessor) ProcessFile(ctx context.Context, formName string, file *filestore.FileHandler, writer *multipart.Writer) error {
	//  ProcessFile for artifacts requires file form-data field name to eq `file`

	if formName != "file" {
		return fmt.Errorf("invalid form field: %q", formName)
	}
	if a.Count() > 0 {
		return fmt.Errorf("artifacts request contains more than one file")
	}
	a.Track(formName, file.LocalPath)

	select {
	case <-ctx.Done():
		return fmt.Errorf("ProcessFile: context done")

	default:
		// TODO: can we rely on disk for shipping metadata? Not if we split workhorse and rails in 2 different PODs
		metadata, err := a.generateMetadataFromZip(ctx, file)
		if err != nil {
			return err
		}

		if metadata != nil {
			for k, v := range metadata.GitLabFinalizeFields("metadata") {
				writer.WriteField(k, v)
			}

			a.Track("metadata", metadata.LocalPath)
		}
	}

	return nil
}

func (a *artifactsUploadProcessor) Name() string {
	return "artifacts"
}

func UploadArtifacts(myAPI *api.API, h http.Handler) http.Handler {
	return myAPI.PreAuthorizeHandler(func(w http.ResponseWriter, r *http.Request, a *api.Response) {
		mg := &artifactsUploadProcessor{opts: filestore.GetOpts(a), SavedFileTracker: upload.SavedFileTracker{Request: r}}

		upload.HandleFileUploads(w, r, h, a, mg)
	}, "/authorize")
}
