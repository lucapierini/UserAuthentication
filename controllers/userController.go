package controllers

import(
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lucapierini/UserAuthentication/models"
	helpers "github.com/lucapierini/UserAuthentication/helpers"
	"golang.org/x/crypto/bcrypt"
	
	var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

	var validate = validator.New()

	func HashPassword(password string) string {
		bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			log.Panic(err)
		}
		return string(bytes)
	}

	func VerifyPassword(hashedPassword, password string) error {
		return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	}

	func SingUp(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validate := validator.New()
		err := validate.Struct(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user.Password = HashPassword(user.Password)
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		helpers.DB.Create(&user)
		c.JSON(http.StatusOK, gin.H{"data": user})
	}

	func Login(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var foundUser models.User
		helpers.DB.Where("email = ?", user.Email).First(&foundUser)
		if foundUser.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		err := VerifyPassword(foundUser.Password, user.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": foundUser})
	}

	func GetUser(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		helpers.DB.Where("id = ?", id).First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			}
			c.JSON(http.StatusOK, gin.H{"data": user})
	}

	func GetUsers(c *gin.Context) {
		var users []models.User
		helpers.DB.Find(&users)
		c.JSON(http.StatusOK, gin.H{"data": users})
	}
)