package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	html "html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	text "text/template"

	"github.com/Masterminds/sprig"
	"github.com/leekchan/gtf"
)

func main() {

	var opts = struct {
		insensitive bool
		strict      bool
		funcs       string
		kind        string
	}{}

	flag.BoolVar(&opts.insensitive, "i", true, "regexp match case insensitive")
	flag.BoolVar(&opts.strict, "s", false, "enable strict variable naming")
	flag.StringVar(&opts.funcs, "funcs", "sprig", "the set of functions to import for the template processing, one of sprig|gtf")
	flag.StringVar(&opts.kind, "kind", "html", "the kind of html processing, one of text|html")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("no command to run")
	}

	//look for "- as word '\s+'"
	sdsn := args[0]
	nameStr := args[2]
	splitStr := args[3]

	cmdLine := append([]string{}, args[4:]...)
	if cmdLine[0] == "--" {
		cmdLine = cmdLine[1:]
	}

	// open src
	var src io.Reader
	{
		if sdsn == "-" {
			src = os.Stdin
		} else {
			f, err := os.Open(sdsn)
			if err != nil {
				log.Fatalf("failed to open the source file %q, err=%v", sdsn, err)
			}
			defer f.Close()
			src = f
		}
	}

	// parse variable names
	var varName string
	{
		varName = strings.TrimSpace(nameStr)
	}
	// if len(names) == 0 {}

	// parse the regexp to split the input
	var splitBy *regexp.Regexp
	{
		if opts.insensitive {
			splitStr = fmt.Sprintf("(?i)%v", splitStr)
		}
		splitBy = regexp.MustCompile(splitStr)
	}

	// parse the command line to execute
	var bin string
	var binArgs []string
	{
		if len(cmdLine) == 0 {
			log.Fatalf("missing command line")
		}
		bin = cmdLine[0]
		binArgs = append([]string{}, cmdLine[1:]...)
	}

	// initialize the func map
	funcMap := map[string]interface{}{}
	{
		if opts.funcs == "sprig" {
			funcMap = sprig.FuncMap()
		} else if opts.funcs == "gtf" {
			funcMap = gtf.GtfFuncMap
		}
	}

	//process the input
	scanner := bufio.NewScanner(src)
	scanner.Split(ScanRegexp(splitBy))
	index := 0
	for scanner.Scan() {
		thing := scanner.Text()

		if thing == "" {
			continue
		}

		callArgs := append([]string{}, binArgs...)
		for i, callArg := range callArgs {
			callArgs[i] = mustExecTemplate(opts.kind, callArg, varName, thing, index, funcMap, opts.strict)
		}

		cmd := exec.Command(bin, callArgs...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if index > 0 {
			log.Println()
			log.Println()
			log.Println("---------------------------------------------------------")
		}

		log.Printf("> %v %v\n", bin, strings.Join(callArgs, " "))

		if err := cmd.Run(); err != nil {
			log.Fatalf("failed to execute the command %v %v, err=%v", bin, strings.Join(callArgs, " "), err)
		}
		index++
	}

	// check for errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to scan the source reader, err=%v", err)
	}

	// exit properly
	os.Exit(0)
}

func mustExecTemplate(kind, tpl, varName, value string, index int, funcMap map[string]interface{}, strict bool) string {

	data := map[string]interface{}{
		"index": index,
	}
	data[varName] = value
	out := new(bytes.Buffer)

	opt := "missingkey=zero"
	if strict {
		opt = "missingkey=error"
	}
	if kind == "html" {
		t := html.Must(
			html.New("").Funcs(funcMap).Option(opt).Parse(tpl),
		)
		t = html.Must(t, t.Execute(out, data))

	} else if kind == "text" {
		t := text.Must(
			text.New("").Funcs(funcMap).Option(opt).Parse(tpl),
		)
		t = text.Must(t, t.Execute(out, data))
	}
	return out.String()
}

// ScanRegexp is based on ScanLines.
func ScanRegexp(r *regexp.Regexp) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		i := r.FindIndex(data)
		if len(i) > 0 {
			index := i[0]
			l := i[1] - i[0]
			return index + l, data[0:index], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}
}
