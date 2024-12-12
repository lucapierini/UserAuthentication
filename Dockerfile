# Usa una imagen base de Go
FROM golang:1.23.1

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el archivo go.mod y go.sum
COPY go.mod go.sum ./

# Descarga las dependencias
RUN go mod download

# Copia el código fuente de la aplicación
COPY . .

# Compila la aplicación
RUN go build -o main .

# Expone el puerto en el que la aplicación se ejecutará
EXPOSE 9000

# Comando para ejecutar la aplicación
CMD ["./main"]