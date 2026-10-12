package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type briefInput struct {
	Goal            string   `json:"goal"`
	ModifiablePaths []string `json:"modifiable_paths"`
	DoneWhen        string   `json:"done_when"`
	VerifyCommand   string   `json:"verify_command"`
	OutOfScope      any      `json:"out_of_scope"`
	Worktree        struct {
		Path   string  `json:"path"`
		Branch *string `json:"branch"`
	} `json:"worktree"`
	ReadFirst []string `json:"read_first"`
}

type templateData struct {
	Goal            string
	ModifiablePaths []string
	DoneWhen        string
	VerifyCommand   string
	OutOfScope      string
	WorktreePath    string
	WorktreeBranch  string
	ReadFirst       []string
}

func BriefFromFile(tmplPath string, inputPath string) (string, error) {
	in, err := os.ReadFile(inputPath)
	if err != nil {
		return "", err
	}
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", err
	}
	return BriefFromBytes(string(tmplBytes), in)
}

func BriefFromBytes(tmpl string, input []byte) (string, error) {
	var bi briefInput
	if err := json.Unmarshal(input, &bi); err != nil {
		return "", err
	}
	scope := ""
	switch v := bi.OutOfScope.(type) {
	case string:
		scope = v
	case []interface{}:
		parts := make([]string, 0, len(v))
		for _, x := range v {
			parts = append(parts, fmt.Sprint(x))
		}
		scope = strings.Join(parts, "\n")
	default:
		scope = fmt.Sprint(bi.OutOfScope)
	}
	branch := ""
	if bi.Worktree.Branch != nil {
		branch = *bi.Worktree.Branch
	}
	data := templateData{
		Goal:            bi.Goal,
		ModifiablePaths: bi.ModifiablePaths,
		DoneWhen:        bi.DoneWhen,
		VerifyCommand:   bi.VerifyCommand,
		OutOfScope:      scope,
		WorktreePath:    bi.Worktree.Path,
		WorktreeBranch:  branch,
		ReadFirst:       bi.ReadFirst,
	}
	t, err := template.New("brief").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
