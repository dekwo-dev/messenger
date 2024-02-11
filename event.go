// TODO -> Swap to Protobuf
package main

import (
	"encoding/json"
	"log"
)

type Serializeable interface {
    serialize() ([]btye, Error)
}

type SubConnEvent struct {
    event string
    viewCount uint8
}

type DBChangeEvent struct {
    event string
}

func (e *SubConnEvent) serialize() ([]byte, error) {
    const f = "SubConnEvent.serialize"
    b, err := json.Marshal(e)
    if err != nil {
        log.Printf("%v (ERROR): Failed to serialize SubConnEvent\n", f)
        log.Printf("%v (ERROR): %v", f, err)

        return nil, nil
    }
    
    return b, nil;
}

func (e *DBChangeEvent) serializer() ([]byte, error) {
    const f = "DBChangeEvent.serialize"

    b, err := json.Marshal(e)
    if err != nil {
        log.Printf("%v (ERROR): Failed to serialize SubConnEvent\n", f)
        log.Printf("%v (ERROR): %v", f, err)

        return nil, nil;
    }

    return b, nil;
}
