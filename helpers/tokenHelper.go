package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lucapierini/UserAuthentication/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Definición de la estructura SignedDetails que incluye los detalles del usuario y los claims estándar de JWT
type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.StandardClaims
}

// Inicialización de la colección de usuarios desde la base de datos
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// Obtención de la clave secreta para firmar los tokens desde las variables de entorno
var SECRET_KEY string = os.Getenv("SECRET_KEY")

// Función para generar tokens de acceso y de refresco
func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {
	// Crea los claims del token con los detalles del usuario y la fecha de expiración
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	// Crea los claims del token de refresco con una fecha de expiración extendida
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	// Genera el token firmado usando los claims y una clave secreta
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

    // Genera el token de refresco firmado usando los claims de refresco y la misma clave secreta
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return
	}
    // Retorna el token firmado, el token de refresco firmado y cualquier error que haya ocurrido
	return token, refreshToken, err
}

// Función para validar un token firmado y devolver los claims del usuario
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	// 
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprint("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg

}

func UpdateAllTokens(singedToken string, signedRefreshToken string, userId string) {
	// Crear un contexto con un tiempo de espera de 100 segundos
    var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
    defer cancel() // Asegurarse de cancelar el contexto al final de la función

    // Crear un objeto de actualización para almacenar los nuevos valores de los tokens
    var updateObj primitive.D

    // Agregar el token firmado al objeto de actualización
    updateObj = append(updateObj, bson.E{Key: "token", Value: singedToken})
    // Agregar el token de refresco firmado al objeto de actualización
    updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

    // Obtener la fecha y hora actual en formato RFC3339
    Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    // Agregar la fecha de actualización al objeto de actualización
    updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

    // Configurar la opción de upsert para que inserte un nuevo documento si no se encuentra uno existente
    upsert := true
    // Crear un filtro para buscar el documento con el user_id especificado
    filter := bson.M{"user_id": userId}
    // Configurar las opciones de actualización
    opt := options.UpdateOptions{
        Upsert: &upsert,
    }

    // Actualizar un documento en la colección userCollection
    _, err := userCollection.UpdateOne(ctx, filter, bson.D{
        {Key: "$set", Value: updateObj},
    }, &opt)

    // Verificar si ocurrió un error durante la actualización
    if err != nil {
        log.Panic(err) // Registrar un mensaje de pánico si ocurre un error
        return
    }
}
