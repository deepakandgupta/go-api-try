package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
)

type RedisCreds struct {
	Name string `json:"name" bson:"name"`
	Username string `json:"username" bson:"username"`
}

func init() {
    gob.Register(RedisCreds{})
}

// go binary encoder
func ToGOB64(name, username string) string {
	var m RedisCreds
	m.Name = name
	m.Username= username
    b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    err := e.Encode(m)
    if err != nil { fmt.Println(`failed gob Encode`, err) }
    return base64.StdEncoding.EncodeToString(b.Bytes())
}

// go binary decoder
func FromGOB64(str string) RedisCreds {
    m := RedisCreds{}
    by, err := base64.StdEncoding.DecodeString(str)
    if err != nil { fmt.Println(`failed base64 Decode`, err); }
    b := bytes.Buffer{}
    b.Write(by)
    d := gob.NewDecoder(&b)
    err = d.Decode(&m)
    if err != nil { fmt.Println(`failed gob Decode`, err); }
    return m
}