package main

import (
	"os"
	"github.com/github-user/cloudgo-io/services"
	"github.com/kataras/iris"
	"github.com/spf13/pflag"
)

const (
	PORT string = "8080"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}
	pPort := pflag.StringP("port", "p", PORT, "http listening port")
	pflag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}
	app := iris.Default()
	app.Logger().SetLevel("debug")
	services.StartServices(app)
	app.Run(iris.Addr(":"+port), iris.WithConfiguration(iris.TOML("./configs/main.tml")))
}
