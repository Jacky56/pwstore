package services

import "pwstore/types"

type Handler struct {
}

func NewHandler() types.Service {
	return &Handler{}
}
