package models

type LogDTO struct {
	BlockHash string `json:"blockHash"`
	EventName string `json:"eventName"`
}
