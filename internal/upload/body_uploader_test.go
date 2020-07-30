package upload

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/filestore"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/testhelper"
)

const (
	fileContent = "A test file content"
	fileLen     = len(fileContent)
)

func TestBodyUploader(t *testing.T) {
	testhelper.ConfigureSecret()

	body := strings.NewReader(fileContent)

	resp := testUpload(&rails{}, nil, echoProxy(t, fileLen), body)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	uploadEcho, err := ioutil.ReadAll(resp.Body)

	require.NoError(t, err, "Can't read response body")
	require.Equal(t, fileContent, string(uploadEcho))
}

func TestBodyUploaderCustomPreparer(t *testing.T) {
	body := strings.NewReader(fileContent)

	resp := testUpload(&rails{}, &alwaysLocalPreparer{}, echoProxy(t, fileLen), body)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	uploadEcho, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err, "Can't read response body")
	require.Equal(t, fileContent, string(uploadEcho))
}

func TestBodyUploaderCustomVerifier(t *testing.T) {
	body := strings.NewReader(fileContent)
	verifier := &mockVerifier{}

	resp := testUpload(&rails{}, &alwaysLocalPreparer{verifier: verifier}, echoProxy(t, fileLen), body)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	uploadEcho, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err, "Can't read response body")
	require.Equal(t, fileContent, string(uploadEcho))
	require.True(t, verifier.invoked, "Verifier.Verify not invoked")
}

func TestBodyUploaderAuthorizationFailure(t *testing.T) {
	testNoProxyInvocation(t, http.StatusUnauthorized, &rails{unauthorized: true}, nil)
}

func TestBodyUploaderErrors(t *testing.T) {
	tests := []struct {
		name     string
		preparer *alwaysLocalPreparer
	}{
		{name: "Prepare failure", preparer: &alwaysLocalPreparer{prepareError: fmt.Errorf("")}},
		{name: "Verify failure", preparer: &alwaysLocalPreparer{verifier: &alwaysFailsVerifier{}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testNoProxyInvocation(t, http.StatusInternalServerError, &rails{}, test.preparer)
		})
	}
}

func testNoProxyInvocation(t *testing.T, expectedStatus int, auth PreAuthorizer, preparer Preparer) {
	proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "request proxied upstream")
	})

	resp := testUpload(auth, preparer, proxy, nil)
	require.Equal(t, expectedStatus, resp.StatusCode)
}

func testUpload(auth PreAuthorizer, preparer Preparer, proxy http.Handler, body io.Reader) *http.Response {
	req := httptest.NewRequest("POST", "http://example.com/upload", body)
	w := httptest.NewRecorder()

	BodyUploader(auth, proxy, preparer).ServeHTTP(w, req)

	return w.Result()
}

func echoProxy(t *testing.T, expectedBodyLength int) http.Handler {
	require := require.New(t)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		require.NoError(err)

		require.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"), "Wrong Content-Type header")

		require.Contains(r.PostForm, "file.md5")
		require.Contains(r.PostForm, "file.sha1")
		require.Contains(r.PostForm, "file.sha256")
		require.Contains(r.PostForm, "file.sha512")

		require.Contains(r.PostForm, "file.path")
		require.Contains(r.PostForm, "file.size")
		require.Equal(strconv.Itoa(expectedBodyLength), r.PostFormValue("file.size"))

		token, err := jwt.ParseWithClaims(r.Header.Get(RewrittenFieldsHeader), &MultipartClaims{}, testhelper.ParseJWT)
		require.NoError(err, "Wrong JWT header")

		rewrittenFields := token.Claims.(*MultipartClaims).RewrittenFields
		if len(rewrittenFields) != 1 || len(rewrittenFields["file"]) == 0 {
			t.Fatalf("Unexpected rewritten_fields value: %v", rewrittenFields)
		}

		path := r.PostFormValue("file.path")
		uploaded, err := os.Open(path)
		require.NoError(err, "File not uploaded")

		//sending back the file for testing purpose
		io.Copy(w, uploaded)
	})
}

type rails struct {
	unauthorized bool
}

func (r *rails) PreAuthorizeHandler(next api.HandleFunc, _ string) http.Handler {
	if r.unauthorized {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r, &api.Response{TempPath: os.TempDir()})
	})
}

type alwaysLocalPreparer struct {
	verifier     Verifier
	prepareError error
}

func (a *alwaysLocalPreparer) Prepare(_ *api.Response) (*filestore.SaveFileOpts, Verifier, error) {
	return filestore.GetOpts(&api.Response{TempPath: os.TempDir()}), a.verifier, a.prepareError
}

type alwaysFailsVerifier struct{}

func (alwaysFailsVerifier) Verify(handler *filestore.FileHandler) error {
	return fmt.Errorf("Verification failed")
}

type mockVerifier struct {
	invoked bool
}

func (m *mockVerifier) Verify(handler *filestore.FileHandler) error {
	m.invoked = true

	return nil
}
