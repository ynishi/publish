package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Msg string
}

func greet(name string) string {
	return fmt.Sprintf("hello, %s!", name)
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	p := viper.GetStringMap("person")

	fmt.Printf("%s\n", greet(p["name"].(string)))
}
