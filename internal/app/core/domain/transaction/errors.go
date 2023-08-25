package transaction

import "fmt"

type NotEnoughMoney struct{}

func (NotEnoughMoney) Error() string {
	return fmt.Sprintf("Not enough money")
}
