package main

import "fmt"

// Code from https://github.com/segmentio/go-prompt/blob/master/prompt.go
// Avoid package due to unstable deps.

func cliConfirm(prompt string, args ...interface{}) bool {
	for {
		switch scanString(prompt, args...) {
		case "Yes", "yes", "y", "Y":
			return true
		case "No", "no", "n", "N":
			return false
		}
	}
}

func scanString(prompt string, args ...interface{}) string {
	var s string
	fmt.Printf(prompt+": ", args...)
	fmt.Scanln(&s)
	return s
}
