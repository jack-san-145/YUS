package models

type RateLimit struct {
	Key             string
	Capacity        int64
	RefillPerSecond int64
	TimeStamp       int64
}
