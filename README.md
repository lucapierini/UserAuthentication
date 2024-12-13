

# Proyect-Go-JWT

Luca Pierini

---

## Descripción

**Proyect-Go-JWT** es un sistema de autenticación que permite la creación de usuarios y la validación de credenciales utilizando JSON Web Tokens (JWT) y una base de datos MongoDB. Está diseñado para ser fácilmente desplegado en entornos locales o en producción gracias a la integración con Docker y Docker Compose.

## Funcionalidades Principales

- **Registro de Usuarios**: Permite registrar nuevos usuarios almacenando sus credenciales de manera segura mediante hashing de contraseñas.
- **Inicio de Sesión**: Valida las credenciales del usuario y genera un token de acceso y un token de refresco.
- **Protección de Endpoints**: Middleware para proteger rutas mediante la validación de JWT.
- **Renovación de Tokens**: Proporciona tokens actualizados para mantener sesiones seguras.

## Endpoints

### Usuarios
- **POST /signup**: Registro de nuevos usuarios.
- **POST /login**: Inicio de sesión y generación de tokens.
- **GET /protected**: Endpoint protegido para probar la autenticación con JWT.

## Tecnologías Utilizadas

- **Backend**: Api en Go implementando Gin Gonic
- **Base de Datos**: MongoDB
- **Contenedores**: Docker y Docker Compose
- **Autenticación**: JWT

## Estructura del Proyecto

- **controllers**: Controladores de usuario.
- **models**: Definición de modelos para la base de datos.
- **helpers**: Funciones de utilidad para la generación y validación de tokens.
- **routes**: Rutas principales del sistema.
- **database**: Configuración para la conexión con MongoDB.
- **middlewares**: Middleware para proteger endpoints.

## Requisitos Técnicos

1. **Docker**: Asegúrate de tener instalado Docker y Docker Compose.
2. **Variables de Entorno**:
   - **PORT**: Puerto para la API.
   - **MONGODB_URL**: Cadena de conexión a la base de datos MongoDB.

## Instrucciones de Instalación

1. Clona este repositorio:

   ```bash
   git clone https://github.com/lucapierini/proyect-go-jwt.git
   cd proyect-go-jwt
   ```

2. Crea un archivo `.env` en la raíz del proyecto con el siguiente contenido:

   ```env
   PORT=9000
   MONGODB_URL=mongodb://localhost:27017/go-auth
   ```

3. Levanta los contenedores utilizando Docker Compose:

   ```bash
   docker-compose up
   ```

4. La API estará disponible en `http://localhost:9000`.

## Uso de la API

1. **Registro de Usuario**:
   - Endpoint: `POST /signup`
   - Body:

     ```json
     {
        "First_name": "John",
        "Last_name": "Doe",
        "Password": "securepassword123",
        "Email": "john.doe@example.com",
        "Phone": "1234567890",
        "User_type": "ADMIN"
     }
     ```

2. **Inicio de Sesión**:
   - Endpoint: `POST /login`
   - Body:

     ```json
     {
        "Password":"securepassword123",
        "Email":"john.doe@example.com"
     }
     ```

3. **Acceso a Rutas Protegidas**:
   - Proporciona el token de acceso en el encabezado mediante un parámetro llamado `token`.



## Instrucciones de Desarrollo

1. Instala las dependencias locales:

   ```bash
   go mod tidy
   ```

2. Ejecuta el servidor:

   ```bash
   go run main.go
   ```

3. Ejecuta una base de datos mongo:

   ```bash
   docker run -d --name mongodb-container-auth -p 27017:27017 -v mongo-data:/data/db mongo:latest
   ```

--- 
