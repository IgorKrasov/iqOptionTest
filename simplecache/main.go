package main

import (
	"github.com/iqOptionTest/simplecache/app"
)

func main() {
	a := app.NewApp()
	a.Initialize()
	a.Run(":9003")
}
