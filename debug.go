// Package debug provides an easy mechanism to enable and streamline "print debugging"
// inspired by Rust's [dbg macro].
//
// [dbg macro]: https://doc.rust-lang.org/stable/std/macro.dbg.html
package debug

import (
	"bytes"
	"fmt"
	"go/ast"
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
//	debug.Debug(something) // DEBUG: [/Users/you/projects/myproject/main.go:30:3] something = hello
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

		// We already know the line from runtime.Caller so that gives us an easy filter straight away
		if start.Line == line {
			// If the node we're currently visiting is a call expression
			if parsed, ok := node.(*ast.CallExpr); ok {
				// If it's specifically a call to debug.Debug
				if isDebugCall(parsed.Fun) {
					arg := parsed.Args[0] // Debug takes a single argument
					buf := &bytes.Buffer{}
					printer.Fprint(buf, fset, arg)
					fmt.Fprintf(os.Stderr, "DEBUG: [%v] %v = %v\n", fset.Position(parsed.Fun.Pos()), buf.String(), value)
					return false // Found it
				}
			}
		}

		// We've not found a debug.Debug call yet, keep going
		return true
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
