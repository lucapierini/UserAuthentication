package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
    // Cargar las variables de entorno desde el archivo .env
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Leer la URL de conexión de MongoDB desde las variables de entorno
    MongoDb := os.Getenv("MONGODB_URL")

    // Crear un contexto con un tiempo de espera de 10 segundos para la conexión
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel() // Asegurarse de cancelar el contexto al final de la función

    // Conectar el cliente de MongoDB usando el contexto y la URL de conexión
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDb))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")
    return client
}

// Crear una variable global para el cliente de MongoDB
var Client *mongo.Client = DBinstance()

// Función para abrir una colección específica en la base de datos
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    // Obtener la colección especificada del cliente de MongoDB
    var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
    return collection
}