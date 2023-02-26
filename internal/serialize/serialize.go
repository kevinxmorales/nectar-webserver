package serialize

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.Info("initializing gob...")
	gob.Register(SX{})
	gob.Register(map[string]interface{}{})
	log.Info("finished initialing gob")
}

type SX map[string]interface{}

func ToGOB64(m SX) (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		return "", fmt.Errorf("failed gob Encode %s", err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func FromGOB64(str string) (SX, error) {
	m := SX{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("failed base64 Decode: %s", err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("failed gob Decode: %s", err)
	}
	return m, nil
}
