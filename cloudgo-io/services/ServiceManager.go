package services

import (
	"github.com/kataras/iris"
)


func StartServices(app *iris.Application) {
	LoadResources(app)
	GetLoginPage(app)
	GetInfoPage(app)
	NotImplement(app)
	GetStaticPage(app)
}
