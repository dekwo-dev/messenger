// TODO -> Swap to Protobuf
package main

import (
	"encoding/json"
	"time"
)

type Comment struct {
    Id      int
    Author  string
    Content string
    Time    time.Time
}

type ViewCountEvent struct {
    Type               string
    ViewCountEventType string
	ViewCount          uint8
}

type DBChangeEvent struct {
    Type        string
    DBEventType string
    Comment     Comment
}

type Event interface {
    getType() string
    serialize() ([]byte, error)
}

func (e *DBChangeEvent) serialize() ([]byte, error) {
	const f = "DBChangeEvent.serialize"
    const file = "event.go"

	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (e *DBChangeEvent) getType() string {
    return e.DBEventType
}

func (e *ViewCountEvent) serialize() ([]byte, error) {
	const f = "ViewCountEvent.serialize"
    const file = "event.go"

	b, err := json.Marshal(e)
	if err != nil {
		return nil, err 
	}

	return b, nil
}

func (e *ViewCountEvent) getType() string {
    return e.Type;
}
