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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/lucapierini/UserAuthentication/database"
)
	
	var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

	var validate = validator.New()

	func HashPassword(password string) string {
		bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			log.Panic(err)
		}
		return string(bytes)
	}

	func VerifyPassword(userPassword string, providedPassword string)(bool, string) {
		err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
		check := true
		msg := " "

		if err != nil {
			msg = "email or password are incorrect"
			check = false
		}
		return check, msg
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

			password := HashPassword(*user.Password)
			user.Password = &password

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
			
			user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.ID = primitive.NewObjectID()
			user.UserId = user.ID.Hex()
			token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *&user.UserId)
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
			if passwordIsValid != true {
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				return
			}

			if foundUser.Email  == nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
				return
			}
			token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, *&foundUser.UserId)
			helpers.UpdateAllTokens(token, refreshToken, foundUser.UserId)
			userCollection.FindOne(ctx,bson.M{"token": token}).Decode(&foundUser)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while updating user token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": foundUser})
		}
	}

	func GetUser() gin.HandlerFunc {
		return func(c *gin.Context) {
			userId := c.Param("user_id")

			if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

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

	func GetUsers() gin.HandlerFunc {
		return func(c *gin.Context) {
			if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

			recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
			if err != nil || recordPerPage < 1 {
				recordPerPage = 10
			}

			page, err1 := strconv.Atoi(c.Query("page"))
			if err1 != nil || page < 1 {
				page = 1
			}

			startIndex := (page - 1) * recordPerPage
			startIndex, err = strconv.Atoi(c.Query("startIndex"))

			matchStage := bson.D{{"$match", bson.D{{}}}}
			
			groupStage := bson.D{
				{"$group", bson.D{
					{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
					{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
					{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
					}},
				}

			projectStage := bson.D{
				{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count", Value: 1}, 
					{Key: "user_item", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
		}},
	}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	var allUsers []bson.M
	if err = result.All(ctx, &allUsers); err != nil {
		log.Fatal(err)

	}
	c.JSON(http.StatusOK, gin.H{"data": allUsers[0]})
}
}
