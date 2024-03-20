package logger

import "log"

func Fatal(l int, file string, f string, m string, e error) {
    if e != nil {
        log.Fatalf(
            "{level: %d, file: %s, function: %s, msg: %s, error: %s}\n",
            l, file, f, m, e.Error(),
        )
    }
    log.Fatalf("{level: %d, file: %s, function: %s, msg: %s}\n", l, file, f, m)
}

func Info(l int, file string, f string, m string, e error) {
    if e != nil {
        log.Printf(
            "{level: %d, file: %s, function: %s, msg: %s, error: %s}\n",
            l, file, f, m, e.Error(),
        )
    }
    log.Printf("{level: %d, file: %s, function: %s, msg: %s}\n", l, file, f, m)
}

