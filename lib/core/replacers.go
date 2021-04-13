package core

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/enr/go-files/files"
)

func envReplacer(r placeholderReplacement, c processContext) (string, error) {
	v := os.Getenv(r.expression)
	if v == `` {
		v = r.defaultValue
	}
	return v, nil
}

func valueReplacer(r placeholderReplacement, c processContext) (string, error) {
	return r.expression, nil
}

func varReplacer(r placeholderReplacement, c processContext) (string, error) {
	k := r.expression
	v, ok := c.variables().get(k)
	if !ok {
		return "", errors.New(`Variable not found key ` + k)
	}
	return v, nil
}

func fileReplacer(r placeholderReplacement, c processContext) (string, error) {
	ui.Confidentialf("Expression  [%s]", r.expression)
	ui.Confidentialf("Indentation [%s]", string(r.indentation))
	return f(r.expression, r, c)
}

func varfileReplacer(r placeholderReplacement, c processContext) (string, error) {
	ui.Confidentialf("Expression  [%s]", r.expression)
	ui.Confidentialf("Indentation [%s]", string(r.indentation))
	re := regexp.MustCompile(`\[(.*?)\]`)
	keys := re.FindAllString(r.expression, -1)
	fp := r.expression
	for _, k := range keys {
		k = strings.TrimPrefix(k, `[`)
		k = strings.TrimSuffix(k, `]`)
		v, _ := c.variables().get(k)
		fp = strings.Replace(fp, `[`+k+`]`, v, 1)
	}
	return f(fp, r, c)
}

func f(f string, r placeholderReplacement, c processContext) (string, error) {
	input := ""
	if filepath.IsAbs(f) {
		input = f
	} else {
		parent := filepath.Dir(c.Source)
		input = filepath.Join(parent, f)
		input = filepath.Clean(input)
	}
	ui.Confidentialf(`input %s`, input)
	inds := string(r.indentation)
	lines := []string{}
	first := true
	files.EachLine(input, func(line string) error {
		if first {
			lines = append(lines, line)
			first = false
			return nil
		}
		lines = append(lines, inds+line)
		return nil
	})
	return strings.Join(lines, "\n"), nil
}

var placeholderReplacersRegistry = map[string]placeholderReplacerFunc{
	`value`:   valueReplacer,
	`env`:     envReplacer,
	`file`:    fileReplacer,
	`var`:     varReplacer,
	`varfile`: varfileReplacer,
}
