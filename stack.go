package misc

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

// Deprecated: This variable will be removed in a future version.
var unknown = []byte("???")

// Stack returns a stack trace with source code lines.
// The skip parameter specifies how many stack frames to skip (0 means the caller of Stack).
//
// Unlike debug.Stack(), this function includes the actual source code line for each frame,
// making it easier to debug. However, it requires access to source files.
//
// Note: If source files are not available or memory is insufficient, some lines may show "???".
//
// Deprecated: This function will be removed in a future version. Use runtime.Stack or debug.Stack instead.
func Stack(skip int) string {
	buf := new(bytes.Buffer)

	// Store the last file we opened as its probable that the preceding stack frame
	// will be in the same file
	var lines [][]byte
	var lastFilename string
	for i := skip + 1; ; i++ { // Skip over frames
		programCounter, filename, lineNumber, ok := runtime.Caller(i)
		// If we can't retrieve the information break - basically we're into go internals at this point.
		if !ok {
			break
		}

		// Print equivalent of debug.Stack()
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", filename, lineNumber, programCounter)
		// Now try to print the offending line
		if filename != lastFilename {
			data, err := os.ReadFile(filename)
			if err != nil {
				// can't read this source file
				// likely we don't have the sourcecode available
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFilename = filename
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", functionName(programCounter), source(lines, lineNumber))
	}
	return buf.String()
}

// functionName converts program counter to a simplified function name.
// It removes the full package path and package name, keeping only the function name.
// For example, "github.com/user/pkg.MyFunc" becomes "MyFunc".
//
// Deprecated: This function will be removed in a future version.
func functionName(programCounter uintptr) []byte {
	function := runtime.FuncForPC(programCounter)
	if function == nil {
		return unknown
	}
	name := []byte(function.Name())

	// Because we provide the filename we can drop the preceding package name.
	if lastslash := bytes.LastIndex(name, []byte("/")); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	// And the current package name.
	if period := bytes.Index(name, []byte(".")); period >= 0 {
		name = name[period+1:]
	}
	// And we should just replace the interpunct with a dot
	name = bytes.ReplaceAll(name, []byte("·"), []byte("."))
	return name
}

// source returns the trimmed text of the nth line (1-based line number).
// Returns "???" if the line number is out of bounds.
//
// Deprecated: This function will be removed in a future version.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return unknown
	}
	return bytes.TrimSpace(lines[n])
}
