package user

import "errors"

// Este archivo son los errores customizados para los campos. Por ejemplo: First Name, Last Name del User, etc...

var ErrFirstNameRequired = errors.New("First Name is required")
var ErrLastNameRequired = errors.New("Last Name is required")

var ErrUserNotFound = errors.New("User not found")
