package app

import (
	"github.com/labstack/echo"
	"github.com/iqOptionTest/simplecacheclient/cacheclient"
)

type App struct {
	appName     string
	port        string
	cacheClient *cacheclient.Client
	echo        *echo.Echo
}

func New(appName string, port string, cacheEndpoint string) *App {
	a := &App{
		appName:     appName,
		port:        port,
		cacheClient: cacheclient.NewClient(cacheEndpoint),
		echo:        echo.New(),
	}

	return a
}

func (a *App) Run() error {
	a.registerHTTPHandlers()
	a.runHTTPServer(a.port)
	return nil
}

func (a *App) registerHTTPHandlers() {
	a.echo.Static("/public", "public")
	a.echo.File("/", "public/index.html")
	a.echo.POST("/set", a.Set)
	a.echo.GET("/get/:key", a.Get)
	a.echo.GET("/keys", a.Keys)
	a.echo.DELETE("/unset/:key", a.Unset)
	a.echo.POST("/rpush", a.Rpush)
	a.echo.GET("/pop/:key", a.Pop)
	a.echo.GET("/lget/:key/:i", a.Lget)
	a.echo.GET("/lgetall/:key", a.Lgetall)
	a.echo.POST("/hset", a.Hset)
	a.echo.GET("/hget/:key/:i", a.Hget)
	a.echo.GET("/hgetall/:key", a.Hgetall)
	a.echo.GET("/ping", a.healthCheck)
}

func (a *App) runHTTPServer(port string) {
	a.echo.Logger.Fatal(a.echo.Start(port))
}
