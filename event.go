// TODO -> Swap to Protobuf
package main

import (
	"encoding/json"
	"log"
)

type ViewCountEvent struct {
    Type      string
    ViewCountEventType string
	ViewCount uint8
}

type DBChangeEvent struct {
    Type        string
    DBEventType string
    Comment Comment
}

type Event interface {
    getType() string
    serialize() ([]byte, error)
}

func (e *ViewCountEvent) serialize() ([]byte, error) {
	const f = "SubConnEvent.serialize"
	b, err := json.Marshal(e)
	if err != nil {
		log.Printf("%v (ERROR): Failed to serialize SubConnEvent\n", f)
		log.Printf("%v (ERROR): %v", f, err)

		return nil, nil
	}

	return b, nil
}

func (e *ViewCountEvent) getType() string {
    return e.Type;
}

func (e *DBChangeEvent) serialize() ([]byte, error) {
	const f = "DBChangeEvent.serialize"

	b, err := json.Marshal(e)
	if err != nil {
		log.Printf("%v (ERROR): Failed to serialize SubConnEvent\n", f)
		log.Printf("%v (ERROR): %v", f, err)

		return nil, nil
	}

	return b, nil
}

func (e *DBChangeEvent) getType() string {
    return e.DBEventType
}
