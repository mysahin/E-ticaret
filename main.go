package main

import (
	"ETicaret/Database"
	"ETicaret/Router"
)

func main() {
	database.Connect()
	app := Router.Routes()

	err := app.Listen(":8000")
	if err != nil {
		panic(err)
	}
}
