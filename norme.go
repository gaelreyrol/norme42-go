package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var NormeDebug bool

type Norme struct {
	File   os.File
	Reader io.Reader
	Errors []*NormeError
	Debug  *bool
}

type NormeError struct {
	Line    int
	Column  int
	Message string
	Level   int
}

func main() {
	norme := new(Norme)
	norme.Debug = flag.Bool("debug", false, "Show debug messages")
	flag.Parse()

	if flag.Args != nil {

		filename := flag.Arg(0)
		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			panic(err)
		}

		if *norme.Debug == true {
			fmt.Println("Open file: " + filename)
		}
		norme.File = bufio.NewReader(f)
	}
}
