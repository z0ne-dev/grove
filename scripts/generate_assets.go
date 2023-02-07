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
	wd, err := os.Getwd()
	must(err)

	println(wd)

	assets := filepath.Join(wd, "frontend", "build")
	templates := filepath.Join(wd, "templates")
	migrations := filepath.Join(wd, "migrations")

	fs := union.New(map[string]http.FileSystem{
		"/templates":  http.Dir(templates),
		"/assets":     http.Dir(assets),
		"/migrations": http.Dir(migrations),
	})

	must(vfsgen.Generate(fs, vfsgen.Options{
		Filename:     filepath.Join(wd, "assets_generated.go"),
		PackageName:  "resource",
		VariableName: "All",
		BuildTags:    "!dev",
	}))
}
