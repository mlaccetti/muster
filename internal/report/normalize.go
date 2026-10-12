package report

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var blockRE = regexp.MustCompile("(?s)```json\\s*(\\{.*?\\})\\s*```")

func NormalizeFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NormalizeBytes(data)
}

func NormalizeBytes(data []byte) ([]byte, error) {
	trim := strings.TrimSpace(string(data))
	if strings.HasPrefix(trim, "{") {
		return data, nil
	}
	m := blockRE.FindStringSubmatch(trim)
	if len(m) != 2 {
		return nil, fmt.Errorf("report is neither JSON nor markdown block")
	}
	var v any
	if err := json.Unmarshal([]byte(m[1]), &v); err != nil {
		return nil, err
	}
	return []byte(m[1]), nil
}
