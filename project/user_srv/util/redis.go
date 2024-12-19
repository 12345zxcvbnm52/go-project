package util

import (
	"encoding"
	"encoding/json"
)

type RedisParser struct {
	Source interface{}
}

var _ encoding.BinaryMarshaler = new(RedisParser)
var _ encoding.BinaryUnmarshaler = new(RedisParser)

func (u *RedisParser) MarshalBinary() ([]byte, error) {
	return json.Marshal(u.Source)
}

func (u *RedisParser) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u.Source)
}

func Marshal(src interface{}) ([]byte, error) {
	r := &RedisParser{
		Source: src,
	}
	return r.MarshalBinary()
}

func Unmarshal(data []byte, dest interface{}) error {
	r := &RedisParser{
		Source: dest,
	}
	return r.UnmarshalBinary(data)
}
