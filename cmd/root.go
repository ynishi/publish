// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ynishi/publish"
)

var RootCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish document to web services.",
	Long: `Publish is a document publisher for multible web services in Golang.

This application is a tool for a working document to set on web services.
Project is available at http://github.com/ynishi/publish`,
	Run: func(cmd *cobra.Command, args []string) {
		r, err := os.Open(content)
		if err != nil {
			fmt.Println(err)
		}
		publish.SetReader(r)
		publishers := []publish.Publisher{
			&publish.PublishGitHub{
				Conf: ghConf,
			},
			&publish.PublishAwsS3{
				Conf: aS3Conf,
			},
		}
		err = publish.Publish(publishers)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var (
	cfgFile, content string
	ghConf, aS3Conf  *viper.Viper
)

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "config.toml", "config file (default in . or $HOME/.publish or /etc/publish)")
	RootCmd.PersistentFlags().StringVarP(&content, "content", "c", "content.md", "content to publish( default is ./content.md")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("content", RootCmd.PersistentFlags().Lookup("content"))

	ghConf = viper.New()
	setupConfPath(ghConf)

	aS3Conf = viper.New()
	setupConfPath(aS3Conf)
}

func initConfig() {
	setupConfPath(viper.GetViper())
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func setupConfPath(v *viper.Viper) {
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		v.AddConfigPath(path.Join(home, ".publish"))
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/publish")
		v.SetConfigName("config")
		v.SetConfigType("toml")
	}
}
