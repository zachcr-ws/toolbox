package code

import (
	"bytes"
	"encoding/gob"
)

// A general gob encoder
// Example:
// byteData, err = GobGeneralEncoder(input)
func GobGeneralEncoder(input interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(input)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// A general gob decoder
// Example:
// 1) base type
// var baseVar string
// err := GobGeneralDecoder(byteData, "", &baseVar)
//
// 2) custom struct
// var youVar yourStruct
// err := GobGeneralDecoder(byteData, yourStruct{}, &youVar)
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
