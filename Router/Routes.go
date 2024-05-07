package Router

import (
	"ETicaret/Controllers"
	"github.com/gofiber/fiber/v2"
)

func Routes() *fiber.App {
	r := fiber.New()
	//Login işlemleri
	r.Post("/signup", Controllers.Login{}.SignUp)
	r.Post("/signin", Controllers.Login{}.SignIn)
	r.Get("/signout", Controllers.Login{}.SignOut)
	//Ürün işlemleri
	r.Post("/add-product", Controllers.Product{}.AddProduct)
	r.Get("/view-my-products", Controllers.Product{}.ViewMyProduct)
	return r
}
