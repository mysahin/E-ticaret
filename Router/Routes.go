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
	r.Post("/add-product", Controllers.Product{}.AddProduct)                           //ekle
	r.Get("/view-my-products/", Controllers.Product{}.ViewMyProduct)                   //kendi ürünlerini görüntüle
	r.Get("/view-product/:id", Controllers.Product{}.ViewProductById)                  //id ye göre tek ürün görüntüleme
	r.Get("/view-by-type/:type", Controllers.Product{}.ViewProductsByType)             //tipe göre görüntüle
	r.Get("/view-by-category/:category", Controllers.Product{}.ViewProductsByCategory) //kategoriye göre görüntüle
	r.Post("/delete-product/:id", Controllers.Product{}.DeleteProduct)                 //sil
	r.Put("/archive/:id", Controllers.Product{}.ArchiveProduct)                        //arşivle veya arşivden çıkar
	r.Put("/edit-product", Controllers.Product{}.EditProduct)                          //düzenle
	r.Put("/rate-product/:productID/:rating", Controllers.Product{}.RateProduct)
	//Sepete ürün ekleme çıkarma işlemleri
	r.Get("/view-cart", Controllers.ViewCart)                              //Sepettekileri görüntüle
	r.Post("/add-cart/:productID", Controllers.AddToCart)                  //Sepete ekle
	r.Delete("/remove-cart/:productID", Controllers.RemoveFromCart)        //Sepetten sil
	r.Put("/decrease-cart/:productID", Controllers.DecreaseQuantityInCart) //Sepete eklenen ürünün sayısını düşür
	//Diğer işlemler
	r.Get("/search/", Controllers.Search)

	return r
}
