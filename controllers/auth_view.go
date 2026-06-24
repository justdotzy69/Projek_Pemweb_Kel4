package controllers

import (
	"Projek_Pemweb_Kel4/database"
	"Projek_Pemweb_Kel4/middlewares"
	"Projek_Pemweb_Kel4/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginPage menampilkan halaman form login.
// Route: GET /login
func LoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"Title": "Login",
		"Page":  "login",
	}, "layouts/base")
}

// LoginSubmit memproses form login dari halaman web.
// Bedanya dengan Login API: kalau berhasil, token disimpan ke cookie (bukan return JSON).
// Route: POST /login
func LoginSubmit(c *fiber.Ctx) error {
	email    := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Render("auth/login", fiber.Map{
			"Title": "Login",
			"Page":  "login",
			"Error": "Email atau password salah",
		}, "layouts/base")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Render("auth/login", fiber.Map{
			"Title": "Login",
			"Page":  "login",
			"Error": "Email atau password salah",
		}, "layouts/base")
	}

	// Buat JWT dan simpan ke cookie
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Redirect("/login")
	}

	middlewares.SetTokenCookie(c, tokenString)
	return c.Redirect("/dashboard")
}

// RegisterPage menampilkan halaman form registrasi.
// Route: GET /register
func RegisterPage(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"Title": "Daftar Akun",
		"Page":  "register",
	}, "layouts/base")
}

// RegisterSubmit memproses form registrasi dari halaman web.
// Route: POST /register
func RegisterSubmit(c *fiber.Ctx) error {
	email    := c.FormValue("email")
	password := c.FormValue("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return c.Redirect("/register")
	}

	user := models.User{Email: email, Password: string(hash), Role: "user"}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Redirect("/register")
	}

	return c.Redirect("/login")
}

// Logout menghapus cookie token dan mengarahkan ke halaman login.
// Route: GET /logout
func Logout(c *fiber.Ctx) error {
	middlewares.ClearTokenCookie(c)
	return c.Redirect("/login")
}
