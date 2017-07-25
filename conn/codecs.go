package conn

import (
	"encoding/json"
	"gopkg.in/vmihailenco/msgpack.v2"
)

const (
	Text   = 1
	Binary = 2
)

type Encoder struct {
	Format    string
	MsgType	  int
	Marshal   func(v interface{}) ([]byte, error)
	Unmarshal func(data []byte, v interface{}) error
}

var (
	JsonCodec *Encoder
	MsgpCodec *Encoder
)

func init() {
	JsonCodec = &Encoder{
		Format: "json",
		MsgType: Text,
		Marshal: func(v interface{}) ([]byte, error) {
			return json.Marshal(v)
		},
		Unmarshal: func(data []byte, v interface{}) error {
			return json.Unmarshal(data, v)
		},
	}

	MsgpCodec = &Encoder{
		Format: "msgpack",
		MsgType: Binary,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(data []byte, v interface{}) error {
			return msgpack.Unmarshal(data, v)
		},
	}

}