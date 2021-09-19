package statuses

import (
	"errors"
)

// Dictionary type
type Dictionary map[string]string

var (
	errorNotFound   = errors.New("Not Found")
)

func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errorNotFound
}