// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package publish

import (
	"context"
	"io"
	"time"
)

type Publisher interface {
	Publish(context.Context, io.Reader) error
}

var (
	reader  io.Reader
	timeout time.Duration
)

func Publish(publishers []Publisher) error {

	errChan := make(chan error, 1)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for _, publisher := range publishers {
		go func() {
			errChan <- publisher.Publish(ctx, reader)
		}()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
	return nil
}

func SetReader(r io.Reader) {
	reader = r
}

func SetTimeout(t time.Duration) {
	timeout = t
}
