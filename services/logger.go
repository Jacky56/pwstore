package services

import "pwstore/types"

type Logger struct {
}

func NewLogger() types.Service {
	return &Logger{}
}
