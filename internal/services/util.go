package services

import "time"

func generateID() string {
	return time.Now().Format("20060102150405")
}