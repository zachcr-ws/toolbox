package code

import (
	"bytes"
	"encoding/gob"
)

func GobGeneralEncoder(input interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(input)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func GobGeneralDecoder(input []byte, structType interface{}, value interface{}) error {
	buffer := bytes.NewBuffer(input)
	gob.Register(structType)

	dec := gob.NewDecoder(buffer)
	err := dec.Decode(value)
	if err != nil {
		return err
	}
	return nil
}
