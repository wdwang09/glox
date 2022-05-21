package glox

import (
	"glox/src"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPassFail(t *testing.T) {
	walkError := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		if !info.IsDir() {
			if info.Name()[len(info.Name())-4:] != ".lox" {
				return nil
			}
			shouldPass := !strings.Contains(info.Name(), "error")
			interpreter := glox.NewGlox()
			returnCode := interpreter.RunFile(path)
			if shouldPass && returnCode != 0 {
				t.Fatalf("file '%v' should pass, but fail", path)
			} else if !shouldPass && returnCode == 0 {
				t.Fatalf("file '%v' should fail, but pass", path)
			}
		}
		return err
	})
	if walkError != nil {
		t.Fatal(walkError)
	}
}
