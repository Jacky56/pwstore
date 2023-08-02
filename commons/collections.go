package commons

import "time"

func IsFull[T any](slice []T) bool {
	return len(slice) == cap(slice)
}

type Alert struct {
	Id   int64
	Text string
}

func NewAlert(text string) Alert {
	return Alert{
		Id:   time.Now().Unix(),
		Text: text,
	}
}
