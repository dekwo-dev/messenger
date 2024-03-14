package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var GetEnv = env()

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func env() func(k string) string {
    var env map[string]string

    return func(k string) string {
        const f = "env"
        const file = "utils.go"

        if env != nil {
            if v, i := env[k]; i {
                return v
            } else {
                return ""
            }
        }

        env := make(map[string]string)

        fd, e := os.Open(".env")
        if e != nil {
            fatal(50, file, f, 
                "Failed to open environmental variable file .env",
                e,
            )
        }

        r := bufio.NewReader(fd)
        for {
            if l, e := r.ReadString('\n'); e == io.EOF {
                break
            } else if e != nil {
                fatal(50, file, f, "Error when reading string in .env", e)
            } else {
                sp := strings.SplitN(strings.TrimSpace(l), "=", 2)
                env[sp[0]] = sp[1]
            }
        }

        if v, i := env[k]; i {
            return v
        } else {
            return ""
        }
    }
}

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

func fatal(l int, file string, f string, m string, e error) {
    log.Fatalf(
        "{level: %d, file: %s, function: %s, msg: %s, error: %s}",
        l, file, f, m, e.Error(),
    )
}

func info(l int, file string, f string, m string, e error) {
    if e != nil {
        log.Printf(
            "{level: %d, file: %s, function: %s, msg: %s, error: %s}",
            l, file, f, m, e.Error(),
        )
    }
    log.Printf("{level: %d, file: %s, function: %s, msg: %s}", l, file, f, m)
}
