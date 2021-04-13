package core

import (
	"errors"
	"strings"
	"unicode"
)

func resolveDelimiters(source string, request ProcessRequest) (string, string) {
	sd := request.StartDelimiter
	ed := request.EndDelimiter
	// se almeno uno manca cerca di risolvere
	if sd == "" || ed == "" {
		// prima cerca nel registry dei file type
		format := fileExtensionNoDot(source)
		dlmt, ok := delimitersRegistry[format]
		if ok {
			// se sono nel registry sovrascrive ma solo quello non valorizzato
			if sd == "" {
				sd = dlmt.start
			}
			if ed == "" {
				ed = dlmt.end
			}
		}
	}
	// se ancora non sono valorizzati fallback con default
	if sd == "" {
		sd = `${`
	}
	if ed == "" {
		ed = `}`
	}
	return sd, ed
}

func extractPlaceholders(c processContext) []placeholder {

	content := c.bytes
	indexes := c.re.FindAllIndex(content, -1)

	var start int
	var end int
	var text string
	var p placeholder

	placeholders := []placeholder{}

	for _, v := range indexes {
		start = v[0]
		end = v[1]
		text = string(content[start:end])
		p = placeholder{start: start, end: end, text: text}
		placeholders = append(placeholders, p)
	}
	return placeholders
}

func manageIndentation(indentation []byte, b byte, c processContext) ([]byte, bool) {
	if c.SkipIndent {
		return make([]byte, 0), true
	}
	if b == byte('\n') || b == byte('\r') {
		// se char EOL indentation viene azzerata
		return make([]byte, 0), true
	} else if b == ' ' || b == '\t' {
		// se carattere SPACE viene aggiunto a indentation
		return append(indentation, b), true
	}
	// se carattere visibile indentation viene azzerata
	return indentation, false
}

func replace(c processContext) ([]byte, error) {
	orig := c.bytes
	placeholders := extractPlaceholders(c)
	ui.Confidentialf("Found placeholders %v", placeholders)
	new := make([]byte, 0)
	placeholdersIndex := 0
	rep := []byte{}
	var loopError error
	indentation := make([]byte, 0)
	proceed := true
	for i, b := range orig {
		if placeholdersIndex >= len(placeholders) {
			new = append(new, b)
			continue
		}

		if proceed || b == byte('\n') || b == byte('\r') {
			indentation, proceed = manageIndentation(indentation, b, c)
		}
		p := placeholders[placeholdersIndex]

		if i == p.start {
			// la stringa da inserire dovrebbe essere ricavata da placeholder.text
			rep, loopError = replacePlaceholder(p.text, indentation, c)
			if loopError != nil {
				ui.Errorf("Error reading bytes: %v", loopError)
				return new, loopError
			}
			new = append(new, rep...)
			continue
		}
		if i > p.start && i < p.end {
			continue
		}
		if i == p.end {
			placeholdersIndex++
			//continue
		}
		new = append(new, b)
	}
	return new, nil
}

func replacePlaceholder(text string, indentation []byte, c processContext) ([]byte, error) {
	key, expression, defaultValue, err := tokenizePlaceholder(text, c)
	if err != nil {
		return []byte{}, err
	}
	replacer, ok := placeholderReplacersRegistry[key]
	if !ok {
		// oppure potrebbe dare default value o lasciare placeholder
		return []byte{}, errors.New(`Replacer not found for key ` + key)
	}
	ui.Confidentialf("Replacer %v", replacer)
	pr := placeholderReplacement{
		expression:   expression,
		defaultValue: defaultValue,
		indentation:  indentation,
	}
	value, err := replacer(pr, c)
	if err != nil {
		return []byte{}, err
	}
	ui.Confidentialf("Replaced placeholder %s with:\n%s", text, value)
	return []byte(value), nil
}

func tokenizePlaceholder(text string, c processContext) (string, string, string, error) {
	ui.Confidentialf("Text: '%s' delimiters: '%s' '%s'", text, c.StartDelimiter, c.EndDelimiter)
	p := strings.ReplaceAll(text, `} `+c.EndDelimiter, ``)
	p = strings.ReplaceAll(p, c.StartDelimiter+` tpl:{`, ``)

	sep := c.PlaceholderSeparator
	if redeclareSeparator(p) {
		sep = p[0:1]
		p = p[1:]
	}
	tokens := strings.SplitN(p, sep, 3)
	tl := len(tokens)
	if tl < 2 {
		return "", "", "", errors.New(`Invalid text ` + text)
	}
	defaultValue := ""
	if tl > 2 {
		defaultValue = tokens[2]
	}
	return tokens[0], tokens[1], defaultValue, nil
}

// primo carattere non alfanumerico quindi e' il nuovo separatore
func redeclareSeparator(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return true
		}
		break
	}
	return false
}
