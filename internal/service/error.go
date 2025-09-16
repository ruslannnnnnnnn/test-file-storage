package service

type InvalidUUidError struct{}

func (i InvalidUUidError) Error() string {
	return "invalid uuid"
}

type FileNotFoundError struct{}

func (err FileNotFoundError) Error() string {
	return "file not found"
}
