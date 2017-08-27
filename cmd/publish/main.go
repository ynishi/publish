package main

import (
	"os"
	"strings"

	"github.com/urfave/cli"
	"github.com/ynishi/publish"
)

func main() {
	app := cli.NewApp()
	app.Name = "publish"
	app.Version = "0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Yutaka Nishimura",
			Email: "ytk.nishimura@gmail.com",
		},
	}
	app.Copyright = "(c) 2017 Yutaka Nishimura"
	app.Usage = "publish doc"
	app.Action = func(c *cli.Context) error {
		publish.SetReader(strings.NewReader("test"))
		publishers := []publish.Publisher{
			&publish.PublishGitHub{},
		}
		publish.Publish(publishers)
		return nil
	}

	app.Run(os.Args)

}
