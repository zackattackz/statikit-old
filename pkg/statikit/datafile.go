package statikit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ParseDataFormat uint

const (
	JsonFormat ParseDataFormat = iota
	TomlFormat

	DataFileName = "renderData"
)

var extToFormat = map[string]ParseDataFormat{
	".json": JsonFormat,
	".toml": TomlFormat,
}

type MoreThanOneError struct {
	amount uint
}

func (e MoreThanOneError) Error() string {
	return fmt.Sprintf("too many data files: %v", e.amount)
}

func (e MoreThanOneError) Is(target error) bool {
	targetCast, ok := target.(MoreThanOneError)
	if !ok {
		return false
	}
	return targetCast.amount == e.amount
}

func moreThanOneError(amount uint) MoreThanOneError {
	return MoreThanOneError{amount: amount}
}

var (
	ErrDataFileNotExist error = errors.New(fs.ErrNotExist.Error())
)

type ParseDataArgs struct {
	Reader io.Reader
	Format ParseDataFormat
}

/*
Searches for a single valid data file in path `root`

Returns:
	(string, ParseDataFormat, nil) The full path to the single data file and its format
	(_, _, MoreThanOneErr) If more than one data file exists
	(_, _, ErrDataFileNotExist) If no data files exist
	(_, _, error) Any generic error from os calls
*/
func GetDataFilePath(root string) (resPath string, f ParseDataFormat, resErr error) {
	// For each valid extention, ext,
	// If the file at path "root/p.ext" exists and is regular,
	// return the path and it's associated format
	// if no valid file is found return fs.ErrNotExist
	count := uint(0)
	for ext := range extToFormat {
		p := filepath.Join(root, DataFileName+ext)
		s, err := os.Stat(p)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		} else if err != nil {
			resErr = err
			return
		} else if !s.Mode().IsRegular() {
			continue
		} else {
			resPath = p
			f = extToFormat[ext]
			count += 1
		}
	}
	if count == 0 {
		resErr = ErrDataFileNotExist
	} else if count > 1 {
		resErr = moreThanOneError(count)
	}
	return
}

func ParseDataFile(a ParseDataArgs) (any, error) {
	result := make(map[string]interface{})
	switch a.Format {
	case JsonFormat:
		d := json.NewDecoder(a.Reader)
		d.Decode(&result)
		return result, nil
	case TomlFormat:
		d := toml.NewDecoder(a.Reader)
		_, err := d.Decode(&result)
		return result, err
	default:
		return nil, errors.New("invalid format")
	}
}
