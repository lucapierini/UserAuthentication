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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/primitive"
	"github.com/lucapierini/UserAuthentication/database"
)
	
	var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

	var validate = validator.New()

	func HashPassword(password string) string {
		
	}

	func VerifyPassword(userPassword string, providedPassword string)(bool, string) {
		err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
		check := true
		msg := " "

		if err != nil {
			msg = "email or password are incorrect"
			check = false
		}
	}

	func Signup()gin.HandlerFunc{
		return func(c *gin.Context) {
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var user models.User
			
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

			count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
			defer cancel()
			if err != nil{
				log.Panic(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking user phone"})
			}

			if count > 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone already in use"})
				return
			}
			
			user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.ID = primitive.NewObjectID()
			user.User_id = user.ID.Hex()
			token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, *user.First_name, *user.Last_name, *user.User_type, *user.User_id)
			user.Token = &token
			user.RefreshToken = &refreshToken
			
			resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
			if insertErr != nil {
				msg := fmt.Sprintf("User item was not created")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"data": resultInsertionNumber})
	}
}

	func Login() gin.HandlerFunc {
		return func(c *gin.Context){
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
			var user models.User
			var foundUser models.User

			if err := c.BindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
			defer cancel()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid credentials"})
				return
			}

			passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
			defer cancel()
		}
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
