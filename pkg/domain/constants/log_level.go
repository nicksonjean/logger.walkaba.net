package constants

const (
	DebugLevel     = 100
	InfoLevel      = 200
	NoticeLevel    = 250
	WarningLevel   = 300
	ErrorLevel     = 400
	CriticalLevel  = 500
	AlertLevel     = 550
	EmergencyLevel = 600
)

var LevelNames = map[int]string{
	DebugLevel:     "DEBUG",
	InfoLevel:      "INFO",
	NoticeLevel:    "NOTICE",
	WarningLevel:   "WARNING",
	ErrorLevel:     "ERROR",
	CriticalLevel:  "CRITICAL",
	AlertLevel:     "ALERT",
	EmergencyLevel: "EMERGENCY",
}
