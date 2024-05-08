package Router

import (
	"ETicaret/Controllers"
	database "ETicaret/Database"
	"github.com/gofiber/fiber/v2"
)

func Routes() *fiber.App {
	r := fiber.New()
	rdb := database.ConnectRedis()
	//Login işlemleri
	r.Post("/signup", Controllers.Login{}.SignUp)
	r.Post("/signin", Controllers.Login{}.SignIn)
	r.Get("/signout", Controllers.Login{}.SignOut)
	//Ürün işlemleri
	r.Post("/add-product", Controllers.Product{}.AddProduct)                           //ekle
	r.Get("/view-my-products", Controllers.Product{}.ViewMyProduct)                    //kendi ürünlerini görüntüle
	r.Get("/view-product/:id", Controllers.Product{}.ViewProductById)                  //id ye göre tek ürün görüntüleme
	r.Get("/view-by-type/:type", Controllers.Product{}.ViewProductsByType)             //tipe göre görüntüle
	r.Get("/view-by-category/:category", Controllers.Product{}.ViewProductsByCategory) //kategoriye göre görüntüle
	r.Post("/delete-product/:id", Controllers.Product{}.DeleteProduct)                 //sil
	r.Put("/archive/:id", Controllers.Product{}.ArchiveProduct)                        //arşivle veya arşivden çıkar
	r.Put("/edit-product", Controllers.Product{}.EditProduct)                          //düzenle
	//Sepete ürün ekleme çıkarma işlemleri
	r.Get("/view-cart", func(ctx *fiber.Ctx) error {
		return Controllers.ViewCart(ctx, rdb)
	})
	r.Post("/add-cart/:productID", func(ctx *fiber.Ctx) error {
		return Controllers.AddToCart(ctx, rdb)
	})
	r.Delete("/remove-cart/:productID", func(ctx *fiber.Ctx) error {
		return Controllers.RemoveFromCart(ctx, rdb)
	})
	r.Put("/decrease-cart/:productID", func(ctx *fiber.Ctx) error {
		return Controllers.DecreaseQuantityInCart(ctx, rdb)
	})
	return r
}
