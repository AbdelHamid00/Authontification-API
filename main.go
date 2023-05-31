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
	"github.com/golang-jwt/jwt"
	"time"
	"github.com/joho/godotenv"
)

var (
	DB *gorm.DB
	FD *os.File
)

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
	var client struct {
		login		string
		password	string
	}
	err := ctx.BindJSON(&client)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Error Reading the body !",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(client.password), 10)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Error Hashing the password",
		})
		return
	}
	user := usermodel.Usermodel{Login: client.login, Password: string(hash), Admin: false}
	result := DB.Create(&user)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": "Error Creating the user",
		})
		return
	}
	ctx.JSON(200, nil)
}


func LoginHandler(ctx *gin.Context) {
	var client struct {
		login		string
		password	string
	}
	err := ctx.BindJSON(&client)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Eroor Parsing the Body",
		})
		return 
	}
	var user usermodel.Usermodel
	DB.First(&user, "email = ?", client.login)
	if user.ID == 0 {
		ctx.JSON(400, gin.H{
			"error": "Invalid Login or Password",
		})
		return 
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(client.password))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid Login or Password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	_, err = token.SignedString(os.Getenv("SECRET"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Failed to create the token",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"token": token,
	})
}
