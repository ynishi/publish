// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package publish

import (
	"io"
	"log"
	"reflect"
	"strings"
	"testing"

	"context"

	"time"

	"os"

	"github.com/spf13/viper"
)

type MockPublisher struct {
	Publisher
	Conf *viper.Viper
}

func (m *MockPublisher) Publish(ctx context.Context, r io.Reader) error {
	logger.Println("mock job start")
	logger.Println("mock job done")
	return nil
}

var mockPublishers []Publisher
var mockPublisher *MockPublisher

func init() {
	SetReader(strings.NewReader(`<html></html>`))
	SetTimeout(3 * time.Second)
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
	errChan := make(chan error, 1)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		errChan <- mockPublisher.Publish(ctx, reader)
	}()
	select {
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	case err := <-errChan:
		if err != nil {
			t.Fatal(err)
		}
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

func TestSetTimeout(t *testing.T) {
	t5 := 5 * time.Second
	SetTimeout(t5)
	if timeout != t5 {
		t.Fatalf("Failed match timeout.\n want: %q,\n have: %q\n", t5, timeout)
	}
}

func TestSetLogger(t *testing.T) {
	tlogger := log.New(os.Stdout, "testSetLog", log.Lshortfile)
	SetLogger(tlogger)
	if !reflect.DeepEqual(logger, tlogger) {
		t.Fatalf("Failed match logger.\n want: %q,\n have: %q\n", tlogger, logger)
	}
}
