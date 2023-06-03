package main

import (
	"API/LoginAPI/initializer"
	"API/LoginAPI/middlewares"
	"API/LoginAPI/usermodel"
	"os"
	"io"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"github.com/joho/godotenv"
	"fmt"
	"net/http"
)

var (
	DB *gorm.DB
	FD *os.File
)
// type Clientmodel struct {
// 	login string `json: "login" binding:"required"`
// 	password string `json: "password" binding:"required"`
// }
func SetupLogOutput() {
	FD, _ = os.Create("LoginAPI.log")
	gin.DefaultWriter = io.MultiWriter(FD, os.Stdout)
}

func CloseDefer() {
	FD.Close()
	// DB.Close()
}

func main() {
    var err error
	err = godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	DB, err = initializer.ConnectDataBase()
	if err != nil {
		panic("Error intializing the connection to the Database")
	}
	SetupLogOutput()
	defer CloseDefer()

	r := gin.New()
	r.Use(gin.Recovery(), middlewares.Logger())
	r.POST("/Login", LoginHandler)
	r.POST("/Signup", SignupHandler)

	r.Run(":8080")
}


func SignupHandler(ctx *gin.Context) {
	var client usermodel.Usermodel
	err := ctx.BindJSON(&client)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Error Reading the body !",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(client.Password), 14)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Error Hashing the password",
		})
		return
	}
	user := usermodel.Usermodel{Login: client.Login, Password: string(hash), Admin: false}
	result := DB.Create(&user)
	if result.Error != nil {
		ctx.JSON(400, gin.H {
			"error": "Error Creating the user",
		})
		return
	}
	ctx.JSON(200, nil)
}


func LoginHandler(ctx *gin.Context) {
	var client usermodel.Usermodel
	err := ctx.BindJSON(&client)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Eroor Parsing the Body",
		})
		return 
	}
	var user usermodel.Usermodel
	DB.First(&user, "Login = ?", client.Login)
	if user.ID == 0 {
		ctx.JSON(400, gin.H{
			"error": "Invalid Login",
		})
		return 
	}
	fmt.Println(user.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(client.Password))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(400, gin.H {
			"error": "Invalid Password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, erno := token.SignedString([]byte(os.Getenv("SECRET")))
	if erno != nil {
		fmt.Println(err)
		ctx.JSON(400, gin.H{
			"error": "Failed to create the token",
		})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	// To change this after we switch to https  true,false
	ctx.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "", "", false, true)
	ctx.JSON(200, gin.H{})
}
