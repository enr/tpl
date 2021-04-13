package core

import (
	"fmt"
	"regexp"

	"github.com/magiconair/properties"
)

type delimiters struct {
	start string
	end   string
}

type placeholder struct {
	start int
	end   int
	text  string
}

func (p *placeholder) String() string {
	return fmt.Sprintf(`%d %d %s`, p.start, p.end, p.text)
}

type placeholderReplacement struct {
	expression   string
	defaultValue string
	indentation  []byte
}

// ProcessRequest represents the aggregate of user options.
type ProcessRequest struct {
	Source               string
	Destination          string
	Format               string
	Inline               bool
	Stdout               bool
	SkipIndent           bool
	PlaceholderSeparator string
	StartDelimiter       string
	EndDelimiter         string
	Varfiles             []string
	Vars                 map[string]string
}

func newProcessContext(r ProcessRequest, pattern string, bytes []byte) processContext {
	return processContext{
		Source: r.Source,
		//Destination: r.Destination,
		Format: r.Format,
		// Inline:               r.Source,
		// Stdout:               r.Source,
		SkipIndent:           r.SkipIndent,
		PlaceholderSeparator: r.PlaceholderSeparator,
		StartDelimiter:       r.StartDelimiter,
		EndDelimiter:         r.EndDelimiter,
		Varfiles:             r.Varfiles,
		Vars:                 r.Vars,
		re:                   regexp.MustCompile(pattern),
		bytes:                bytes,
	}
}

type destination struct {
	Path   string
	Inline bool
	Stdout bool
}

func (d *destination) isInitialized() bool {
	return d.Inline || d.Stdout || d.Path != ""
}

func stdoutDestination() (destination, error) {
	return destination{
		Inline: false,
		Stdout: true,
	}, nil
}
func inlineDestination(source string) (destination, error) {
	ui.Warnf(`inline destination path %s`, source)
	return destination{
		Inline: true,
		Stdout: false,
		Path:   source,
	}, nil
}
func pathDestination(d string) (destination, error) {
	return destination{
		Inline: false,
		Stdout: false,
		Path:   d,
	}, nil
}

type processContext struct {
	Source      string
	Destination destination
	Format      string
	// Inline               bool
	// Stdout               bool
	SkipIndent           bool
	PlaceholderSeparator string
	StartDelimiter       string
	EndDelimiter         string
	Varfiles             []string
	Vars                 map[string]string
	re                   *regexp.Regexp
	bytes                []byte
}

func (c *processContext) variables() *variables {
	ignoreMissing := true
	p := properties.MustLoadFiles(c.Varfiles, properties.UTF8, ignoreMissing)
	return &variables{
		fileProperties:  p,
		optionVariables: c.Vars,
	}
}

type placeholderReplacerFunc func(r placeholderReplacement, c processContext) (string, error)

type variables struct {
	fileProperties *properties.Properties
	// option vars hanno priorita'
	optionVariables map[string]string
}

func (v *variables) addVariable(key string, value string) {
	v.optionVariables[key] = value
}

func (v *variables) get(key string) (string, bool) {
	// check in vars
	value, ok := v.optionVariables[key]
	if ok {
		return value, true
	}
	// check in properties
	return v.fileProperties.Get(key)
}
