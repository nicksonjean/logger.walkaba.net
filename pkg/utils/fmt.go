package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func PrintServerBanner(host string, port int, handlersCount int) {
	lineWidth := 51
	goVersion := runtime.Version()
	httpVersion := strings.Replace(goVersion, "go", "net/http ", 1)

	fmt.Println("                                                      ")
	fmt.Println(" ┌───────────────────────────────────────────────────┐")
	fmt.Printf(" │%s│\n", centerText(httpVersion, lineWidth))
	fmt.Printf(" │%s│\n", centerText(fmt.Sprintf("http://%s:%d", host, port), lineWidth))
	fmt.Printf(" │%s│\n", centerText(fmt.Sprintf("(bound on host 0.0.0.0 and port %d)", port), lineWidth))
	fmt.Println(" │                                                   │")
	handlersPart := fmt.Sprintf("Handlers ............. %-2d", handlersCount)
	processesPart := "Processes .......... 1"
	fmt.Printf(" │ %-25s  %-21s │\n", handlersPart, processesPart)
	preforkPart := "Prefork ....... Disabled"
	pidStr := fmt.Sprintf("%d", os.Getpid())
	availableDotsSpace := 22 - 5 - len(pidStr)
	dots := strings.Repeat(".", availableDotsSpace)
	pidPart := fmt.Sprintf("PID %s %s", dots, pidStr)
	fmt.Printf(" │ %-25s  %-22s │\n", preforkPart, pidPart)
	fmt.Println(" └───────────────────────────────────────────────────┘")
	fmt.Println("                                                      ")
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}

	leftPadding := (width - len(text)) / 2
	rightPadding := width - len(text) - leftPadding

	return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)
}
