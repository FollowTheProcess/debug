// Package debug provides an easy mechanism to enable and streamline "print debugging"
// inspired by Rust's [dbg macro].
//
// [dbg macro]: https://doc.rust-lang.org/stable/std/macro.dbg.html
package debug

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"runtime"
)

// Debug prints the file, line number and the value to stderr in a nice human-readable format.
//
// Any internal error is printed to stderr and the process halted.
//
//	something := "hello"
//	debug.Debug(something) // DEBUG: [/Users/you/projects/myproject/main.go:30:3] something = "hello"
func Debug(value any) {
	_, file, line, ok := runtime.Caller(1) // Skip: 1 so this file gets skipped
	if !ok {
		fmt.Fprintln(os.Stderr, "DEBUG: Unable to determine caller")
		return
	}

	// Parse the file in question
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUG: Failed to parse %s: %v\n", file, err)
		return
	}

	ast.Inspect(f, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		start := fset.Position(node.Pos())

		// We already know the line from runtime.Caller so if it's not the right line, keep looking
		if start.Line != line {
			return true
		}

		// If the node we're currently visiting is not a function call, keep looking
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		// If the function call is not a call to debug.Debug, keep looking
		if !isDebugCall(call.Fun) {
			return true
		}

		// By now we know it's a function call, and we know the function the user is calling
		// is debug.Debug

		// Debug takes a single argument and it's not variadic so len(parsed.Args) will
		// be enforced at compile time to be 1
		arg := call.Args[0]
		val := fmt.Sprintf("%#v", value)
		buf := &bytes.Buffer{}
		printer.Fprint(buf, fset, arg)
		formatted, err := format.Source([]byte(val))
		if err != nil {
			// If we couldn't format it nicely, just print the raw value
			fmt.Fprintf(os.Stderr, "DEBUG: [%v] %v = %s\n", fset.Position(call.Fun.Pos()), buf.String(), val)
		} else {
			// We could format the value with gofmt, so use that
			fmt.Fprintf(os.Stderr, "DEBUG: [%v] %v = %s\n", fset.Position(call.Fun.Pos()), buf.String(), string(formatted))
		}

		return false // Found it
	})
}

// isDebugCall takes an arbitrary AST expression and determines if it
// was a call to debug.Debug.
func isDebugCall(expr ast.Expr) bool {
	selector, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(selector.X, "debug") && isIdent(selector.Sel, "Debug")
}

// isIdent takes an arbitrary AST expression and a name and determines if
// the expression was assigned to that name.
func isIdent(expr ast.Expr, name string) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == name
}
