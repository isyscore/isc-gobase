package goid

import (
	uuid "github.com/iris-contrib/go.uuid"
)

func GenerateUUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}
