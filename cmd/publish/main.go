package main

import (
	"os"

	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"github.com/ynishi/publish"
)

var conf, ghConf *viper.Viper

func init() {
	// init conf
	conf = viper.New()
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conf.AddConfigPath(home)
	conf.AddConfigPath(".")
	conf.SetConfigName(".publish")
	conf.SetConfigType("toml")
	err = conf.ReadInConfig()
	if err != nil {
		fmt.Println("failed read .publish: ", err)
		os.Exit(1)
	}
	conf.SetEnvPrefix("publish")
	conf.BindEnv("content")
	viper.BindPFlag("content", pflag.CommandLine.Lookup("content"))
	conf.SetDefault("content", "content.md")

	ghConf = viper.New()

	ghConf.AddConfigPath(home)
	ghConf.AddConfigPath(".")
	ghConf.SetConfigName(".publish")
	ghConf.SetConfigType("toml")
	err = ghConf.ReadInConfig()
	if err != nil {
		fmt.Println("failed read .publish: ", err)
		os.Exit(1)
	}
	ghConf.SetEnvPrefix("publish_gh")
}

func main() {

	content := conf.GetString("content")

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
		r, err := os.Open(content)
		if err != nil {
			fmt.Println(err)
		}
		publish.SetReader(r)
		publishers := []publish.Publisher{
			&publish.PublishGitHub{
				Conf: ghConf,
			},
		}
		publish.Publish(publishers)
		return nil
	}

	app.Run(os.Args)

}
