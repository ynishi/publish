package publish

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

var githubHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{}`)
})

func TestPublishGitHub(t *testing.T) {
	ts := httptest.NewServer(githubHandler)
	publishGitHub := &PublishGitHub{
		Conf: viper.New(),
	}
	publishGitHub.Conf.SetDefault("Owner", "o")
	publishGitHub.Conf.SetDefault("Repo", "r")
	publishGitHub.Conf.SetDefault("Token", "t")
	publishGitHub.Conf.SetDefault("Branch", "b")
	publishGitHub.Conf.Set("endpoint", ts.URL)
	err := publishGitHub.Publish(reader)
	if err != nil {
		t.Fatal(err)
	}
}
