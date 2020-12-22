package workflow

type BusinessError struct {
	Message string
}

func (be BusinessError) Error() string {
	return be.Message
}
