package pkg

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type HTTPError struct {
	status int
	msg    string
}

func (e *HTTPError) AsFiber(ctx *fiber.Ctx) error {
	if err := ctx.Status(e.status).JSON(&fiber.Map{
		"error": e.msg,
	}); err != nil {
		return fmt.Errorf("failed to write JSON response: %w", err)
	}

	return nil
}

func ToHTTPError(err error) *HTTPError {
	var unwrapper interface{ Unwrap() error }
	if errors.As(err, &unwrapper) {
		unwrappedErr := unwrapper.Unwrap()
		if unwrappedErr != nil {
			return toHTTPError(unwrappedErr)
		}
	}

	var multiUnwrapper interface{ Unwrap() []error }
	if errors.As(err, &multiUnwrapper) {
		var finalErr *HTTPError
		for _, e := range multiUnwrapper.Unwrap() {
			herr := toHTTPError(e)
			if finalErr == nil {
				finalErr = herr

				continue
			}
			if finalErr.status > herr.status {
				finalErr = herr
			}
		}

		return finalErr
	}

	return toHTTPError(err)
}

//nolint:gocyclo,maintidx // Complex error mapping function with many error types
func toHTTPError(err error) *HTTPError {
	var status int
	var msg string

	switch {
	case errors.Is(err, ErrInternalServerError),
		errors.Is(err, ErrFailedToConsumeData),
		errors.Is(err, ErrFailedToUnmarshalEventData),
		errors.Is(err, ErrFailedToLoadConfig),
		errors.Is(err, ErrFailedToGetMotherService),
		errors.Is(err, ErrFailedToGetMotherServices),
		errors.Is(err, ErrNegativePageOrPerPageNotAllowed),
		errors.Is(err, ErrFailedToGetTestScenario),
		errors.Is(err, ErrFailedToGetTestScenarios),
		errors.Is(err, ErrFailedToCreateTestScenario),
		errors.Is(err, ErrFailedToGetTestCategory),
		errors.Is(err, ErrFailedToGetTestCategoryFromRepository),
		errors.Is(err, ErrFailedToGetTestCategoriesFromRepository),
		errors.Is(err, ErrFailedToCreateMotherService),
		errors.Is(err, ErrFailedToSendProvisioningEvent),
		errors.Is(err, ErrFailedToSendEventData):
		status = http.StatusInternalServerError
		msg = InternalServerErrorMessage
	case errors.Is(err, ErrMotherServiceNotFound):
		status = http.StatusNotFound
		msg = MotherServiceNotFound
	case errors.Is(err, ErrTestScenarioNotFound):
		status = http.StatusNotFound
		msg = TestScenarioNotFound
	case errors.Is(err, ErrTestCategoryNotFound):
		status = http.StatusNotFound
		msg = TestCategoryNotFound
	case errors.Is(err, ErrMotherServiceAlreadyExist):
		status = http.StatusConflict
		msg = MotherServiceAlreadyExist
	case errors.Is(err, ErrBadRequest):
		status = http.StatusBadRequest
		msg = InvalidReqBody
	case errors.Is(err, ErrInvalidIDInParams):
		status = http.StatusBadRequest
		msg = InvalidIDInParams

	case errors.Is(err, ErrInvalidMotherServiceName):
		status = http.StatusUnprocessableEntity
		msg = InvalidMotherServiceName
	case errors.Is(err, ErrInvalidExceptionRate):
		status = http.StatusUnprocessableEntity
		msg = InvalidExceptionRate
	case errors.Is(err, ErrInvalidResponseDelayRate):
		status = http.StatusUnprocessableEntity
		msg = InvalidResponseDelayRate
	case errors.Is(err, ErrInvalidDelayConfiguration):
		status = http.StatusUnprocessableEntity
		msg = InvalidDelayConfiguration
	case errors.Is(err, ErrInvalidRandomDelayRange):
		status = http.StatusUnprocessableEntity
		msg = InvalidRandomDelayRange
	case errors.Is(err, ErrInvalidDatabaseName):
		status = http.StatusUnprocessableEntity
		msg = InvalidDatabaseName
	case errors.Is(err, ErrInvalidDatabaseTableName):
		status = http.StatusUnprocessableEntity
		msg = InvalidDatabaseTableName
	case errors.Is(err, ErrExecutionDurationLessThanOne):
		status = http.StatusUnprocessableEntity
		msg = ExecutionDurationLessThanOne
	case errors.Is(err, ErrMaxTestServiceCountLessThanOne):
		status = http.StatusUnprocessableEntity
		msg = MaxTestServiceCountLessThanOne
	case errors.Is(err, ErrAutoStepChangeRateLessThanOne):
		status = http.StatusUnprocessableEntity
		msg = AutoStepChangeRateLessThanOne
	case errors.Is(err, ErrMaxTestServiceCountNotSet):
		status = http.StatusUnprocessableEntity
		msg = MaxTestServiceCountNotSet
	case errors.Is(err, ErrNoNeedMaxTestServiceCount):
		status = http.StatusUnprocessableEntity
		msg = NoNeedMaxTestServiceCount
	case errors.Is(err, ErrExecutionDurationNotSet):
		status = http.StatusUnprocessableEntity
		msg = ExecutionDurationNotSet
	case errors.Is(err, ErrNoNeedExecutionDuration):
		status = http.StatusUnprocessableEntity
		msg = NoNeedExecutionDuration
	case errors.Is(err, ErrAutoStepChangeNotSet):
		status = http.StatusUnprocessableEntity
		msg = AutoStepChangeNotSet
	case errors.Is(err, ErrNoNeedAutoStepChange):
		status = http.StatusUnprocessableEntity
		msg = NoNeedAutoStepChange
	case errors.Is(err, ErrOnlyPendingScenariosCanBeStarted):
		status = http.StatusUnprocessableEntity
		msg = OnlyPendingScenariosCanBeStarted
	case errors.Is(err, ErrPageNotFound):
		status = http.StatusNotFound
		msg = PageNotFound

	case errors.Is(err, ErrInvalidMaxRequest):
		status = http.StatusUnprocessableEntity
		msg = InvalidMaxRequest
	case errors.Is(err, ErrInvalidMaxDuration):
		status = http.StatusUnprocessableEntity
		msg = InvalidMaxDuration
	case errors.Is(err, ErrInvalidBadValueRate):
		status = http.StatusUnprocessableEntity
		msg = InvalidBadValueRate
	case errors.Is(err, ErrInvalidNegativeValueRate):
		status = http.StatusUnprocessableEntity
		msg = InvalidNegativeValueRate
	case errors.Is(err, ErrInvalidRealValueRate):
		status = http.StatusUnprocessableEntity
		msg = InvalidRealValueRate
	case errors.Is(err, ErrInvalidZeroValueRate):
		status = http.StatusUnprocessableEntity
		msg = InvalidZeroValueRate
	case errors.Is(err, ErrInvalidStringValueRate):
		status = http.StatusUnprocessableEntity
		msg = InvalidStringValueRate
	case errors.Is(err, ErrInvalidLongStringValueRate):
		status = http.StatusUnprocessableEntity
		msg = InvalidLongStringValueRate
	case errors.Is(err, ErrInvalidNullValueRate):
		status = http.StatusUnprocessableEntity
		msg = InvalidNullValueRate
	case errors.Is(err, ErrInvalidRequestDelayDurationConfig):
		status = http.StatusUnprocessableEntity
		msg = InvalidRequestDelayDurationConfig
	case errors.Is(err, ErrInvalidRequestDelayDuration):
		status = http.StatusUnprocessableEntity
		msg = InvalidRequestDelayDuration
	case errors.Is(err, ErrMinDelayDurationMoreThanMax):
		status = http.StatusUnprocessableEntity
		msg = MinDelayDurationMoreThanMax
	case errors.Is(err, ErrInvalidTestNumberConfig):
		status = http.StatusUnprocessableEntity
		msg = InvalidTestNumberConfig
	case errors.Is(err, ErrInvalidFixedTestNumber):
		status = http.StatusUnprocessableEntity
		msg = InvalidFixedTestNumber
	case errors.Is(err, ErrInvalidFixedTestNumberConfig):
		status = http.StatusUnprocessableEntity
		msg = InvalidFixedTestNumberConfig
	case errors.Is(err, ErrInvalidMinOrMaxRandomTestNumber):
		status = http.StatusUnprocessableEntity
		msg = InvalidMinOrMaxRandomTestNumber
	case errors.Is(err, ErrMinRandomTestNumberMoreThanMax):
		status = http.StatusUnprocessableEntity
		msg = MinRandomTestNumberMoreThanMax
	case errors.Is(err, ErrInvalidZeroSumOfBadValues):
		status = http.StatusUnprocessableEntity
		msg = InvalidZeroSumOfBadValues
	case errors.Is(err, ErrInvalid100SumOfBadValues):
		status = http.StatusUnprocessableEntity
		msg = Invalid100SumOfBadValues
	case errors.Is(err, ErrTestServiceConfigIsRequired):
		status = http.StatusUnprocessableEntity
		msg = TestServiceConfigIsRequired

	default:
		status = http.StatusInternalServerError
		msg = InternalServerErrorMessage
	}

	return &HTTPError{status, msg}
}

