package utils

import (
	"strings"

	"github.com/docker/docker/pkg/stdcopy"
)

func ShortenID(imageID string) string {
	cleanedID := strings.TrimPrefix(imageID, `sha256:`)
	if len(cleanedID) > 12 {
		return cleanedID[0:12]
	}
	return imageID
}

func TrimHeader(message string) string {
	//fmt.Printf(`%d|%d|%d|%d|%d|%d|%d|%d|%d`+"\n", message[0], message[1], message[2], message[3], message[4], message[5], message[6], message[7], message[8])
	if len(message) > 9 && (message[0] == 0 /*stdin*/ ||
		message[0] == 1 /*stdout*/ ||
		message[0] == 2 /*stderr*/) {
		message = message[9:]
	}
	return message
}

var StdCopy = stdcopy.StdCopy
