package repository

type NoResultError struct{}

func (NoResultError) Error() string {
	return "No result"
}
