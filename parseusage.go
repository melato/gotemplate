package gotemplate

import (
	"bufio"
	"bytes"
	"strings"
)

/*
ParseFuncUsage parses function usage in the format

name1

	description
	...

name2

	...

There may be additional indentation.
This is the format used by "go doc text/template",
under "Predefined global functions..."

if startPrefix is not empty, the parsing begins
after the first line that begins with startPrefix,
and ends at the first line that does not start with white space.
*/
func ParseUsage(data []byte, startPrefix string, usage map[string]string) {
	inFuncs := startPrefix == ""
	var indentName int
	var indentDesc int
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var name string
	var desc bytes.Buffer
	insert := func() {
		if name != "" {
			usage[name] = strings.TrimSpace(desc.String())
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		if !inFuncs {
			if strings.HasPrefix(line, startPrefix) {
				inFuncs = true
			}
			continue
		} else {
			trimmed := strings.TrimLeft(line, " \t")
			indent := len(line) - len(trimmed)
			if name != "" && indent > indentName {
				if indentDesc == 0 {
					indentDesc = indent
				}
				if indent >= indentDesc {
					line = line[indentDesc:]
				}
				if desc.Len() > 0 {
					desc.WriteString("\n")
				}
				desc.WriteString(line)
				continue
			} else {
				insert()
				name = strings.TrimSpace(trimmed)
				indentName = indent
				indentDesc = 0
				desc = bytes.Buffer{}
			}
		}
	}
	insert()
}
