package services

import (
	"fmt"
	"github.com/kataras/iris"
)

func LoadResources(app *iris.Application) {
	app.RegisterView(iris.HTML("./templates", ".html").Reload(true))
	app.HandleDir("/public", "./static")
}

func GetLoginPage(app *iris.Application) {
	app.Get("/login", func(ctx iris.Context) {
		ctx.View("login.html")
	})
}

type User struct {
	Username string
	Password string
}

func GetInfoPage(app *iris.Application) {
	app.Post("/info", func(ctx iris.Context) {
		form := User{}
		err := ctx.ReadForm(&form)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.WriteString(err.Error())
		}
		fmt.Println(form)
		username := form.Username
		password := form.Password
		ctx.ViewData("username", username)
		ctx.ViewData("password", password)
		ctx.View("info.html")
	})
}

func NotImplement(app *iris.Application) {
	app.Get("/unknown", func(ctx iris.Context) {
		ctx.StatusCode(501)
		ctx.JSON(iris.Map{
			"error": "501 not implement error",
		})
	})
}

func GetStaticPage(app *iris.Application) {
	app.Get("/public", func(ctx iris.Context) {
		ctx.HTML(`<a href='/public/css/main.css'>/public/css/main.css</a><br/><br/>
			<a href='/public/img/test.jpg'>/public/img/test.jpg</a><br/><br/>
			<a href='/public/js/showStatic.js'>/public/js/showStatic.js</a>`)
	})
}
