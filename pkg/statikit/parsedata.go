package statikit

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/BurntSushi/toml"
)

type ParseDataFormat uint

const (
	JsonFormat ParseDataFormat = iota
	TomlFormat
)

type ParseDataArgs struct {
	r      io.Reader
	format ParseDataFormat
}

func ParseData(a ParseDataArgs) (any, error) {
	result := make(map[string]interface{})
	switch a.format {
	case JsonFormat:
		d := json.NewDecoder(a.r)
		d.Decode(&result)
		return result, nil
	case TomlFormat:
		d := toml.NewDecoder(a.r)
		_, err := d.Decode(&result)
		return result, err
	default:
		return nil, errors.New("invalid format")
	}
}
