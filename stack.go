package misc

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

// Deprecated: This variable will be removed in a future version.
var unknown = []byte("???")

// Stack will skip back the provided number of frames and return a stack trace with source code.
// Although we could just use debug.Stack(), this routine will return the source code and
// skip back the provided number of frames - i.e. allowing us to ignore preceding function calls.
// A skip of 0 returns the stack trace for the calling function, not including this call.
// If the problem is a lack of memory of course all this is not going to work...
// Stack returns stack information with source code lines. The skip parameter is used to skip preceding frames (0 means caller).
// The difference from debug.Stack is: this function attempts to output corresponding source code lines simultaneously, making it easier to locate.
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

// functionName converts the provided programCounter into a function name
// functionName converts program counter to simplified function name (removing package path and package name).
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

// source returns a space-trimmed slice of the n'th line.
// source returns the trimmed text of line n (1-based), returns "???" when out of bounds.
//
// Deprecated: This function will be removed in a future version.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return unknown
	}
	return bytes.TrimSpace(lines[n])
}
