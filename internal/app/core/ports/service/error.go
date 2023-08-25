package service

import "fmt"

type NoData struct{}

func (NoData) Error() string {
	return fmt.Sprintf("No data")
}
