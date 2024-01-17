# Debug

[![License](https://img.shields.io/github/license/FollowTheProcess/debug)](https://github.com/FollowTheProcess/debug)
[![Go Reference](https://pkg.go.dev/badge/github.com/FollowTheProcess/debug.svg)](https://pkg.go.dev/github.com/FollowTheProcess/debug)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/debug)](https://goreportcard.com/report/github.com/FollowTheProcess/debug)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/debug?logo=github&sort=semver)](https://github.com/FollowTheProcess/debug)
[![CI](https://github.com/FollowTheProcess/debug/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/debug/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/debug/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/debug)

Like Rust's [dbg macro], but for Go!

## Project Description

`debug` provides a simple, elegant mechanism to streamline "print debugging", inspired by Rust's [dbg macro]

> [!NOTE]
> This is not intended for logging in production, more to enable quick local debugging while iterating. `debug.Debug` should **not** make it into your production code

## Installation

```shell
go get github.com/FollowTheProcess/debug@latest
```

## Quickstart

### Before

```go
package main

import "fmt"

func main() {
    something := "hello"
    fmt.Printf("something = %s\n", something)
}
```

```shell
something = something
```

- Not ideal, need to manually type variable name
- No source info, could easily get lost in a large program
- Need to deal with fmt verbs

### With `debug.Debug`

```go
package main

import "github.com/FollowTheProcess/debug"

func main() {
    something := "hello"
    debug.Debug(something)
}
```

```shell
DEBUG: [/Users/you/projects/myproject/main.go:7:3] something = "hello"
```

- Source info, yay! 🎉
- Variable name is handled for you
- Can handle any valid go expression or value
- No fmt print verbs

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
[dbg macro]: https://doc.rust-lang.org/stable/std/macro.dbg.html
