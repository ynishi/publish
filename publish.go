// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package publish

import (
	"io"
)

type Publisher interface {
	Publish(io.Reader) error
}

var reader io.Reader

func Publish(publishers []Publisher) error {

	for _, publisher := range publishers {
		err := publisher.Publish(reader)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetReader(r io.Reader) {
	reader = r
}
