package middleware

import (
	"fmt"
	"net/http"
	helper "github.com/lucapierini/UserAuthentication/helpers"
	"github.com/gin-gonic/gin"
)
func Authenticate() gin.HandlerFunc {
    // Esta función devuelve un middleware de Gin que autentica las solicitudes
    return func(c *gin.Context) {
        // Obtiene el token del encabezado de la solicitud
        clientToken := c.Request.Header.Get("token")
        // Si no se proporciona el token, responde con un error y aborta la solicitud
        if clientToken == "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprint("No authorization header provided")})
            c.Abort()
            return
        }

        // Valida el token utilizando una función de ayuda (helper)
        claims, err := helper.ValidateToken(clientToken)

        // Si hay un error en la validación del token, responde con un error y aborta la solicitud
        if err != "" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err})
            c.Abort()
            return
        }

        // Si el token es válido, establece los valores de las reclamaciones (claims) en el contexto
        c.Set("email", claims.Email)
        c.Set("first_name", claims.First_name)
        c.Set("last_name", claims.Last_name)
        c.Set("uid", claims.Uid)
        c.Set("user_type", claims.User_type)
        
        // Continúa con el siguiente manejador en la cadena de middleware
        c.Next()
    }
}
