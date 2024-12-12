# UserAuthentication

Para correr el proyecto directamente en go se levanta un contenedor de mongo para la base de datos y se ejecuta luego el go run main.go desde la raiz.
docker run -d --name mongodb-container-auth -p 27017:27017 -v mongo-data:/data/db mongo:latest

Para correr el proyecto con docker-compose se ejecuta el comando docker-compose up desde la raiz del proyecto.