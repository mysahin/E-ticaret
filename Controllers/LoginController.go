package Controllers

import (
	database "ETicaret/Database"
	"ETicaret/Helpers"
	"ETicaret/Models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type Login struct {
}

var SecretKey = "secret"

func (login Login) SignUp(c *fiber.Ctx) error {
	db := database.DB.Db
	var user = new(Models.User)

	if err := c.BodyParser(&user); err != nil {
		return err
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "İşlem sırasında bir hata oluştu.",
		})
	}

	newUser := Models.User{
		Name:     user.Name,
		LastName: user.LastName,
		UserName: user.UserName,
		Password: string(password[:]),
		Adress:   user.Adress,
		TcNo:     user.TcNo,
	}

	newUserLogin := Models.Login{
		UserName: user.UserName,
		Password: string(password[:]),
	}
	if err := db.Create(&newUserLogin).Error; err != nil {
		return err
	}
	if err := db.Create(&newUser).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"NewUser": newUser,
	})
}
func (login Login) SignIn(c *fiber.Ctx) error {
	db := database.DB.Db
	var user = new(Models.Login)
	var compareUser = new(Models.Login)
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	if err := db.First(&compareUser, "user_name=?", user.UserName).Error; err != nil {
		return c.JSON(fiber.Map{
			"error": "",
		})
	}
	if compareUser.ID == 0 {
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(compareUser.Password), []byte(user.Password)); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(compareUser.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"token": token,
	})
}
func (login Login) SignOut(c *fiber.Ctx) error {
	isLogin := Helpers.IsLogin(c)
	if isLogin {
		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(time.Hour * 1),
			HTTPOnly: true,
		}

		c.Cookie(&cookie)
		return c.JSON(fiber.Map{
			"message": "Başarıyla çıkış yapıldı",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Zaten giriş yapılmamış!",
	})
}
