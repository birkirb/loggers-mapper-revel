# loggers-mapper-revel
Golang [Loggers](https://gopkg.in/birkirb/loggers.v1) mapper for [Revel](https://github.com/revel/revel/).

[![GoDoc](https://godoc.org/github.com/birkirb/loggers-mapper-revel?status.svg)](https://godoc.org/github.com/birkirb/loggers-mapper-revel)
[![Build Status](https://travis-ci.org/birkirb/loggers-mapper-revel.svg?branch=master)](http://travis-ci.org/birkirb/loggers-mapper-revel)

## Pre-recquisite

See https://gopkg.in/birkirb/loggers.v1

## Installation

    go get github.com/birkirb/loggers-mapper-revel

## Usage

Include loggers in the same project as your Revel app. After initialization / configuration of your revel app assign the mapper to your own logger (direct or embedded).
It directly accesses and maps Revel's loggers.

### Example

```Go
package app

import (
    "math/rand"
    "time"

    "github.com/revel/revel"

    "gopkg.in/birkirb/loggers.v1/log"
    revlog "github.com/birkirb/loggers-mapper-revel/"
)

func init() {
    // Filters is the default set of global filters.
    revel.Filters = []revel.Filter{
        // Revel filter setup.
    }

    // register startup functions with OnAppStart
    // ( order dependent )
    rand.Seed(time.Now().UTC().UnixNano()) // Seed the Pseudo-Random Generator

    revel.OnAppStart( func() { log.Logger = revlog.NewLogger() } )
    // Other startup
    log.Info("My app has started")
}
```
