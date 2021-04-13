package core

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"
)

func TestReplacement(t *testing.T) {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	os.Setenv(`TPL_TEST`, `ohoh`)

	sd := `<!--`
	ed := `-->`
	r := ProcessRequest{
		PlaceholderSeparator: `:`,
		StartDelimiter:       sd,
		EndDelimiter:         ed,
	}

	pattern := fmt.Sprintf(regexpTemplate, sd, ed)

	input := `seafood <!-- tpl:{env:TPL_TEST} --> <!-- tpl:{file:/tmp/test} -->  `
	expectedOutput := `seafood ohoh   `

	// f, err := os.OpenFile(r.Destination, os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	t.Errorf("error opening destination file %v", err)
	// }
	// defer f.Close()

	c := processContext{
		Source: r.Source,
		//Destination: f,
		Format: r.Format, //string
		// Inline:               r.Source, //bool
		// Stdout:               r.Source, //bool
		SkipIndent:           r.SkipIndent,           //bool
		PlaceholderSeparator: r.PlaceholderSeparator, //string
		StartDelimiter:       r.StartDelimiter,       //string
		EndDelimiter:         r.EndDelimiter,         //string
		Varfiles:             r.Varfiles,             //[]string
		Vars:                 r.Vars,                 //map[string]string
		re:                   regexp.MustCompile(pattern),
		bytes:                []byte(input),
	}
	result, err := replace(c)
	if err != nil {
		t.Errorf("expected error nil, got %v", err)
	}

	if string(result) != expectedOutput {
		t.Errorf("expected [%s], got [%s]", expectedOutput, string(result))
	}
}

func TestIndentation(t *testing.T) {
	c := processContext{}
	expectedIndentation := `	  `
	orig := []byte(expectedIndentation + `t est`)
	indentation := make([]byte, 0)
	proceed := true
	for _, b := range orig {
		if proceed {
			indentation, proceed = manageIndentation(indentation, b, c)
		}
	}
	fmt.Printf("IND=[%s] \n", string(indentation))
	if string(indentation) != expectedIndentation {
		ef := `
For string           [%s]
Expected indentation [%s]
But got              [%s]`
		t.Errorf(ef, orig, expectedIndentation, indentation)
	}
}
