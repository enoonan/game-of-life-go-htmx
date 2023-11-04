package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	density   = 0.3 // Density of alive cells
	dimension = 12  // Dimension of the board
	host      = "localhost:1323"
)

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e := echo.New()
	e.Debug = true

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}
	e.Use(middleware.CORS())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"lion-shining-rarely.ngrok-free.app/"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Renderer = t
	e.Static("/static", "assets")

	e.GET("/", home)
	e.GET("/games", listGames)
	e.GET("/games/:id", show)
	e.POST("/games", start)
	e.POST("/games/:id/next", next)
	e.GET("/games/:id/runner", getRunner)
	e.GET("/games/:id/cancel-runner", cancelRunner)
	e.Logger.Debug(e.Start(":1323"))
}

func pageData() PageData {
	return PageData{Host: host}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
