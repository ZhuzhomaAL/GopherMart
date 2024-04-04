package service

type NoData struct{}

func (NoData) Error() string {
	return "No data"
}
