package models

type Exception struct {
	Message string `json:"message"`
	File    string `json:"file"`
}
