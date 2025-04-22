package models

type CustomLog struct {
	Message   string      `json:"message"`
	Context   Context     `json:"context"`
	Level     int         `json:"level"`
	LevelName string      `json:"level_name"`
	Channel   string      `json:"channel"`
	Datetime  string      `json:"datetime"`
	Extra     interface{} `json:"extra"`
}
