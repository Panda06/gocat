package cat

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type state struct {
	lineNum      int
	lastSym      byte
	prevWasEmpty bool
}

func Run(opt Options, files []string) {
	s := state{lineNum: 1, lastSym: '\n'}

	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			continue
		}
		processFile(file, opt, &s)
		file.Close()
	}
}

func processFile(file *os.File, opt Options, s *state) {
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			processLine(line, opt, s)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			break
		}
	}
}

func processLine(line []byte, opt Options, s *state) {
	currentEmpty := len(line) == 0 || (len(line) == 1 && line[0] == '\n')

	if opt.SqueezeBlank && s.prevWasEmpty && currentEmpty {
		s.prevWasEmpty = currentEmpty
		return
	}

	if (opt.Number || (opt.NumberNonBlank && !currentEmpty)) && s.lastSym == '\n' {
		fmt.Printf("%6d\t", s.lineNum)
		s.lineNum++
	}

	printLine(line, opt)
	s.prevWasEmpty = currentEmpty
	s.lastSym = line[len(line)-1]
}

func printLine(line []byte, opt Options) {
	for _, c := range line {
		if opt.ShowEnds && c == '\n' {
			fmt.Print("$")
		}
		if c == '\t' && opt.ShowTabs {
			fmt.Print("^I")
		} else if opt.ShowNonPrinting && (c < 32 || c > 126) && c != '\n' && c != '\t' {
			printNonVisible(c)
		} else {
			fmt.Printf("%c", c)
		}
	}
}

func printNonVisible(c byte) {
	if c == 127 {
		fmt.Print("^?")
	} else if c < 32 {
		fmt.Printf("^%c", c+64)
	} else if c > 127 && c < 160 {
		fmt.Printf("M-^%c", c-64)
	} else if c == 255 {
		fmt.Print("M-^?")
	} else if c >= 160 {
		fmt.Printf("M-%c", c-128)
	}
}
