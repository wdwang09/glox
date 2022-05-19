package glox

import (
	"glox/src"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestState(t *testing.T) {
	walkError := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		if !info.IsDir() {
			if info.Name()[len(info.Name())-4:] != ".lox" {
				return nil
			}
			shouldSuccess := !strings.Contains(info.Name(), "error")
			interpreter := glox.NewGlox()
			returnCode := interpreter.RunFile(path)
			if shouldSuccess && returnCode != 0 {
				t.Fatalf("file '%v' should success, but fail", path)
			} else if !shouldSuccess && returnCode == 0 {
				t.Fatalf("file '%v' should fail, but success", path)
			}
		}
		return err
	})
	if walkError != nil {
		t.Fatal(walkError)
	}
}
