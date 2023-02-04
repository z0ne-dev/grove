package main

import (
	"github.com/shurcooL/httpfs/union"
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

	assets := filepath.Join(cwd, "frontend", "build")
	templates := filepath.Join(cwd, "templates")

	fs := union.New(map[string]http.FileSystem{
		"/templates": http.Dir(templates),
		"/assets":    http.Dir(assets),
	})

	must(vfsgen.Generate(fs, vfsgen.Options{
		Filename:     filepath.Join(cwd, "assets_generated.go"),
		PackageName:  "resource",
		VariableName: "All",
		BuildTags:    "!dev",
	}))
}
