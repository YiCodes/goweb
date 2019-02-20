package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/YiCodes/goweb/codegen"
)

var (
	input string
)

func init() {
	flag.StringVar(&input, "in", ".", "source directory")
}

func genCode() error {
	var err error

	if !filepath.IsAbs(input) {
		input, err = filepath.Abs(input)

		if err != nil {
			return err
		}
	}

	fmt.Printf("in: %v\n", input)

	inputFileInfo, err := os.Stat(input)

	if err != nil {
		return err
	}

	if !inputFileInfo.IsDir() {
		return fmt.Errorf("input is not a directory")
	}

	files, err := ioutil.ReadDir(input)

	if err != nil {
		return err
	}

	parseContext := codegen.NewParseContext()

	var srcFiles []string

	for _, f := range files {
		if f.IsDir() || strings.HasSuffix(f.Name(), ".gen.go") {
			continue
		}

		src := filepath.Join(input, f.Name())

		srcFiles = append(srcFiles, src)
	}

	err = codegen.Parse(parseContext, srcFiles)

	if err != nil {
		return err
	}

	dest := filepath.Join(input, inputFileInfo.Name()) + ".gen.go"

	fmt.Printf("out: %v\n", dest)
	err = codegen.Compile(parseContext, dest)

	if err != nil {
		os.Remove(dest)
		return err
	}

	fmt.Println("complete.")

	return nil
}

func main() {
	flag.Parse()

	err := genCode()

	if err != nil {
		fmt.Println(err)
	}
}
