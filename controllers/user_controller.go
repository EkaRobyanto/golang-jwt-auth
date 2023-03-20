package controllers

import (
	"golang-auth/helper"
	"golang-auth/initializers"
	"golang-auth/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	///get Email
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	_, isEmailValidated := helper.ValidateEmail(body.Email)

	if !isEmailValidated {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Masukan email dengan benar",
		})
		return
	}

	if body.Email == "" {
		c.JSON(400, gin.H{"error": "email tidak boleh kosong"})
		return
	}

	if body.Password == "" {
		c.JSON(400, gin.H{"error": "password tidak boleh kosong"})
		return
	}

	//hash pw
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to Hash pw",
		})
		return
	}

	//create user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Where(models.User{Email: body.Email}).Attrs(user).FirstOrCreate(&user)
	if result.RowsAffected == 1 {
		c.JSON(201, gin.H{
			"message":  "User Created",
			"email":    body.Email,
			"password": body.Password,
		})
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	//resp
	c.JSON(400, gin.H{
		"error": "Email telah terdaftar",
	})
}

func Login(c *gin.Context) {
	//get email and password
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	//look up requested user
	user := models.User{}
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(400, gin.H{
			"error": "Email belum terdaftar",
		})
		return
	}

	//compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Password salah",
		})
		return
	}
	//generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		// log.Fatal(err)
		c.JSON(400, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	//send it back

	//cookie
	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	//token
	c.JSON(200, gin.H{
		"token": tokenString,
	})
}

func GetProfile(c *gin.Context) {
	data, _ := c.Get("user")
	c.JSON(200, gin.H{
		"data": data,
	})
}
