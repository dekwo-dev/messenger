package main

import "time"

const debug bool = true
const remote string = "172.105.103.168" // Testing / Debug
const ipify string = "https://api.ipify.org?format=text"
const port string = "8000"

const NewComment = "on-new-comment"
const DelComment = "on-del-comment"
const NewSub = "on-new-sub"
const DelSub = "on-del-sub"

const writeWait = 10 * time.Second
