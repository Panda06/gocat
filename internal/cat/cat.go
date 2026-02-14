package cat

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	space         = 32  // первый печатный ASCII
	tilde         = 126 // последний печатный ASCII (~)
	del           = 127 // DEL
	metaCtrlStart = 128 // M-^@ начало extended control
	metaStart     = 160 // начало extended printable
	metaDel       = 255 // M-^?
	ctrlOffset    = 64  // смещение для ^-нотации (^A = 1+64 = 'A')
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
		} else if opt.ShowNonPrinting && (c < space || c > tilde) && c != '\n' && c != '\t' {
			printNonVisible(c)
		} else {
			fmt.Printf("%c", c)
		}
	}
}

func printNonVisible(c byte) {
	if c == del {
		fmt.Print("^?")
	} else if c < space {
		fmt.Printf("^%c", c+ctrlOffset)
	} else if c > del && c < metaStart {
		fmt.Printf("M-^%c", c-ctrlOffset)
	} else if c == metaDel {
		fmt.Print("M-^?")
	} else if c >= metaStart {
		fmt.Printf("M-%c", c-metaCtrlStart)
	}
}
