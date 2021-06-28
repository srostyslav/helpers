# Helpers package for Go language

This package provides some helpers Go types & functions.

## Installation

Use the `go` command:

	$ go get github.com/srostyslav/helpers

## Requirements

Helpers package tested against Go 1.16.

## Example

```go
package main

import (
    "github.com/srostyslav/helpers"
)

func main() {
    helpers.InitLogger()

    req := &helpers.Request{Url: "https://www.google.com"}
    if err := req.Get(); err != nil {
        helpers.ErrorLogger.Println(err)
    } else {
        helpers.InfoLogger.Println(req.ResponseCode)
    }
}

```
