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
