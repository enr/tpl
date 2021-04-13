package core

import (
	"fmt"
	"regexp"
	"testing"
)

type placeholderTokenizationTestCase struct {
	sd                   string
	ed                   string
	sep                  string
	text                 string
	expectedKey          string
	expectedExpression   string
	expectedDefaultValue string
	expectedSuccess      bool
}

var tokenizationSpecs = []placeholderTokenizationTestCase{
	{
		sd:                   `<!--`,
		ed:                   `-->`,
		sep:                  `:`,
		text:                 `<!-- tpl:{env:ENV_VAR} -->`,
		expectedKey:          `env`,
		expectedExpression:   `ENV_VAR`,
		expectedDefaultValue: ``,
		expectedSuccess:      true,
	},
	{
		sd:                   `<!--`,
		ed:                   `-->`,
		sep:                  `:`,
		text:                 `<!-- tpl:{|file|c:\a\path|c:\default} -->`,
		expectedKey:          `file`,
		expectedExpression:   `c:\a\path`,
		expectedDefaultValue: `c:\default`,
		expectedSuccess:      true,
	},
}

func TestPlaceholdersTokenization(t *testing.T) {

	for idx, spec := range tokenizationSpecs {
		sd := spec.sd
		ed := spec.ed
		r := ProcessRequest{
			PlaceholderSeparator: spec.sep,
			StartDelimiter:       sd,
			EndDelimiter:         ed,
		}

		pattern := fmt.Sprintf(regexpTemplate, sd, ed)
		c := newProcessContext(r, pattern, []byte(`ininfluent`))

		key, expression, defaultValue, err := tokenizePlaceholder(spec.text, c)

		fmt.Printf(">>> k   %s  \n", key)
		fmt.Printf(">>> e   %s  \n", expression)
		fmt.Printf(">>> d   %s  \n", defaultValue)
		fmt.Printf(">>> err   %v  \n", err)

		if spec.expectedSuccess && err != nil {
			t.Errorf(`spec %d unexpected error for success %t: %v`, idx, spec.expectedSuccess, err)
		}
		if !spec.expectedSuccess && err == nil {
			t.Errorf(`spec %d missing expected error for success %t`, idx, spec.expectedSuccess)
		}
		if key != spec.expectedKey {
			t.Errorf(`spec %d key expected [%s] but got [%s]`, idx, spec.expectedKey, key)
		}
		if expression != spec.expectedExpression {
			t.Errorf(`spec %d expression expected [%s] but got [%s]`, idx, spec.expectedExpression, expression)
		}
		if defaultValue != spec.expectedDefaultValue {
			t.Errorf(`spec %d default value expected [%s] but got [%s]`, idx, spec.expectedDefaultValue, defaultValue)
		}

	}

}

func TestPlaceholdersExtraction(t *testing.T) {
	pattern := fmt.Sprintf(regexpTemplate, `<!--`, `-->`)
	c := processContext{
		bytes: []byte(`seafood <!-- tpl:{env:ENV_VAR} --> <!-- tpl:{file:/tmp/test} -->  `),
		re:    regexp.MustCompile(pattern),
	}
	result := extractPlaceholders(c)
	if len(result) != 2 {
		t.Errorf("expected 2 projects, got %d", len(result))
	}
	first := result[0]
	if first.text != `<!-- tpl:{env:ENV_VAR} -->` {
		t.Errorf("expected placeholder text , got %s", first.text)
	}
}

type placeholderResolutionTestCase struct {
	context      processContext
	text         string
	key          string
	expression   string
	defaultValue string
}

var specs = []placeholderResolutionTestCase{
	{
		context:      processContext{PlaceholderSeparator: `:`},
		text:         `env:HOME:~`,
		key:          `env`,
		expression:   `HOME`,
		defaultValue: `~`,
	},
	{
		context:      processContext{PlaceholderSeparator: `:`},
		text:         `var:var.value:foo: bar`,
		key:          `var`,
		expression:   `var.value`,
		defaultValue: `foo: bar`,
	},
	{
		context:      processContext{PlaceholderSeparator: `:`},
		text:         `xxx:zzz`,
		key:          `xxx`,
		expression:   `zzz`,
		defaultValue: ``,
	},
}

func TestPlaceholderResolution(t *testing.T) {

	for idx, spec := range specs {
		c := spec.context
		key, expression, defaultValue, err := tokenizePlaceholder(spec.text, c)
		if err != nil {
			t.Errorf("spec(%d) expected error nil, got %v", idx, err)
		}
		if key != spec.key {
			t.Errorf("spec(%d) expected key %s, got %s", idx, spec.key, key)
		}
		if expression != spec.expression {
			t.Errorf("spec(%d) expected expression %s, got %s", idx, spec.expression, expression)
		}
		if defaultValue != spec.defaultValue {
			t.Errorf("spec(%d) expected defaultValue %s, got %s", idx, spec.defaultValue, defaultValue)
		}
	}
}
