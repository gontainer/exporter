[![Go Reference](https://pkg.go.dev/badge/github.com/gontainer/exporter.svg)](https://pkg.go.dev/github.com/gontainer/exporter)
[![Tests](https://github.com/gontainer/exporter/actions/workflows/tests.yml/badge.svg)](https://github.com/gontainer/exporter/actions/workflows/tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/gontainer/exporter/badge.svg?branch=main)](https://coveralls.io/github/gontainer/exporter?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/gontainer/exporter)](https://goreportcard.com/report/github.com/gontainer/exporter)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gontainer_exporter&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gontainer_exporter)

# Exporter

This package provides sets of functions to export variables to a GO code.

```go
s, _ := exporter.Export([3]any{nil, 1.5, "hello world"})
fmt.Println(s)
// Output: [3]interface{}{nil, float64(1.5), "hello world"}
```

See [examples](examples_test.go).
