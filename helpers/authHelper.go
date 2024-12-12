package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"

)

func CheckUserType(c *gin.Context, role string) (err error) {
    // Obtener el tipo de usuario desde el contexto
    userType := c.GetString("user_type")
    err = nil

    if userType != role {
        err = errors.New("unauthorized to access this resource")
        return err
    }
    return err // Devolver nil si el tipo de usuario coincide
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
    // Obtener el tipo de usuario y el UID desde el contexto
    userType := c.GetString("user_type")
    uid := c.GetString("uid")
    err = nil

    // Verificar si el tipo de usuario es "USER" y el UID no coincide con el userId proporcionado
    if userType == "USER" && uid != userId {
        err = errors.New("unauthorized to access this resource")
        return err // Devolver el error
    }

    // Verificar el tipo de usuario llamando a la función CheckUserType
    err = CheckUserType(c, userType)
    return err
}