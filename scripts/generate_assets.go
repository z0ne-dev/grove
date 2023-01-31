package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/shurcooL/vfsgen"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cwd, err := os.Getwd()
	must(err)

	dir := filepath.Join(cwd, os.Args[1])

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		must(os.MkdirAll(dir, os.ModePerm))
	}

	fs := http.Dir(dir)

	must(vfsgen.Generate(fs, vfsgen.Options{
		Filename:     filepath.Join(cwd, "assets_generated.go"),
		PackageName:  "resource",
		VariableName: "vfs",
		BuildTags:    "!dev",
	}))
}
