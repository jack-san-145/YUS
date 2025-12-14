package models

type RateLimit struct {
	Key        string
	Capacity   int
	RefillRate int
}
