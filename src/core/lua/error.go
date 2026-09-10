package lua

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func describeLuaError(path string, err error) error {
	src, readErr := os.ReadFile(path)
	if readErr != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	lines := strings.Split(string(src), "\n")

	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", err.Error())

	if line, ok := extractLine(err.Error()); ok && line > 0 && line <= len(lines) {
		writeSnippet(&b, lines, line)
	} else {
		writeTail(&b, lines)
		writeHints(&b, string(src))
	}

	return fmt.Errorf("%s", b.String())
}

func extractLine(msg string) (int, bool) {
	for p := range strings.SplitSeq(msg, ":") {
		if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			return n, true
		}
	}
	return 0, false
}

func writeSnippet(b *strings.Builder, lines []string, line int) {
	start, end := max(0, line-3), min(len(lines), line+2)
	for i := start; i < end; i++ {
		marker := "  "
		if i == line-1 {
			marker = "> "
		}
		fmt.Fprintf(b, "%s%4d | %s\n", marker, i+1, lines[i])
	}
}

func writeTail(b *strings.Builder, lines []string) {
	start := max(0, len(lines)-6)
	fmt.Fprintln(b, "  (error at EOF, showing end of file)")
	for i := start; i < len(lines); i++ {
		fmt.Fprintf(b, "  %4d | %s\n", i+1, lines[i])
	}
}

func writeHints(b *strings.Builder, src string) {
	if o, c := strings.Count(src, "{"), strings.Count(src, "}"); o != c {
		fmt.Fprintf(b, "\n  hint: %d '{' vs %d '}' — likely an unclosed block\n", o, c)
	}
	if f, e := strings.Count(src, "function"), strings.Count(src, "end"); f > e {
		fmt.Fprintf(b, "  hint: %d 'function' vs %d 'end' — likely a missing 'end'\n", f, e)
	}
	if strings.Count(src, `"`)%2 != 0 {
		fmt.Fprintln(b, "  hint: odd number of '\"' — likely an unterminated string")
	}
}
