package passenger

import "yus/internal/storage"

type PassengerHandler struct {
	Store *storage.Store
}

func NewPassengerHandler(store *storage.Store) *PassengerHandler {
	return &PassengerHandler{
		Store: store,
	}
}
