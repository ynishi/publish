// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package publish

import (
	"errors"
	"io"
	"io/ioutil"

	"context"
	"fmt"
	"net/url"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type PublishGitHub struct {
	Publisher
	GitHub *PublishGitHubOpts
}

type PublishGitHubOpts struct {
	Owner    string
	Repo     string
	Token    string
	Branch   string
	Endpoint string
	Encoding string
	Path     string
}

func (gh *PublishGitHub) String() string {
	return "PublishGitHub"
}

func InitConfGitHub(gh *PublishGitHub, c *viper.Viper) (err error) {
	if c == nil {
		return errors.New("error: conf is nil. pointer to viper is needed.")
	}
	c.SetDefault("Encoding", "utf-8")
	err = c.UnmarshalKey("GitHub", &gh.GitHub)
	if err != nil {
		return err
	}
	return nil
}

func (pgh *PublishGitHub) Publish(ctx context.Context, r io.Reader) error {
	logger.Println("start publish github")

	pgho := pgh.GitHub

	if pgho.Owner == "" || pgho.Repo == "" || pgho.Token == "" || pgho.Branch == "" || pgho.Path == "" {
		return errors.New("error: cannot fetch conf vars.")
	}
	// make client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pgho.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	if pgho.Endpoint != "" {
		e := ""
		if pgho.Endpoint[len(pgho.Endpoint)-1:] != "/" {
			e = "/"
		}
		u, err := url.Parse(pgho.Endpoint + e)
		if err != nil {
			return err
		}
		client.BaseURL = u
		client.UploadURL = u
	}
	service := client.Git

	// prepare
	pRef, _, err := service.GetRef(context.Background(), pgho.Owner, pgho.Repo, fmt.Sprintf("heads/%s", pgho.Branch))
	if err != nil {
		return err
	}
	if pRef.Object == nil {
		return errors.New("error: cannot fetch parent ref.")
	}
	pTree, _, err := client.Git.GetTree(context.Background(), pgho.Owner, pgho.Repo, *pRef.Object.SHA, true)
	if err != nil {
		return err
	}
	if pTree.Entries == nil {
		return errors.New("error: cannot fetch parent tree.")
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	str := string(buf)
	blob := &github.Blob{
		Content:  github.String(str),
		Encoding: github.String(pgho.Encoding),
		Size:     github.Int(len(str)),
	}
	b, _, err := service.CreateBlob(context.Background(), pgho.Owner, pgho.Repo, blob)
	if err != nil {
		return err
	}
	entry := github.TreeEntry{
		Path: github.String(pgho.Path),
		Mode: github.String("100644"),
		Type: github.String("blob"),
		SHA:  b.SHA,
	}
	entries := []github.TreeEntry{entry}
	tree, _, err := service.CreateTree(context.Background(), pgho.Owner, pgho.Repo, *pRef.Object.SHA, entries)

	// commit
	parent, _, err := client.Git.GetCommit(context.Background(), pgho.Owner, pgho.Repo, *pRef.Object.SHA)
	if err != nil {

		return err
	}
	if parent == nil {

		return errors.New("error: cannot fetch parent commit.")
	}
	input := &github.Commit{
		Message: github.String(fmt.Sprintf("Change %s(by publishGitHub)", pgho.Path)),
		Tree:    tree,
		Parents: []github.Commit{{SHA: parent.SHA}},
	}
	commit, _, err := client.Git.CreateCommit(context.Background(), pgho.Owner, pgho.Repo, input)
	if err != nil {
		return err
	}
	nRef := &github.Reference{
		Ref: github.String("refs/heads/" + pgho.Branch),
		Object: &github.GitObject{
			Type: github.String("commit"),
			SHA:  commit.SHA,
		},
	}
	service.UpdateRef(context.Background(), pgho.Owner, pgho.Repo, nRef, false)
	logger.Println("end publish github")
	return nil
}
