// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a Apache License 2.0
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

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
		}
		publish.Publish(publishers)
	},
}

var (
	cfgFile, content string
	ghConf           *viper.Viper
)

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", ".publish.toml", "config file (default is . or $HOME/.publish.toml)")
	RootCmd.PersistentFlags().StringVarP(&content, "content", "c", "content.md", "content to publish( default is ./content.md")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("content", RootCmd.PersistentFlags().Lookup("content"))

	initGitHubConfig()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".publish")
		viper.SetConfigType("toml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func initGitHubConfig() {
	ghConf = viper.New()

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ghConf.AddConfigPath(home)
	ghConf.AddConfigPath(".")
	ghConf.SetConfigName(".publishGitHub")
	ghConf.SetConfigType("toml")

	if err := ghConf.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
