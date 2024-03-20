package env

import (
	"bufio"
	"fmt"
	"io"
	"os"
    "strconv"
	"strings"

    . "dekwo.dev/messager/logger"
)

var GetEnv = env()

func env() func(k string) string {
    var env map[string]string

    return func(k string) string {
        const f = "env"
        const file = "env/env.go"

        if env != nil {
            if v, i := env[k]; i {
                return v
            } else {
                return ""
            }
        }

        env := make(map[string]string)

        fd, err := os.Open(".env")
        if err != nil {
            Fatal(50, file, f, "Failed to open environmental variable file .env",
                err)
        }

        defer fd.Close()

        reader := bufio.NewReader(fd)
        for {
            if line, err := reader.ReadString('\n'); err == io.EOF {
                break
            } else if err != nil {
                Fatal(50, file, f, "Error when reading string in .env", err)
            } else {
                sp := strings.SplitN(strings.TrimSpace(line), "=", 2)
                
                if len(sp) != 2 {
                    Fatal(50, file, f, fmt.Sprintf("Error when parsing an entry of .env variable. Problem lines: %s",
                        line), nil)
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

func Debug() bool {
    const f = "Debug"
    const file = "env/env.go"

    debug, e := strconv.ParseBool(GetEnv("DEBUG")); 
    if e != nil {
        Fatal(50, file, f, "Failed to environmental variable DEBUG", e)
    }

    return debug
}

func Prod() bool {
    const f = "Prod"
    const file = "env/env.go"

    prod, e := strconv.ParseBool(GetEnv("PROD")); 
    if e != nil {
        Fatal(50, file, f, "Failed to environmental variable PROD", e)
    }

    return prod 
}

func Port() int64 {
    const f = "Port"
    const file = "env/env.go"

    port, e := strconv.ParseInt(GetEnv("PORT"), 10, 16); 
    if e != nil {
        Fatal(50, file, f, "Failed to environmental variable PROD", e)
    }

    return port 
}
