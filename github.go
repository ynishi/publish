package publish

import (
	"errors"
	"io"

	"context"
	"fmt"
	"net/url"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type PublishGitHub struct {
	Publisher
	Conf *viper.Viper
}

type PublishGitHubOpts struct {
	Owner    string
	Repo     string
	Token    string
	Branch   string
	Endpoint string
}

func (pgh *PublishGitHub) Publish(r io.Reader) error {
	var pgho PublishGitHubOpts
	pgh.Conf.ReadInConfig()
	pgh.Conf.Unmarshal(&pgho)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pgho.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	u, err := url.Parse(pgho.Endpoint)
	if err != nil {
		return err
	}
	client.BaseURL = u
	service := client.Git

	pRef, _, err := service.GetRef(pgho.Owner, pgho.Repo, fmt.Sprintf("heads/%s", pgho.Branch))
	if err != nil {
		return err
	}
	if pRef.Object == nil {
		return errors.New("error: cannot fetch parent ref.")
	}
	pTree, _, err := client.Git.GetTree(pgho.Owner, pgho.Repo, *pRef.Object.SHA, true)
	if err != nil {
		return err
	}
	if pTree.Entries == nil {
		return errors.New("error: cannot fetch parent tree.")
	}

	parent, _, err := client.Git.GetCommit(pgho.Owner, pgho.Repo, *pRef.Object.SHA)
	if err != nil {
		return err
	}
	if parent == nil {
		return errors.New("error: cannot fetch parent commit.")
	}
	input := &github.Commit{
		Message: github.String("m"),
		Tree:    &github.Tree{SHA: pTree.SHA},
		Parents: []github.Commit{{SHA: parent.SHA}},
	}
	commit, _, err := client.Git.CreateCommit(pgho.Owner, pgho.Repo, input)
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
	service.UpdateRef(pgho.Owner, pgho.Repo, nRef, false)
	return nil
}