var (
	ErrNotImplemented                  = errors.New("NOT IMPLEMENTED")
	ErrInternalServerError             = errors.New("internal server error")
	ErrBadRequest                      = errors.New("bad request")
	ErrPageNotFound                    = errors.New("404 page not found")
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
	ErrTestScenarioNotFound            = errors.New("test scenario not found")
	ErrFailedToGetTestScenario         = errors.New("failed to get test scenario")
	ErrFailedToGetTestScenarios        = errors.New("failed to get test scenarios")
	ErrFailedToCreateTestScenario      = errors.New("failed to create test scenario")
	ErrTestCategoryNotFound            = errors.New("test category not found")
	ErrFailedToGetTestCategory         = errors.New("failed to get test category record")
	ErrFailedToValidateTestSvcCfg      = errors.New("failed to validate test service config data")
	ErrInvalidIDInParams               = errors.New("invalid id in params")

	// Validation errors.
	ErrInvalidMotherServiceName                = errors.New("mother service name is required")
	ErrInvalidExceptionRate                    = errors.New("exception rate must be between 0 and 100")
	ErrInvalidResponseDelayRate                = errors.New("response delay rate must be between 0 and 100")
	ErrInvalidDelayConfiguration               = errors.New("invalid delay configuration: must be either no delay, fixed delay, or random delay")
	ErrInvalidRandomDelayRange                 = errors.New("random delay min must be less than max")
	ErrInvalidDatabaseName                     = errors.New("database name is required")
	ErrInvalidDatabaseTableName                = errors.New("database table name is required")
	ErrFailedToGetTestCategoriesFromRepository = errors.New("failed to get test categories from repository")
	ErrFailedToGetTestCategoryFromRepository   = errors.New("failed to get test category from repository")
	ErrExecutionDurationLessThanOne            = errors.New("execution duration should be more than one")
	ErrMaxTestServiceCountLessThanOne          = errors.New("max test service count should be more than 1")
	ErrAutoStepChangeRateLessThanOne           = errors.New("auto step change rate should be more than one")
	ErrMaxTestServiceCountNotSet               = errors.New("for this test scenario max test service count should be set")
	ErrNoNeedMaxTestServiceCount               = errors.New("no need to set max service count")
	ErrExecutionDurationNotSet                 = errors.New("for this test scenario execution duration should be set")
	ErrNoNeedExecutionDuration                 = errors.New("no need to set execution duration")
	ErrAutoStepChangeNotSet                    = errors.New("for this test scenario auto step change should be set")
	ErrNoNeedAutoStepChange                    = errors.New("no need to set auto step change rate")
	ErrFailedToSendProvisioningEvent           = errors.New("failed to send provisioning event")
	ErrInvalidMaxRequest                       = errors.New("invalid max request number, couldn't be negative")
	ErrInvalidMaxDuration                      = errors.New("invalid max duration time, couldn't be negative")
	ErrInvalidRequestDelayDuration             = errors.New("invalid request delay duration")
	ErrInvalidRequestDelayDurationConfig       = errors.New("invalid request delay duration configuration")
	ErrMinDelayDurationMoreThanMax             = errors.New("min request random delay duration not be more than max request delay duration")
	ErrInvalidFixedTestNumber                  = errors.New("fixed test number couldn't be negative")
	ErrInvalidFixedTestNumberConfig            = errors.New("invalid fixed test number configuration")
	ErrMinRandomTestNumberMoreThanMax          = errors.New("min random test number couldn't be more than max random test number")
	ErrInvalidMinOrMaxRandomTestNumber         = errors.New("min random test number or max random test number couldn't be negative")
	ErrInvalidBadValueRate                     = errors.New("bad value rate should be between 0 and 100")
	ErrInvalidNegativeValueRate                = errors.New("negative value rate should be between 0 and 100")
	ErrInvalidRealValueRate                    = errors.New("real value rate should be between 0 and 100")
	ErrInvalidZeroValueRate                    = errors.New("zero value rate should be between 0 and 100")
	ErrInvalidStringValueRate                  = errors.New("string value rate should be between 0 and 100")
	ErrInvalidLongStringValueRate              = errors.New("long string value rate should be between 0 and 100")
	ErrInvalidNullValueRate                    = errors.New("null value rate should be between 0 and 100")
	ErrInvalidTestNumberConfig                 = errors.New("invalid test number configuration")
	ErrInvalidZeroSumOfBadValues               = errors.New("sum of all bad values rate should be 0 if bad value rate field is zero")
	ErrInvalid100SumOfBadValues                = errors.New("sum of all bad values should be 100 if bad value field has value")
	ErrTestServiceConfigIsRequired             = errors.New("test service config is required")
	ErrOnlyPendingScenariosCanBeStarted        = errors.New("only pending scenarios can be started")
	ErrFailedToSetScenarioStatusAsRunning      = errors.New("failed to set scenario status as running")
	ErrGettingRunningTestServicesByScenario    = errors.New("failed to get running test services by scenario")
	ErrFailedToDeprovisionTestServices         = errors.New("failed to deprovision test services")
)
