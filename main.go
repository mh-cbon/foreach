package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {

	var opts = struct {
		insensitive bool
	}{}

	flag.BoolVar(&opts.insensitive, "i", true, "regexp match case insensitive")

	flag.Parse()

	args := flag.Args()

	//look for "- as word '\s+'"
	sdsn := args[0]
	nameStr := args[2]
	splitStr := args[3]
	cmdLine := append([]string{}, args[4:]...)

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

	//process the input
	scanner := bufio.NewScanner(src)
	scanner.Split(ScanRegexp(splitBy))
	index := 0
	for scanner.Scan() {
		thing := scanner.Text()

		callArgs := append([]string{}, binArgs...)
		for i, callArg := range callArgs {
			callArgs[i] = mustExecTemplate(callArg, varName, thing, index)
		}

		cmd := exec.Command(bin, callArgs...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			cmd := strings.Join(cmdLine, " ")
			log.Fatalf("failed to execute the command %q, err=%v", cmd, err)
		}
		index++
	}
}

func mustExecTemplate(tpl, varName, value string, index int) string {

	t := template.Must(
		template.New("").Parse(tpl),
	)
	data := map[string]interface{}{
		"index": index,
	}
	data[varName] = value
	out := new(bytes.Buffer)
	t = template.Must(
		t,
		t.Execute(out, data),
	)
	return out.String()
}

// ScanRegexp is based on ScanLines.
func ScanRegexp(r *regexp.Regexp) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := r.FindIndex(data); len(i) >= 0 {
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
