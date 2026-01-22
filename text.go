package gotemplate

import (
	"bufio"
	"iter"
	"strings"
)

func firstLine(s string) string {
	i := strings.IndexAny(s, "\r\n")
	if i >= 0 {
		return s[:i]
	} else {
		return s
	}
}

// count the number of runes in a string
func runeCount(s string) int {
	var i int
	for i, _ = range s {
	}
	return i + 1
}

func maxRunes(names []string) int {
	var maxlen int
	for _, name := range names {
		w := runeCount(name)
		if w > maxlen {
			maxlen = w
		}
	}
	return maxlen
}

func iterLines(text string) iter.Seq[string] {
	return func(yield func(string) bool) {
		scanner := bufio.NewScanner(strings.NewReader(text))
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}
}
