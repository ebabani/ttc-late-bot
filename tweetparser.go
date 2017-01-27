package main

import (
	"fmt"
	"regexp"
	"strings"
)

func isSubway(text string) bool {
	regex := regexp.MustCompile("(line 1|line 2|Line 2|Line 1)")

	return regex.Match([]byte(text))
}

func isDelay(text string) bool {
	regex := regexp.MustCompile("(?i:no service|holding on|no subway service)")

	return regex.Match([]byte(text))
}

func isClear(text string) bool {
	regex := regexp.MustCompile("(?i:clear|Clear|CLEAR)")

	return regex.Match([]byte(text))
}

func getStation(text string) string {
	for _, station := range Line1 {
		fmt.Println(station)
		if strings.Contains(text, station) {
			return station
		}
	}

	for _, station := range Line2 {
		if strings.Contains(text, station) {
			return station
		}
	}

	return ""
}
