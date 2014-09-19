package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var builddir string
var imports []string = []string{}
var modules []string = []string{}
var modulesName string = "modules"

func main() {
	dir := flag.String("d", "", "please input build dir with")
	flag.Parse()

	l := len(*dir)
	if l == 0 {
		fmt.Println("please input build dir with -d")
		os.Exit(0)
	}
	builddir = fmt.Sprintf("%s/src/%s/", *dir, modulesName)
	filepath.Walk(builddir, walkpath)
	// fmt.Println(imports, modules)

	var istr string = "	\"github.com/rainkid/dogo\"\n"
	for _, v := range imports {
		istr += fmt.Sprintf("	%s \"%s/%s\"\n", v, modulesName, v)
	}

	var rstr string
	for _, v := range modules {
		sm := strings.Split(v, "/")
		c := strings.Title(strings.Split(sm[1], ".")[0])
		if c != "Base" {
			rstr += fmt.Sprintf("	router.AddSampleRoute(\"%s\", &%s.%s{})\n", sm[0], sm[0], c)
		}
	}
	fd, err := os.Create(*dir + "/routers.go")
	defer fd.Close()
	if err != nil {
		fmt.Println("error create file router.go")
		return
	}
	d := "package main\n\nimport (\n" + istr + ")\n\nfunc AddSampleRoute(router *dogo.Router) {\n" + rstr + "}"
	n, err := io.WriteString(fd, d)
	if err != nil {
		fmt.Println(n, err)
	}
	fmt.Println("build success")
}

func walkpath(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() && fi.Name() != modulesName {
		imports = append(imports, fi.Name())
	}
	if !fi.Mode().IsDir() && filepath.Ext(path) == ".go" {
		modules = append(modules, strings.Replace(path, builddir, "", -1))
	}
	return nil
}
