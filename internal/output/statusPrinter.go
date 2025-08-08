package output

import (
	"fmt"
	"github.com/fatih/color"
)

func PrintStatus(category, key string, args ...interface{}) {
	printMessage(category, key, color.New(color.FgWhite), args...)
}

func PrintSuccess(category, key string, args ...interface{}) {
	printMessage(category, key, color.New(color.FgGreen), args...)
}

func PrintError(category, key string, args ...interface{}) {
	printMessage(category, key, color.New(color.FgRed), args...)
}

func PrintWarning(category, key string, args ...interface{}) {
	printMessage(category, key, color.New(color.FgYellow), args...)
}

func printMessage(category, key string, c *color.Color, args ...interface{}) {
	cat, ok := StatusMessages[category]
	if !ok {
		fmt.Println("Unknown category:", category)
		return
	}
	msg, ok := cat.Messages[key]
	if !ok {
		fmt.Println("Unknown message key:", key)
		return
	}
	full := fmt.Sprintf("%s %s", cat.Prefix, fmt.Sprintf(msg, args...))
	c.Println(full)
}
