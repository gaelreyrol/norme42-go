package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

var Debug *bool

type Norme struct {
	Name    string
	File    *os.File
	Reader  *bufio.Reader
	Errors  []*NormeError
	InScope bool
}

type NormeError struct {
	Line    int
	Column  int
	Message string
	Level   int
}

func DebugMessage(message string) {
	if *Debug {
		fmt.Printf("[Debug] %s\n", message)
	}
}

func (this *Norme) PrintErrors() {
	for _, n_error := range this.Errors {
		fmt.Printf("Line: %d | %s\n", n_error.Line, n_error.Message)
	}
	if len(this.Errors) == 0 {
		fmt.Printf("File: %s - OK\n", this.Name)
	}
}

func (this *Norme) CheckHeader(line string, line_number int) {
	var validCommentLine = regexp.MustCompile(`^\/\*.*\*\/\n$`)

	if len(line) > 81 {
		n_error := new(NormeError)
		n_error.Line = line_number
		n_error.Message = "Header: more than 80 chars"
		this.Errors = append(this.Errors, n_error)
	} else if validCommentLine.MatchString(line) == false {
		n_error := new(NormeError)
		n_error.Line = line_number
		n_error.Message = "Header: Not a comment line"
		this.Errors = append(this.Errors, n_error)
	}
}

func (this *Norme) CheckInclude(line string, line_number int) {
	var validIncludeLine = regexp.MustCompile(`^#include [<|"].*[>|"]\n$`)

	if validIncludeLine.MatchString(line) == false {
		n_error := new(NormeError)
		n_error.Line = line_number
		n_error.Message = "Include: bad format"
		this.Errors = append(this.Errors, n_error)
	}
}

func (this *Norme) CheckPrototype(line string, line_number int) {
	var validPrototypeLine = regexp.MustCompile(`^[a-z0-9 _]+[ \t]+[a-z0-9_]+\([a-z0-9 *_]+\)\;\n$`)

	if validPrototypeLine.MatchString(line) == false {
		n_error := new(NormeError)
		n_error.Line = line_number
		n_error.Message = "Prototype: bad format"
		this.Errors = append(this.Errors, n_error)
	}
}

func CheckFile(filename string) {
	norme := new(Norme)
	norme.Name = filename
	f, err := os.Open(filename)
	DebugMessage("Open file: " + filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	norme.File = f
	norme.Reader = bufio.NewReader(f)
	norme.InScope = false
	DebugMessage("New bufio Reader")

	line_count := 1
	for {
		line, err := norme.Reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if line_count < 12 {
			norme.CheckHeader(line, line_count)
		}

		pre_include_check := regexp.MustCompile(`^#.*`)
		pre_func_check := regexp.MustCompile(`^.*\)\n{$`)
		pre_proto_check := regexp.MustCompile(`^.*\)\;\n$`)

		if pre_include_check.MatchString(line) {
			norme.CheckInclude(line, line_count)
		} else if pre_func_check.MatchString(line) {
			DebugMessage("Detect function entry line")
			norme.InScope = true
		} else if pre_proto_check.MatchString(line) && norme.InScope == false {
			DebugMessage("Detect prototype line")
			norme.CheckPrototype(line, line_count)
		}

		line_count++
	}
	norme.PrintErrors()
}

func main() {

	Debug = flag.Bool("debug", false, "Show debug messages")
	flag.Parse()

	for _, filename := range flag.Args() {
		CheckFile(filename)
	}
}
