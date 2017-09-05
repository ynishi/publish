// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

/*
   The publish package implements a simple library and clients for
   document publish tool to web services.

   Publish do Publish method each publishers. Support context and
   concurrent execution. Count goroutines and return when all jobs
   has done.
 */
package publish

import (
	"context"
	"io"
	"log"
	"os"
	"time"
)

// Publisher is able to Publish document.
type Publisher interface {
	Publish(context.Context, io.Reader) error
}

var (
	reader  io.Reader
	timeout time.Duration
	logger  *log.Logger
)

func init() {
	SetLogger(log.New(os.Stdout, "publish: ", log.LstdFlags|log.Lshortfile))
}

// Publish do Publish method of each publishers.
func Publish(publishers []Publisher) error {

	n := len(publishers)
	logger.Printf("all publishers: %d", n)
	errc := make(chan error, 1)
	ctx := context.Background()

	for _, publisher := range publishers {
		go func(publisher Publisher) {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			ctx = context.WithValue(ctx, "name", publisher)
			defer cancel()

			go func() {
				errc <- publisher.Publish(ctx, reader)
			}()

			select {
			case <-ctx.Done():
				logger.Printf("%s, %s", ctx.Value("name"), ctx.Err())
				return
			}
		}(publisher)
	}

	for {
		err := <-errc
		if err != nil {
			n--
			logger.Printf("todo: %d, err: %q\n", n, err)
		} else {
			n--
			logger.Printf("todo: %d, job done", n)
		}
		if n < 1 {
			logger.Println("all publishers done")
			return nil
		}
	}
}

// SetReader set reader
func SetReader(r io.Reader) {
	reader = r
}

// SetTimeout set timeout
func SetTimeout(t time.Duration) {
	timeout = t
}

// SetLogger set logger
func SetLogger(l *log.Logger) {
	logger = l
}
