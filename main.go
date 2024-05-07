package main

import (
	"ETicaret/Database"
	"ETicaret/Router"
)

func main() {
	database.Connect()
	app := Router.Routes()

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
