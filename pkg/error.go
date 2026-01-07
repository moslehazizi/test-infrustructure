package pkg

import (
	"errors"
)

var (
	ErrInternalServerError        = errors.New("internal server error")
	ErrFailedToSendEventData      = errors.New("failed to send event data")
	ErrFailedToConsumeData        = errors.New("failed to consume data")
	ErrFailedToUnmarshalEventData = errors.New("failed to unmarshal factorial event")
	ErrFailedToLoadConfig         = errors.New("failed to load config from env")
)
