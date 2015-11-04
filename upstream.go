/*
The upstream type implements http.Handler.

In this file we handle request routing and interaction with the authBackend.
*/

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

type upstream struct {
	httpClient  *http.Client
	authBackend string
}

type gitService struct {
	method         string
	regex          *regexp.Regexp
	middlewareFunc func(u *upstream, w http.ResponseWriter, r *http.Request, handleFunc func(w http.ResponseWriter, r *gitRequest, rpc string), rpc string)
	handleFunc     func(w http.ResponseWriter, r *gitRequest, rpc string)
	rpc            string
}

// A gitReqest is an *http.Request decorated with attributes returned by the
// GitLab Rails application.
type gitRequest struct {
	*http.Request
	// GL_ID is an environment variable used by gitlab-shell hooks during 'git
	// push' and 'git pull'
	GL_ID string
	// RepoPath is the full path on disk to the Git repository the request is
	// about
	RepoPath string
	// ArchivePath is the full path where we should find/create a cached copy
	// of a requested archive
	ArchivePath string
	// ArchivePrefix is used to put extracted archive contents in a
	// subdirectory
	ArchivePrefix string
	// CommitId is used do prevent race conditions between the 'time of check'
	// in the GitLab Rails app and the 'time of use' in gitlab-workhorse.
	CommitId string
}

// Routing table
var gitServices = [...]gitService{
	gitService{"GET", regexp.MustCompile(`/info/refs\z`), repoPreAuth, handleGetInfoRefs, ""},
	gitService{"POST", regexp.MustCompile(`/git-upload-pack\z`), repoPreAuth, handlePostRPC, "git-upload-pack"},
	gitService{"POST", regexp.MustCompile(`/git-receive-pack\z`), repoPreAuth, handlePostRPC, "git-receive-pack"},
	gitService{"GET", regexp.MustCompile(`/repository/archive\z`), repoPreAuth, handleGetArchive, "tar.gz"},
	gitService{"GET", regexp.MustCompile(`/repository/archive.zip\z`), repoPreAuth, handleGetArchive, "zip"},
	gitService{"GET", regexp.MustCompile(`/repository/archive.tar\z`), repoPreAuth, handleGetArchive, "tar"},
	gitService{"GET", regexp.MustCompile(`/repository/archive.tar.gz\z`), repoPreAuth, handleGetArchive, "tar.gz"},
	gitService{"GET", regexp.MustCompile(`/repository/archive.tar.bz2\z`), repoPreAuth, handleGetArchive, "tar.bz2"},
}

func newUpstream(authBackend string, authTransport http.RoundTripper) *upstream {
	return &upstream{&http.Client{Transport: authTransport}, authBackend}
}

func (u *upstream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var g gitService

	log.Printf("%s %q", r.Method, r.URL)

	// Look for a matching Git service
	foundService := false
	for _, g = range gitServices {
		if r.Method == g.method && g.regex.MatchString(r.URL.Path) {
			foundService = true
			break
		}
	}
	if !foundService {
		// The protocol spec in git/Documentation/technical/http-protocol.txt
		// says we must return 403 if no matching service is found.
		http.Error(w, "Forbidden", 403)
		return
	}

	g.middlewareFunc(u, w, r, g.handleFunc, g.rpc)
}

func repoPreAuth(u *upstream, w http.ResponseWriter, r *http.Request, handleFunc func(w http.ResponseWriter, r *gitRequest, rpc string), rpc string) {
	url := u.authBackend + r.URL.RequestURI()
	authReq, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		fail500(w, "doAuthRequest", err)
		return
	}
	// Forward all headers from our client to the auth backend. This includes
	// HTTP Basic authentication credentials (the 'Authorization' header).
	for k, v := range r.Header {
		authReq.Header[k] = v
	}
	// Also forward the Host header, which is excluded from the Header map by the http libary.
	// This allows the Host header received by the backend to be consistent with other
	// requests not going through gitlab-workhorse.
	authReq.Host = r.Host
	// Set a custom header for the request. This can be used in some
	// configurations (Passenger) to solve auth request routing problems.
	authReq.Header.Set("GitLab-Git-HTTP-Server", Version)

	authResponse, err := u.httpClient.Do(authReq)

	if err != nil {
		fail500(w, "doAuthRequest", err)
		return
	}
	defer authResponse.Body.Close()

	if authResponse.StatusCode != 200 {
		// The Git request is not allowed by the backend. Maybe the
		// client needs to send HTTP Basic credentials.  Forward the
		// response from the auth backend to our client. This includes
		// the 'WWW-Authenticate' header that acts as a hint that
		// Basic auth credentials are needed.
		for k, v := range authResponse.Header {
			// Accomodate broken clients that do case-sensitive header lookup
			if k == "Www-Authenticate" {
				w.Header()["WWW-Authenticate"] = v
			} else {
				w.Header()[k] = v
			}
		}
		w.WriteHeader(authResponse.StatusCode)
		io.Copy(w, authResponse.Body)
		return
	}

	// The auth backend validated the client request and told us additional
	// request metadata. We must extract this information from the auth
	// response body.
	gitReq := &gitRequest{Request: r}
	if err := json.NewDecoder(authResponse.Body).Decode(gitReq); err != nil {
		fail500(w, "decode JSON GL_ID", err)
		return
	}
	// Don't hog a TCP connection in CLOSE_WAIT, we can already close it now
	authResponse.Body.Close()

	// Negotiate authentication (Kerberos) may need to return a WWW-Authenticate
	// header to the client even in case of success as per RFC4559.
	for k, v := range authResponse.Header {
		// Case-insensitive comparison as per RFC7230
		if strings.EqualFold(k, "WWW-Authenticate") {
			w.Header()[k] = v
		}
	}

	if !looksLikeRepo(gitReq.RepoPath) {
		http.Error(w, "Not Found", 404)
		return
	}

	handleFunc(w, gitReq, rpc)
}

func looksLikeRepo(p string) bool {
	// If /path/to/foo.git/objects exists then let's assume it is a valid Git
	// repository.
	if _, err := os.Stat(path.Join(p, "objects")); err != nil {
		log.Print(err)
		return false
	}
	return true
}