package user

import "fmt"

type IncorrectLoginOrPassword struct {
	Login    string
	Password string
}

func (IncorrectLoginOrPassword) Error() string {
	return "Incorrect login or password"
}

type LoginAlreadyExists struct {
	Login string
}

func (e LoginAlreadyExists) Error() string {
	return fmt.Sprintf("Login %s already exists", e.Login)
}
