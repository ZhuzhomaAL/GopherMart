package repository

import "fmt"

type NoResultError struct{}

func (NoResultError) Error() string {
	return fmt.Sprintf("No result")
}
