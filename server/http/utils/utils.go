package utils

import (
	"io"
	"log"
	"net/http"
	"strings"

	"example.com/guestbook/config"
)

var ipBuilder strings.Builder

func GetSelfPublicIP() string {
    ipBuilder.Reset()

    ipBuilder.WriteByte(':')

    var (
        buffer []byte
        err error
    )

    if !config.Debug {
        ipBuilder.Reset()
        var response *http.Response
        response, err = http.Get(config.IPLookupService) 
        CheckError(err)

        defer response.Body.Close()
        buffer, err = io.ReadAll(response.Body)
        CheckError(err)
        
        ipBuilder.Write(buffer)
        ipBuilder.WriteByte(':')
    }

    ip := ipBuilder.String()

    return ip
}

func CheckError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
