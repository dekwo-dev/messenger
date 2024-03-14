package main

import (
	"strconv"
	"time"
)

const remote string = "172.105.103.168" // Testing / Debug
const ipify string = "https://api.ipify.org?format=text"
const port string = "8000"

const NewComment = "on-new-comment"
const DelComment = "on-del-comment"
const NewSub = "on-new-sub"
const DelSub = "on-del-sub"

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
