package utils

import "log"

var DebugMode = true

func DebugPrint(msg string) {
	if DebugMode {
		log.Println("[DEBUG]", msg)
	}
}
