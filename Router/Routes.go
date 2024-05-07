package Router

import (
	"ETicaret/Controllers"
	"github.com/gofiber/fiber/v2"
)

func Routes() *fiber.App {
	r := fiber.New()
	r.Post("/signup", Controllers.Login{}.SignUp)
	r.Post("/signin", Controllers.Login{}.SignIn)
	r.Get("/signout", Controllers.Login{}.SignOut)
	return r
}
