package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

//Structure used for communication between server and client
//TODO: Change to JSON
type User struct {
	UserType int
	UserData []string
}

//Pack User struct to byte array
func EncodeToBytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Println(err)
	}
	return buf.Bytes()
}

//Pack User struct to string
func EncodeToString(p interface{}) string {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

//Pack byte array to User struct
func DecodeToUser(s []byte) User {
	p := User{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	dec.Decode(&p)
	return p
}
