package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func getSelfPublicIP() string {
	var builder strings.Builder

	builder.Reset()

	var (
		buffer []byte
		err    error
	)

	if !debug {
		builder.Reset()
		var r *http.Response
		r, err = http.Get(ipify)
		checkError(err)

		defer r.Body.Close()
		buffer, err = io.ReadAll(r.Body)
		checkError(err)

		builder.Write(buffer)
	}

	ip := builder.String()

	return ip
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
