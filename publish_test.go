// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package publish

import (
	"io"
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
