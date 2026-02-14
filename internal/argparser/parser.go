package argparser

import (
	"errors"
	"gocat/internal/cat"
	"strings"
)

func ParseArgs(args []string) (cat.Options, []string, error) {
	var options cat.Options
	var files []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") && arg != "--" {
			err := parseLongArgs(arg, &options)
			if err != nil {
				return cat.Options{}, nil, err
			}
		} else if strings.HasPrefix(arg, "-") && arg != "-" {
			err := parseShortArgs(arg[1:], &options)
			if err != nil {
				return cat.Options{}, nil, err
			}
		} else {
			files = append(files, arg)
		}
	}
	return options, files, nil
}

func parseShortArgs(arg string, opt *cat.Options) error {
	for _, ch := range arg {
		switch ch {
		case 'b':
			opt.NumberNonBlank = true
			opt.Number = false
		case 'n':
			opt.Number = !opt.NumberNonBlank
		case 's':
			opt.SqueezeBlank = true
		case 'e':
			opt.ShowEnds = true
			opt.ShowNonPrinting = true
		case 'E':
			opt.ShowEnds = true
		case 't':
			opt.ShowTabs = true
			opt.ShowNonPrinting = true
		case 'T':
			opt.ShowTabs = true
		case 'v':
			opt.ShowNonPrinting = true
		default:
			return errors.New("Unknown flag")
		}
	}
	return nil
}
func parseLongArgs(arg string, opt *cat.Options) error {
	switch arg {
	case "--number-nonblank":
		opt.NumberNonBlank = true
		opt.Number = false
	case "--number":
		opt.Number = !opt.NumberNonBlank
	case "--squeeze-blank":
		opt.SqueezeBlank = true
	default:
		return errors.New("Unknown flag")
	}
	return nil
}
