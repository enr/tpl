package core

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/enr/go-files/files"
)

type processTestCase struct {
	testDir            string
	template           string
	expectedOutputFile string
	varfiles           []string
	vars               map[string]string
}

var processSpecs = []processTestCase{

	{
		testDir:            `hello`,
		template:           `template.txt`,
		expectedOutputFile: `expected_output.txt`,
	},
	{
		testDir:            `yaml`,
		template:           `template.yaml`,
		expectedOutputFile: `expected_output.yaml`,
	},
	{
		testDir:            `js`,
		template:           `template.js`,
		expectedOutputFile: `expected_output.js`,
	},
	{
		testDir:            `complete`,
		template:           `template.txt`,
		expectedOutputFile: `expected_output.txt`,
		varfiles: []string{
			`test.properties`,
		},
		vars: map[string]string{
			`foo`: `bar bar`,
		},
	},
}

func TestProcess(t *testing.T) {
	testdataBase := `../../testdata`
	for idx, spec := range processSpecs {
		destination := fmt.Sprintf(`/tmp/test-%s.txt`, spec.testDir)

		p, err := filepath.Abs(filepath.Join(testdataBase, spec.testDir, spec.template))
		if err != nil {
			log.Println(err)
		}
		exp, err := filepath.Abs(filepath.Join(testdataBase, spec.testDir, spec.expectedOutputFile))
		if err != nil {
			log.Println(err)
		}
		var f string
		vf := []string{}
		if len(spec.varfiles) > 0 {
			for _, v := range spec.varfiles {
				f, err = filepath.Abs(filepath.Join(testdataBase, spec.testDir, v))
				if err != nil {
					log.Println(err)
				}
				vf = append(vf, f)
			}
		}
		r := ProcessRequest{
			Source:               p,
			Destination:          destination,
			Varfiles:             vf,
			Vars:                 spec.vars,
			PlaceholderSeparator: `:`,
		}
		err = Process(r)

		if err != nil {
			t.Errorf("spec %d <%s> expected error nil, got %v", idx, spec.testDir, err)
		}

		assertTextFilesEqual(t, destination, exp)
	}
}

func TestInline(t *testing.T) {
	testdataBase := `../../testdata`

	source, err := filepath.Abs(filepath.Join(testdataBase, `hello`, `template.txt`))
	if err != nil {
		t.Errorf(`%v`, err)
	}
	err = files.Copy(source, `/tmp/test-inline-hello.txt`)
	if err != nil {
		t.Errorf(`%v`, err)
	}
	r := ProcessRequest{
		Source:               `/tmp/test-inline-hello.txt`,
		Inline:               true,
		PlaceholderSeparator: `:`,
	}
	err = Process(r)

	if err != nil {
		t.Errorf("%v", err)
	}

	assertTextFilesEqual(t, `/tmp/test-inline-hello.txt`, `../../testdata/hello/expected_output.txt`)
}
