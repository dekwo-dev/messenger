package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type E struct {

}

var GetEnv = env()

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func fatal(l int, file string, f string, m string, e error) {
    if e != nil {
        log.Fatalf(
            "{level: %d, file: %s, function: %s, msg: %s, error: %s}\n",
            l, file, f, m, e.Error(),
        )
    }
    log.Fatalf("{level: %d, file: %s, function: %s, msg: %s}\n", l, file, f, m)
}

func info(l int, file string, f string, m string, e error) {
    if e != nil {
        log.Printf(
            "{level: %d, file: %s, function: %s, msg: %s, error: %s}\n",
            l, file, f, m, e.Error(),
        )
    }
    log.Printf("{level: %d, file: %s, function: %s, msg: %s}\n", l, file, f, m)
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

        defer fd.Close()

        r := bufio.NewReader(fd)
        for {
            if l, e := r.ReadString('\n'); e == io.EOF {
                break
            } else if e != nil {
                fatal(50, file, f, "Error when reading string in .env", e)
            } else {
                sp := strings.SplitN(strings.TrimSpace(l), "=", 2)
                
                if len(sp) != 2 {
                    fatal(50, file, f, fmt.Sprintf("Error when parsing an entry of .env variable. Problem lines: %s",
                        l), nil)
                }

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
