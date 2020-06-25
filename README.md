![leego](.github/leego.jpeg?raw=true)

[![Go Report Card](https://goreportcard.com/badge/github.com/go-wyvern/leego)](https://goreportcard.com/report/github.com/go-wyvern/leego)
[![Build Status](https://api.travis-ci.org/go-wyvern/leego.svg?branch=master)](https://travis-ci.org/github/go-wyvern/leego)
[![codecov](https://codecov.io/gh/go-wyvern/leego/branch/master/graph/badge.svg)](https://codecov.io/gh/go-wyvern/leego)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/go-wyvern/leego/master/LICENSE)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/go-wyvern/leego?tab=doc)

# Supported Go versions

Leego is available as a Go module. Therefore a Go version capable of understanding /vN suffixed imports is required:

1.6.3+

# Feature Overview

- Optimized HTTP router which smartly prioritize routes
- Build robust and scalable RESTful APIs
- Group APIs
- Extensible middleware framework
- Define middleware at root, group or route level
- Data binding for JSON, XML and form payload
- Handy functions to send variety of HTTP responses
- Centralized HTTP error handling
- Template rendering with any template engine
- Define your format for the logger
- Highly customizable
- Automatic TLS via Let’s Encrypt
- HTTP/2 support

# Installation

go get github.com/go-wyvern/leego

# Example

```
package main

import (
  "net/http"

  "github.com/go-wyvern/leego"
  "github.com/go-wyvern/leego/middleware"
)

func main() {
  // Leego instance
  lee:=leego.New()

  // Middleware
  lee.Use(middleware.AddTrailingSlash())

  // Routes
  lee.GET("/", hello)

  // Start server
  e.Start(":1323")
}

// Handler
func hello(c leego.Context) error {
  return c.String(http.StatusOK, "Hello, World!")
}
```
