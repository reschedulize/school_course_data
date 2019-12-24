package school_course_data

import (
	"bytes"
	"encoding/json"
	"strconv"
)

func trimJSON(input []byte) ([]byte, error) {
	var trimmed bytes.Buffer
	err := json.Compact(&trimmed, input)

	if err != nil {
		return nil, err
	}

	return trimmed.Bytes(), nil
}

func parseTime(input string) (uint16, error) {
	if len(input) != 4 {
		return 0, TimeParseError
	}

	hour, err := strconv.ParseUint(input[:2], 10, 16)

	if err != nil {
		return 0, err
	}

	min, err := strconv.ParseUint(input[2:4], 10, 16)

	if err != nil {
		return 0, nil
	}

	return uint16((hour * 60) + min), nil
}


func maxInt(a int64, b int64) int64 {
	if a < b {
		return b
	} else {
		return a
	}
}

func minInt(a int64, b int64) int64 {
	if a > b {
		return b
	} else {
		return a
	}
}

func createMask(low int64, high int64, start int64, end int64) uint64 {
	var output uint64

	startBit := int64(float64(start) / 5.625)
	endBit := int64(float64(end) / 5.625)

	for i := maxInt(startBit-low, 0); i < minInt(endBit-low, high); i++ {
		output |= 1 << uint64(i)
	}

	return output
}

func createMaskBroad(low int64, high int64, start int64, end int64) uint64 {
	var output uint64

	startBit := int64(float64(start) / 22.5)
	endBit := int64(float64(end) / 22.5)

	for i := maxInt(startBit-low, 0); i < minInt(endBit-low, high); i++ {
		output |= 1 << uint64(i)
	}

	return output
}