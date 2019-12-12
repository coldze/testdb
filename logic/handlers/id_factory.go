package handlers

import "github.com/google/uuid"

type IDFactory func() string

func GuidIDFactory() string {
	return uuid.New().String()
}
