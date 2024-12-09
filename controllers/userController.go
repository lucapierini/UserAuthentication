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
		
	}

	func VerifyPassword(hashedPassword, password string) error {
		
	}

	func Signup()gin.HandlerFunc{
		return func(c *gin.Context) {
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var user models. User
			
			if err := c.BindJSON(&user); err != nil {
				c.JSON(http. StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			
			validationErr := validate. Struct(user)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
			defer cancel()
			if err != nil{
				log.Panic(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking user email"})
			}

			count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
			defer cancel()
			if err != nil{
				log.Panic(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking user phone"})
			}

			if count > 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone already in use"})
				return
			}
		}

	func Login(c *gin.Context) {
		
	}

	func GetUser() gin.HandlerFunc {
		return func(c *gin.Context) {
			userId := c.Param("user_id")

			if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			var ctx, cancel = context.WhthTimeout(context.Background(), 100*time.Second)

			var user models.User
			err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		
			defer cancel()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": user})
		
		}
	}

	func GetUsers(c *gin.Context) {
		
	}
)