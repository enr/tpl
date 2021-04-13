package core

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/enr/go-files/files"
)

func assertFilesEqual(t *testing.T, file1 string, file2 string) {
	hash1, _ := hash(file1)
	hash2, _ := hash(file2)
	if hash1 != hash2 {
		t.Errorf("File %s [%s] differs from\n%s [%s]", file1, hash1, file2, hash2)
	}
}

func assertTextFilesEqual(t *testing.T, actual string, exp string) {
	filelines := []string{}
	files.EachLine(actual, func(line string) error {
		filelines = append(filelines, line)
		return nil
	})
	expectedlines := []string{}
	files.EachLine(exp, func(line string) error {
		expectedlines = append(expectedlines, line)
		return nil
	})
	if len(filelines) != len(expectedlines) {
		t.Errorf("EachLine(%s), expected %d lines but got %d", actual, len(expectedlines), len(filelines))
	}
	if len(filelines) == 0 || len(expectedlines) == 0 {
		// probably a missing/unexpected file
		return
	}
	for index, actual := range filelines {
		if len(expectedlines) <= index {
			t.Errorf(`unexpected line %d in file %s`, (index + 1), actual)
			continue
		}
		expected := expectedlines[index]
		if actual != expected {
			t.Errorf(`line %d expected %q but got %q`, (index + 1), expected, actual)
		}
	}
}

func assertStringContains(t *testing.T, s string, substr string) {
	if substr != "" && !strings.Contains(s, substr) {
		t.Fatalf("expected output\n%s\n  does not contain\n%s\n", s, substr)
	}
}

func assertDirectoryContainsOnlyListedFiles(t *testing.T, d string, files []string) {
	aa := []string{}
	err := filepath.Walk(d, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return err
		}
		aa = append(aa, path)
		return err
	})
	if err != nil {
		panic(err)
	}
	if len(files) != len(aa) {
		t.Fatalf("expected %s contains %d files but got %d\n%v%v", d, len(files), len(aa), files, aa)
	}
}

func hash(fullpath string) (string, error) {
	fh, err := os.Open(fullpath)
	defer fh.Close()
	if err != nil {
		return "", err
	}
	h := sha1.New()
	io.Copy(h, fh)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
