package pkg

import (
	"errors"
)

var (
	ErrInternalServerError             = errors.New("internal server error")
	ErrFailedToSendEventData           = errors.New("failed to send event data")
	ErrFailedToConsumeData             = errors.New("failed to consume data")
	ErrFailedToUnmarshalEventData      = errors.New("failed to unmarshal factorial event")
	ErrFailedToLoadConfig              = errors.New("failed to load config from env")
	ErrFailedToCreateMotherService     = errors.New("failed to create mother service item")
	ErrMotherServiceNotFound           = errors.New("mother service not found")
	ErrFailedToGetMotherService        = errors.New("failed to get mother service instance")
	ErrFailedToGetMotherServices       = errors.New("failed to get mother service instances")
	ErrMotherServiceAlreadyExist       = errors.New("mother service already exist")
	ErrNegativePageOrPerPageNotAllowed = errors.New("negative value for page or per page are not allowed")

	// Test Scenario.
	ErrTestScenarioAlreadyExist = errors.New("test scenario already exist")

	// Validation errors.
	ErrInvalidName                             = errors.New("mother service name is required")
	ErrInvalidExceptionRate                    = errors.New("exception rate must be between 0 and 100")
	ErrInvalidResponseDelayRate                = errors.New("response delay rate must be between 0 and 100")
	ErrInvalidDelayConfiguration               = errors.New("invalid delay configuration: must be either no delay, fixed delay, or random delay")
	ErrInvalidRandomDelayRange                 = errors.New("random delay min must be less than max")
	ErrInvalidDatabaseName                     = errors.New("database name is required")
	ErrInvalidDatabaseTableName                = errors.New("database table name is required")
	ErrInvalidKafkaLivefeedTopic               = errors.New("kafka livefeed topic is required")
	ErrInvalidKafkaFactorialTopic              = errors.New("kafka factorial topic is required")
	ErrFailedToGetTestCategoriesFromRepository = errors.New("failed to get test categories from repository")
)
