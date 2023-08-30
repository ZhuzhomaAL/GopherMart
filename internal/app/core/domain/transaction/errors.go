package transaction

type NotEnoughMoney struct{}

func (NotEnoughMoney) Error() string {
	return "Not enough money"
}
