package pkg

import (
	"errors"
)

var (
	ErrInternalServerError         = errors.New("internal server error")
	ErrFailedToSendEventData       = errors.New("failed to send event data")
	ErrFailedToConsumeData         = errors.New("failed to consume data")
	ErrFailedToUnmarshalEventData  = errors.New("failed to unmarshal factorial event")
	ErrFailedToLoadConfig          = errors.New("failed to load config from env")
	ErrFailedToCreateMotherService = errors.New("failed to create mother service item")
	ErrMotherServiceNotFound       = errors.New("mother service not found")
	ErrFailedToGetMotherService    = errors.New("failed to get mother service instance")
)
