package publish

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

type MockPublisher struct {
	Publisher
	Conf *viper.Viper
}

func (m *MockPublisher) Publish(r io.Reader) error {
	return nil
}

var mockPublishers []Publisher
var mockPublisher *MockPublisher

func init() {
	SetReader(strings.NewReader(`<html></html>`))
	mockPublisher = &MockPublisher{
		Conf: viper.New(),
	}
	mockPublisher.Conf.Set("apikey", "test")
	mockPublishers = []Publisher{mockPublisher}
}

func TestPublisher(t *testing.T) {
	conf := viper.New()
	conf.Set("apikey", "test")
	if !reflect.DeepEqual(mockPublisher.Conf, conf) {
		t.Fatalf("Failed match reader.\n want: %q,\n have: %q\n", conf, mockPublisher.Conf)
	}
	err := mockPublisher.Publish(reader)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPublish(t *testing.T) {
	err := Publish(mockPublishers)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetReader(t *testing.T) {
	r := strings.NewReader(`test`)
	SetReader(r)
	if !reflect.DeepEqual(reader, r) {
		t.Fatalf("Failed match reader.\n want: %q,\n have: %q\n", r, reader)
	}
}

var githubHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{}`)
})

func TestPublishGitHub(t *testing.T) {
	ts := httptest.NewServer(githubHandler)
	publishGitHub := &PublishGitHub{
		Conf: viper.New(),
	}
	publishGitHub.Conf.Set("endpoint", ts.URL)
	err := publishGitHub.Publish(reader)
	if err != nil {
		t.Fatal(err)
	}
}
