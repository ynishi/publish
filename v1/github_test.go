// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package publish

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"strings"

	"context"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
)

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method(%q): %v, want %v", r, got, want)
	}
}

type createCommit struct {
	Author    *github.CommitAuthor `json:"author,omitempty"`
	Committer *github.CommitAuthor `json:"committer,omitempty"`
	Message   *string              `json:"message,omitempty"`
	Tree      *string              `json:"tree,omitempty"`
	Parents   []string             `json:"parents,omitempty"`
}
type createTree struct {
	BaseTree string             `json:"base_tree,omitempty"`
	Entries  []github.TreeEntry `json:"tree"`
}
type updateRefRequest struct {
	SHA   *string `json:"sha"`
	Force *bool   `json:"force"`
}

const (
	content         = `# test`
	contentBase64   = `IyB0ZXN0`
	filename        = `doc.md`
	contentEncoding = `utf-8`
	message         = `Change doc.md(by publishGitHub)`
	testTomlTmpl    = `[GitHub]
Owner = "o"
Repo = "r"
Token = "t"
Branch = "b"
Endpoint = "%s"
Encoding = "%s"
Path = "%s"`
)

var (
	mux    *http.ServeMux
	server *httptest.Server
)

func init() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

func TestInitConfGitHub(t *testing.T) {

	testEndpoint := "endpoint"

	testGitHubOpts := &PublishGitHubOpts{
		Owner:    "o",
		Repo:     "r",
		Token:    "t",
		Branch:   "b",
		Endpoint: testEndpoint,
		Encoding: contentEncoding,
		Path:     filename,
	}

	github_toml := fmt.Sprintf(testTomlTmpl, testEndpoint, contentEncoding, filename)

	publishGitHub := &PublishGitHub{}

	c := viper.New()
	c.SetConfigType("toml")
	c.ReadConfig(strings.NewReader(github_toml))

	InitConfGitHub(publishGitHub, c)

	if !reflect.DeepEqual(publishGitHub.GitHub, testGitHubOpts) {
		t.Fatalf("GitHub not matchted\n want: %q,\n have: %q", testGitHubOpts, publishGitHub.GitHub)
	}
}

func TestPublishGitHub(t *testing.T) {

	mux.HandleFunc("/repos/o/r/git/commits/s", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"sha":"s","message":"m","author":{"name":"n"}}`)
	})

	mux.HandleFunc("/repos/o/r/git/refs/heads/b", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			testMethod(t, r, "GET")
			fmt.Fprint(w, `
		  {
		    "ref": "refs/heads/b",
		    "url": "https://api.github.com/repos/o/r/git/refs/heads/b",
		    "object": {
		      "type": "commit",
		      "sha": "s",
		      "url": "https://api.github.com/repos/o/r/git/commits/s"
		    }
		  }`)
		case "PATCH":
			args := &updateRefRequest{
				SHA:   github.String("s2"),
				Force: github.Bool(false),
			}
			v := new(updateRefRequest)
			json.NewDecoder(r.Body).Decode(v)

			testMethod(t, r, "PATCH")
			if !reflect.DeepEqual(v, args) {
				t.Errorf("Request body = %+v, want %+v", v, args)
			}
			fmt.Fprint(w, `
		  {
		    "ref": "refs/heads/b",
		    "url": "https://api.github.com/repos/o/r/git/refs/heads/b",
		    "object": {
		      "type": "commit",
		      "sha": "s2",
		      "url": "https://api.github.com/repos/o/r/git/commits/s2"
		    }
		  }`)
		}
	})

	mux.HandleFunc("/repos/o/r/git/blobs", func(w http.ResponseWriter, r *http.Request) {

		v := new(github.Blob)

		json.NewDecoder(r.Body).Decode(v)

		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}

		want := &github.Blob{
			Content:  github.String(content),
			Encoding: github.String(contentEncoding),
			Size:     github.Int(len(content)),
		}

		if !reflect.DeepEqual(v, want) {
			t.Errorf("Git.CreateBlob request body: %+v, want %+v", v, want)
		}

		fmt.Fprintf(w, `{
		 "sha": "sn",
		 "content": "%s",
		 "encoding": "%s",
		 "size": %d
		}`, content, contentEncoding, len(content))
	})

	mux.HandleFunc("/repos/o/r/git/trees/s", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{
			  "sha": "s",
			  "tree": [ { "type": "blob" } ]
			}`)
	})

	mux.HandleFunc("/repos/o/r/git/trees", func(w http.ResponseWriter, r *http.Request) {
		v := new(createTree)
		json.NewDecoder(r.Body).Decode(v)

		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}

		input := []github.TreeEntry{
			{
				Path: github.String(filename),
				Mode: github.String("100644"),
				Type: github.String("blob"),
				SHA:  github.String("sn"),
			},
		}

		want := &createTree{
			BaseTree: "s",
			Entries:  input,
		}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Git.CreateTree request body: %+v, want %+v", v, want)
		}

		fmt.Fprintf(w, `{
		  "sha": "st2",
		  "tree": [
		    {
		      "path": "%s",
		      "mode": "100644",
		      "type": "blob",
		      "size": %d,
		      "sha": "s2"
		    }
		  ]
		}`, filename, len(content))
	})

	mux.HandleFunc("/repos/o/r/git/commits", func(w http.ResponseWriter, r *http.Request) {
		v := new(createCommit)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")

		want := &createCommit{
			Message: github.String(message),
			Tree:    github.String("st2"),
			Parents: []string{"s"},
		}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"sha":"s2"}`)
	})

	github_toml := fmt.Sprintf(testTomlTmpl, server.URL, contentEncoding, filename)

	publishGitHub := &PublishGitHub{}

	c := viper.New()
	c.SetConfigType("toml")
	c.ReadConfig(strings.NewReader(github_toml))

	InitConfGitHub(publishGitHub, c)
	r := strings.NewReader(content)

	errChan := make(chan error, 1)
	ctx := context.Background()

	go func() {
		errChan <- publishGitHub.Publish(ctx, r)
	}()

	select {
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	case err := <-errChan:
		if err != nil {
			t.Fatal(err)
		}
	}

	ctxc, cancel := context.WithCancel(ctx)
	defer cancel()
	go publishGitHub.Publish(ctxc, r)
	cancel()
	select {
	case <-ctxc.Done():
		// do nothing
	default:
		t.Fatal("failed cancel: %q", ctx)
	}
}
