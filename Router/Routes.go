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
	r.Get("/view-by-type/:type", Controllers.Product{}.ViewProductsByType)
	r.Get("/view-by-category/:category", Controllers.Product{}.ViewProductsByCategory)
	r.Post("/delete-product/:id", Controllers.Product{}.DeleteProduct)
	r.Put("/archive/:id", Controllers.Product{}.ArchiveProduct)
	r.Put("/edit-product", Controllers.Product{}.EditProduct)
	return r
}
