package pcre_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ando-masaki/go-pcre"
)

func TestCompileGlob(t *testing.T) {
	r, err := pcre.CompileGlob("/**/bin")
	if err != nil {
		t.Fatal(err)
	}

	if !r.MatchString("/bin") {
		t.Error("expected /bin to match")
	}

	if !r.MatchString("/usr/bin") {
		t.Error("expected /usr/bin to match")
	}

	if !r.MatchString("/usr/local/bin") {
		t.Error("expected /usr/local/bin to match")
	}

	if r.MatchString("/usr") {
		t.Error("expected /usr not to match")
	}

	if r.MatchString("/usr/local") {
		t.Error("expected /usr/local not to match")
	}

	if r.MatchString("/home") {
		t.Error("expected /home not to match")
	}
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

func TestGlob(t *testing.T) {
	err := os.MkdirAll("pcretest/dir1", 0o755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll("pcretest/dir2", 0o755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll("pcretest/test1/dir4", 0o755)
	if err != nil {
		t.Fatal(err)
	}

	err = touch("pcretest/file1")
	if err != nil {
		t.Fatal(err)
	}

	err = touch("pcretest/file2")
	if err != nil {
		t.Fatal(err)
	}

	err = touch("pcretest/test1/dir4/text.txt")
	if err != nil {
		t.Fatal(err)
	}

	matches, err := pcre.Glob("pcretest")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 || matches[0] != "pcretest" {
		t.Errorf("expected [pcretest], got %v", matches)
	}

	paths := dirwalk("pcretest")
	log.Printf("paths:%#v\n", paths)

	matches, err = pcre.Glob("pcretest/dir*")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 2 ||
		matches[0] != "pcretest/dir1" ||
		matches[1] != "pcretest/dir2" {
		t.Errorf("expected [pcretest/dir1 pcretest/dir2], got %v", matches)
	}

	matches, err = pcre.Glob("pcretest/file*")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 2 ||
		matches[0] != "pcretest/file1" ||
		matches[1] != "pcretest/file2" {
		t.Errorf("expected [pcretest/file1 pcretest/file2], got %v", matches)
	}

	matches, err = pcre.Glob("pcretest/**/*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 ||
		matches[0] != "pcretest/test1/dir4/text.txt" {
		t.Errorf("expected [pcretest/test1/dir4/text.txt], got %v", matches)
	}

	err = os.RemoveAll("pcretest")
	if err != nil {
		t.Fatal(err)
	}
}

func touch(path string) error {
	fl, err := os.OpenFile(path, os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	return fl.Close()
}
