package Router

import (
	"ETicaret/Handlers"
	"github.com/gofiber/fiber/v2"
)

func Routes() *fiber.App {
	r := fiber.New()

	//Login işlemleri
	r.Post("/signup", Handlers.Login{}.SignUp)
	r.Post("/signin", Handlers.Login{}.SignIn)
	r.Get("/signout", Handlers.Login{}.SignOut)
	//Ürün işlemleri
	r.Post("/add-product", Handlers.Product{}.AddProduct)                           //ekle
	r.Get("/view-my-products/", Handlers.Product{}.ViewMyProduct)                   //kendi ürünlerini görüntüle
	r.Get("/view-product/:id", Handlers.Product{}.ViewProductById)                  //id ye göre tek ürün görüntüleme
	r.Get("/view-by-type/:type", Handlers.Product{}.ViewProductsByType)             //tipe göre görüntüle
	r.Get("/view-by-category/:category", Handlers.Product{}.ViewProductsByCategory) //kategoriye göre görüntüle
	r.Post("/delete-product/:id", Handlers.Product{}.DeleteProduct)                 //sil
	r.Put("/archive/:id", Handlers.Product{}.ArchiveProduct)                        //arşivle veya arşivden çıkar
	r.Put("/edit-product", Handlers.Product{}.EditProduct)                          //düzenle
	r.Put("/rate-product/:productID/:rating", Handlers.Product{}.RateProduct)       //puan ver
	r.Put("/rate-product/:productID/:comment", Handlers.Product{}.CommentProduct)   //yorum yap
	//Sepete ürün ekleme çıkarma işlemleri
	r.Get("/view-cart", Handlers.ViewCart)                              //Sepettekileri görüntüle
	r.Post("/add-cart/:productID", Handlers.AddToCart)                  //Sepete ekle
	r.Delete("/remove-cart/:productID", Handlers.RemoveFromCart)        //Sepetten sil
	r.Put("/decrease-cart/:productID", Handlers.DecreaseQuantityInCart) //Sepete eklenen ürünün sayısını düşür
	//Diğer işlemler
	r.Get("/search/", Handlers.Search)

	return r
}
