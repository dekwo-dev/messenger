package main

import (
	"strconv"
	"time"
)

const port int = 8000

const writeWait = 10 * time.Second

func debug() bool {
    const f = "debug"
    const file = "const.go"

    debug, e := strconv.ParseBool(GetEnv("DEBUG")); 
    if e != nil {
        fatal(50, file, f, "Failed to environmental variable DEBUG", e)
    }

    return debug
}

func prod() bool {
    const f = "prod"
    const file = "const.go"

    prod, e := strconv.ParseBool(GetEnv("PROD")); 
    if e != nil {
        fatal(50, file, f, "Failed to environmental variable PROD", e)
    }

    return prod 
}
