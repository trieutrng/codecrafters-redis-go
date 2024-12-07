package main

import (
	"fmt"
	"sort"
	"strings"
)

func RespTypeString(respType RESPType) string {
	switch respType {
	case SimpleString:
		return "SimpleString"
	case BulkString:
		return "BulkString"
	case Arrays:
		return "Arrays"
	}
	return "Not found"
}

func ToLowerCase(input string) string {
	return strings.ToLower(input)
}

func ValidateStreamId(streamEntry StreamEntry, id string) error {
	lastTime, lastSeq := "0", "0"

	if len(streamEntry) > 0 {
		keys := make([]string, 0, len(streamEntry))
		for k := range streamEntry {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		latest := keys[len(keys)-1]

		fmt.Println(keys, latest, id)

		splitedLatest := strings.Split(latest, "-")
		lastTime, lastSeq = splitedLatest[0], splitedLatest[1]
	}

	splitedNow := strings.Split(id, "-")
	time, seq := splitedNow[0], splitedNow[1]

	if time == "0" && seq == "0" {
		return fmt.Errorf("ERR The ID specified in XADD must be greater than 0-0")
	}

	if time < lastTime || time == lastTime && seq <= lastSeq {
		return fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	return nil
}
