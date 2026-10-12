package adapter

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	limitHeadRE        = regexp.MustCompile(`(?i)usage limit|rate.?limit|hit your limit|limit reached|limit will reset`)
	limitErrRE         = regexp.MustCompile(`(?i)usage limit|rate.?limit|\b429\b`)
	cursorLimRE        = regexp.MustCompile(`(?i)resource_exhausted|usage limit`)
	rateLimitBlockedRE = regexp.MustCompile(`(?i)^(rejected|blocked|limit_reached|exceeded)$`)
)

func normReset(t float64, now int64) int {
	if t > 1e11 {
		t = t / 1000
	}
	maxR := int64(8 * 86400)
	if int64(t) > now && int64(t) <= now+maxR {
		return int(t)
	}
	return 0
}

func parseReset(s string, now int64) int {
	if s == "" {
		return 0
	}
	if m := regexp.MustCompile(`\|\s*(\d{9,13})\b`).FindStringSubmatch(s); len(m) == 2 {
		var n float64
		_ = json.Unmarshal([]byte(m[1]), &n)
		if r := normReset(n, now); r != 0 {
			return r
		}
	}
	return 0
}

func textOf(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	default:
		b, _ := json.Marshal(x)
		return string(b)
	}
}

func isTrue(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case float64:
		return x != 0
	default:
		return false
	}
}

func classifyClaude(events []NormalizedEvent, rc int) EndVerdict {
	now := time.Now().Unix()
	var res NormalizedEvent
	rejSeen := false
	rej := 0
	for _, e := range events {
		switch textOf(e["type"]) {
		case "rate_limit_event":
			info, _ := e["rate_limit_info"].(map[string]any)
			if info == nil {
				info = e
			}
			st := textOf(info["status"])
			if rateLimitBlockedRE.MatchString(st) {
				r := normReset(floatFrom(info["resetsAt"]), now)
				rejSeen = true
				if r != 0 || rej == 0 {
					rej = r
				}
			}
		case "result":
			res = e
		}
	}
	subBad := res != nil && res["subtype"] != nil && textOf(res["subtype"]) != "success"
	if res != nil && rc == 0 && !isTrue(res["is_error"]) && !subBad {
		return EndVerdict{Kind: "ok"}
	}
	if rejSeen {
		return EndVerdict{Kind: "limit", Reset: rej}
	}
	if res != nil && isTrue(res["is_error"]) {
		txt := textOf(res["result"])
		head := txt
		if len(head) > 200 {
			head = head[:200]
		}
		if limitHeadRE.MatchString(head) || textOf(res["api_error_status"]) == "429" {
			return EndVerdict{Kind: "limit", Reset: parseReset(txt, now)}
		}
		return EndVerdict{Kind: "error", Reason: "result is_error: " + txt}
	}
	if subBad {
		return EndVerdict{Kind: "error", Reason: "subtype " + textOf(res["subtype"])}
	}
	if rc != 0 {
		return EndVerdict{Kind: "error", Reason: "exit " + itoa(rc)}
	}
	return EndVerdict{Kind: "error", Reason: "no result event"}
}

func classifyCodex(events []NormalizedEvent, rc int) EndVerdict {
	now := time.Now().Unix()
	term := ""
	var errs []string
	for _, e := range events {
		t := textOf(e["type"])
		switch t {
		case "turn.started":
			term = "open"
		case "turn.completed":
			term = "ok"
		case "turn.failed":
			term = "fail"
			errs = append(errs, msgOf(e["error"]))
		case "error":
			errs = append(errs, msgOf(e["message"]))
		default:
			if strings.HasPrefix(t, "item.") {
				item, _ := e["item"].(map[string]any)
				if item != nil && textOf(item["type"]) == "error" {
					errs = append(errs, msgOf(item["message"]))
				}
			}
		}
	}
	if term == "ok" && rc == 0 {
		return EndVerdict{Kind: "ok"}
	}
	for _, err := range errs {
		head := err
		if len(head) > 200 {
			head = head[:200]
		}
		if limitErrRE.MatchString(head) {
			return EndVerdict{Kind: "limit", Reset: parseReset(err, now)}
		}
	}
	if term == "fail" && len(errs) > 0 {
		return EndVerdict{Kind: "error", Reason: "turn failed: " + errs[len(errs)-1]}
	}
	if len(errs) > 0 {
		return EndVerdict{Kind: "error", Reason: "error: " + errs[len(errs)-1]}
	}
	return EndVerdict{Kind: "error", Reason: "no turn.completed"}
}

func classifyCursor(events []NormalizedEvent, rc int) EndVerdict {
	var res NormalizedEvent
	var errs []string
	for _, e := range events {
		switch textOf(e["type"]) {
		case "result":
			res = e
		case "error":
			if e["error"] != nil {
				errs = append(errs, textOf(e["error"]))
			} else {
				errs = append(errs, textOf(e["message"]))
			}
		}
	}
	if res != nil && rc == 0 && !isTrue(res["is_error"]) {
		return EndVerdict{Kind: "ok"}
	}
	bad := append([]string{}, errs...)
	if res != nil && isTrue(res["is_error"]) {
		if res["result"] != nil {
			bad = append(bad, textOf(res["result"]))
		} else {
			bad = append(bad, textOf(res["error"]))
		}
	}
	for _, b := range bad {
		head := b
		if len(head) > 200 {
			head = head[:200]
		}
		if cursorLimRE.MatchString(head) {
			return EndVerdict{Kind: "limit", Reset: 0}
		}
	}
	if len(bad) > 0 {
		return EndVerdict{Kind: "error", Reason: "error: " + bad[len(bad)-1]}
	}
	if rc != 0 {
		return EndVerdict{Kind: "error", Reason: "exit " + itoa(rc)}
	}
	return EndVerdict{Kind: "error", Reason: "no result event"}
}

func msgOf(v any) string {
	if v == nil {
		return ""
	}
	if m, ok := v.(map[string]any); ok {
		if m["message"] != nil {
			return textOf(m["message"])
		}
		return textOf(m)
	}
	return textOf(v)
}

func floatFrom(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case json.Number:
		f, _ := x.Float64()
		return f
	default:
		return 0
	}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
