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
	r.Post("/add-product", Handlers.AddProduct)                              //ekle
	r.Get("/view-my-products/", Handlers.ViewMyProduct)                      //kendi ürünlerini görüntüle
	r.Get("/view-product/:id", Handlers.ViewProductById)                     //id ye göre tek ürün görüntüleme
	r.Get("/view-by-type/:type", Handlers.ViewProductsByType)                //tipe göre görüntüle
	r.Get("/view-by-category/:category", Handlers.ViewProductsByCategory)    //kategoriye göre görüntüle
	r.Get("/view-by-seller/:seller", Handlers.ViewProductBySeller)           //Satıcının tüm ürünlerini görme
	r.Post("/delete-product/:id", Handlers.DeleteProduct)                    //sil
	r.Put("/archive/:id", Handlers.ArchiveProduct)                           //arşivle veya arşivden çıkar
	r.Put("/edit-product", Handlers.EditProduct)                             //düzenle
	r.Put("/rate-product/:productID/:rating", Handlers.RateProduct)          //puan ver
	r.Put("/comment-product/:productID/", Handlers.CommentProduct)           //yorum yap
	r.Get("/view-product-comments/:productId", Handlers.ViewProductComments) //yorumlara bak

	//Sepete ürün ekleme çıkarma işlemleri
	r.Get("/view-cart", Handlers.ViewCart)                              //Sepettekileri görüntüle
	r.Post("/add-cart/:productID", Handlers.AddToCart)                  //Sepete ekle
	r.Delete("/remove-cart/:productID", Handlers.RemoveFromCart)        //Sepetten sil
	r.Put("/decrease-cart/:productID", Handlers.DecreaseQuantityInCart) //Sepete eklenen ürünün sayısını düşür

	//Diğer işlemler
	r.Get("/home-page/", Handlers.HomePage)                      //arama ve ana sayfa
	r.Get("/search-by-categories", Handlers.SearchPageCategorie) //Arama sonuçları filtreleme
	r.Get("/search-by-type", Handlers.SearchPageType)            //Arama sonucu filtreleme (tip)
	return r
}
