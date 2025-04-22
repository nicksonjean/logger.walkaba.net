package models

type Context struct {
	CorrelationID string     `json:"correlation_id"`
	RequestID     string     `json:"request_id"`
	AppName       string     `json:"app_name"`
	TagName       string     `json:"tag_name,omitempty"`
	Exception     *Exception `json:"exception,omitempty"`
}
