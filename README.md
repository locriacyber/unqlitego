unqlitego [![Build Status](https://travis-ci.org/GJRTimmer/unqlitego.svg?branch=master)](https://travis-ci.org/GJRTimmer/unqlitego)
=========

UnQLite Binding for golang.

Install
---------

```sh
$ go get github.com/GJRTimmer/unqlitego
```

Dependencies
------------

Dependencies are management by [Golang Dep](https://github.com/golang/dep/cmd/dep)

```sh
$ go get github.com/golang/dep/cmd/dep
$ dep ensure
```

Test
---------

```sh
$ cd ${GOPATH/:*/}/src/github.com/GJRTimmer/unqlitego
$ go test .
```

Benchmark
----------

```sh
$ cd ${GOPATH/:*/}/src/github.com/GJRTimmer/unqlitego
$ go test -bench Bench*
```

Output:(Lenovo T560 i7-6600U, 16GB Memory, Win:10)
```bash
BenchmarkFileStore-4      500000              5946 ns/op
BenchmarkFileFetch-4      500000              2941 ns/op
BenchmarkMemStore-4      1000000              2695 ns/op
BenchmarkMemFetch-4      1000000              2287 ns/op
```

Output:(Clevo-N850HK1 i7-7700HQ, 32GB Memory, Ubuntu:17.10)
```bash
BenchmarkFileStore-8      500000              3746 ns/op
BenchmarkFileFetch-8     1000000              1940 ns/op
BenchmarkMemStore-8      1000000              1714 ns/op
BenchmarkMemFetch-8      1000000              1715 ns/op
```
