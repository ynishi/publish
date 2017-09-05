# Publish [![Build Status](https://travis-ci.org/ynishi/publish.svg?branch=master)](https://travis-ci.org/ynishi/publish)

The publish is command and its libs publish documents to multi web service.
Document block is human can manage, not raw log.

See Godoc at https://godoc.org/github.com/ynishi/publish

## Current status

* Version 1.0(v1)
* basic feature is implemented

## Install

```
go get "github.com/ynihsi/publish/..."
```

## Example

* prepare config file
```
$ cp $GOPATH/src/github.com/ynishi/publish/config.toml .
$ vi config.toml
```
* prepare document file
```
$ touch doc.md
```
* do pubish in same dir
```
$ publish --content=doc.md
```

## Contribute

* Welcome to participate develop, send pull request, add issue(question, bugs, wants and so on).

### Start develop

* fork, clone, develop and pull request.
```
$ git clone ...
$ cd publish
$ go test
```

## Credit and License

* License is Apache-2.0
* Copyright (c) 2017, Yutaka Nishimura.
