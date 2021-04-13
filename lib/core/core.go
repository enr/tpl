package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/enr/go-files/files"
)

var delimitersRegistry = map[string]delimiters{
	`xml`:  {start: `<!--`, end: `-->`},
	`html`: {start: `<!--`, end: `-->`},
	`yml`:  {start: `#!--`, end: `--#`},
	`yaml`: {start: `#!--`, end: `--#`},
	`txt`:  {start: `#!--`, end: `--#`},
	`sh`:   {start: `#!--`, end: `--#`},
	`java`: {start: `/*!--`, end: `--*/`},
	`js`:   {start: `/*!--`, end: `--*/`},
}

const regexpTemplate string = `%s tpl:\{([^\{\}]*)\} %s`

// Process ...
func Process(r ProcessRequest) error {
	c, err := buildProcessContext(r)
	if err != nil {
		return err
	}
	ui.Confidentialf("Using destination path %s", r.Destination)
	bytes, err := replace(c)
	if err != nil {
		return err
	}

	err = writeToDestination(c.Destination, bytes)
	if err != nil {
		return err
	}
	// if n != len(bytes) {
	// 	return fmt.Errorf(`Error writing to %s`, c.Destination.Name())
	// }
	// defer func() {
	// 	w.Flush()
	// }()
	return nil
}

func writeToDestination(d destination, bytes []byte) error {
	if d.Stdout {
		n, err := os.Stdout.Write(bytes)
		if err != nil {
			return err
		}
		if n != len(bytes) {
			return fmt.Errorf(`Error writing to /dev/stdout`)
		}
		return nil
	}
	// w := bufio.NewWriterSize(c.Destination, 4096*2)
	// n, err := c.Destination.Write(bytes) //w.Write(bytes) //
	return ioutil.WriteFile(d.Path, bytes, 0644)
}

func fileExtensionNoDot(p string) string {
	extension := filepath.Ext(p)
	return strings.TrimLeft(extension, ".")
}

func buildSourcePath(p string) (string, error) {
	s, err := filepath.Abs(p)
	if err != nil {
		return p, err
	}
	s = filepath.Clean(s)
	if !files.Exists(s) {
		ui.Errorf("Source file not found: %s", s)
		return s, fmt.Errorf(`Source file not found: %s`, s)
	}
	return s, nil
}

func buildProcessContext(request ProcessRequest) (processContext, error) {
	// --- source
	s, err := buildSourcePath(request.Source)
	if err != nil {
		return processContext{}, errors.New(`Error reading source ` + err.Error())
	}
	ui.Confidentialf(`Source file: %s`, s)
	// --- dest
	err = validateDestinationOptions(request)
	if err != nil {
		return processContext{}, err
	}
	dest, err := resolveDestination(request, s)
	if err != nil || !dest.isInitialized() {
		return processContext{}, errors.New(`No destination set`)
	}

	dat, err := ioutil.ReadFile(s)
	if err != nil {
		return processContext{}, fmt.Errorf(`Error reading source file %s : %v`, s, err)
	}
	sd, ed := resolveDelimiters(s, request)
	pattern := fmt.Sprintf(regexpTemplate, regexp.QuoteMeta(sd), regexp.QuoteMeta(ed))
	ui.Confidentialf("Start process using pattern: %s", pattern)

	pc := newProcessContext(request, pattern, dat)
	pc.Source = s
	pc.Destination = dest
	pc.StartDelimiter = sd
	pc.EndDelimiter = ed
	return pc, nil
}

func resolveDestination(request ProcessRequest, alreadyCheckedSource string) (destination, error) {
	if request.Inline {
		ui.Lifecyclef(`Using source as destination %s`, alreadyCheckedSource)
		return inlineDestination(alreadyCheckedSource)
	} else if request.Stdout {
		ui.Lifecycle(`Using /dev/stdout as destination`)
		return stdoutDestination()
	} else if request.Destination != "" {
		d, err := filepath.Abs(request.Destination)
		if err != nil {
			return destination{}, errors.New(`Error reading destination ` + err.Error())
		}
		d = filepath.Clean(d)
		ui.Lifecyclef(`Using file %s as destination`, d)
		if files.Exists(d) {
			return pathDestination(d)
		}
		err = ensureDestinationExists(d)
		if err != nil {
			return destination{}, err
		}
		return pathDestination(d)
	}
	return destination{}, errors.New(`No destination set`)
}

func validateDestinationOptions(request ProcessRequest) error {
	i := 0
	if request.Inline {
		i++
	}
	if request.Stdout {
		i++
	}
	if request.Destination != "" {
		i++
	}
	if i > 1 {
		return errors.New(`inline stdout and destination path are mutually exclusive`)
	}
	return nil
}

func ensureDestinationExists(d string) error {
	destDir := filepath.Dir(d)
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return errors.New(`Error creating output path ` + err.Error())
	}
	f, err := os.OpenFile(d, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf(`Error creating destination file %s : %v`, d, err)
	}
	defer f.Close()
	return nil
}
