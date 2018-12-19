package common

import uuid "github.com/satori/go.uuid"

// NewID returns a new id
func NewID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
