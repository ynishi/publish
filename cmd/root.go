// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

/*
   The publish/cmd package implement command line interface for
   publish. It use config file for each publisher.
 */
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ynishi/publish"
)

// RootCmd is root command.
var RootCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish document to web services.",
	Long: `Publish is a document publisher for multible web services in Golang.

This application is a tool for a working document to set on web services.
Project is available at http://github.com/ynishi/publish`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("start root command")

		r, err := os.Open(content)
		if err != nil {
			fmt.Println(err)
		}

		publish.SetReader(r)
		publish.SetTimeout(time.Duration(timeout) * time.Second)

		pas3 := &publish.PublishAwsS3{}
		pgh := &publish.PublishGitHub{}

		err = publish.InitConfAwsS3(pas3, aS3Conf)
		if err != nil {
			fmt.Println(err)
		}
		err = publish.InitConfGitHub(pgh, ghConf)
		if err != nil {
			fmt.Println(err)
		}
		publishers := []publish.Publisher{
			pgh,
			pas3,
		}
		err = publish.Publish(publishers)
		if err != nil {
			fmt.Println(err)
		}
		log.Print("end root command")
	},
}

var (
	cfgFile, content string
	timeout          int
	ghConf, aS3Conf  *viper.Viper
)

func init() {
	log.SetPrefix("publish")
	log.SetFlags(log.LstdFlags)

	cobra.OnInitialize(func() {
		setupConfPath(viper.GetViper())

		ghConf = viper.New()
		setupConfPath(ghConf)

		aS3Conf = viper.New()
		setupConfPath(aS3Conf)
	})

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "config.toml", "config file (default in . or $HOME/.publish or /etc/publish)")
	RootCmd.PersistentFlags().StringVarP(&content, "content", "c", "content.md", "content to publish( default is ./content.md")
	RootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 300, "timeout each publish in seconds(default is 300 Seconds)")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("content", RootCmd.PersistentFlags().Lookup("content"))
	viper.BindPFlag("timeout", RootCmd.PersistentFlags().Lookup("timeout"))
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
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
