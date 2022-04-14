package models

type FibbonaciQuery struct {
	X uint64 `json:"x"`
	Y uint64 `json:"y"`
}

type FibbonaciIndex struct {
	MaxIndex uint64 `json:"maxIndex"`
}
