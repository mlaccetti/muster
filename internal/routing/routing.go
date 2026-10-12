package routing

import (
	"encoding/json"
	"fmt"
	"io"
)

type Line struct {
	Verdict    string  `json:"verdict"`
	Harness    *string `json:"harness"`
	Family     *string `json:"family"`
	Model      *string `json:"model"`
	Report     string  `json:"report"`
	ReasonCode string  `json:"reason_code"`
	Reset      int     `json:"reset"`
	DispatchID string  `json:"dispatch_id,omitempty"`
	Baseline   string  `json:"baseline_sha256,omitempty"`
	Ended      *bool   `json:"ended,omitempty"`
}

func Emit(w io.Writer, line Line) error {
	b, err := json.Marshal(line)
	if err != nil {
		return err
	}
	if len(b) > 1024 {
		return fmt.Errorf("routing line exceeds 1 KiB")
	}
	_, err = fmt.Fprintf(w, "%s\n", string(b))
	return err
}

func ParseLine(s string) (Line, error) {
	var line Line
	if err := json.Unmarshal([]byte(s), &line); err != nil {
		return Line{}, err
	}
	return line, nil
}

func Ptr(s string) *string { return &s }
