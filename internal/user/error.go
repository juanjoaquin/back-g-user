package user

import (
	"errors"
	"fmt"
)

// Este archivo son los errores customizados para los campos. Por ejemplo: First Name, Last Name del User, etc...

var ErrFirstNameRequired = errors.New("First Name is required")
var ErrLastNameRequired = errors.New("Last Name is required")

// Manejo de Errores con Parametros Dinamicos
type ErrUserNotFound struct {
	UserID string
}

func (e ErrUserNotFound) Error() string {
	return fmt.Sprintf("user '%s' doesnt exists", e.UserID)
}
