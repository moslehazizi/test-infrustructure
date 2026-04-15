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
		errors.Is(err, ErrFailedToGetTablesOfDatabase),
		errors.Is(err, ErrFailedToGetDatabases),
		errors.Is(err, ErrFailedToSendEventData),
		errors.Is(err, ErrFailedToUpdateTestScenario),
		errors.Is(err, ErrFailedToGetTestScenariosByStatus),
		errors.Is(err, ErrFailedToGetTestServiceConfig),
		errors.Is(err, ErrFailedToGetReadyForTesting),
		errors.Is(err, ErrFailedToPauseTestService),
		errors.Is(err, ErrFailedToResumeTestService),
		errors.Is(err, ErrFailedToStopTestService):
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
	case errors.Is(err, ErrRunAlreadyInProgress):
		status = http.StatusConflict
		msg = TestServiceAlreadyInProgress
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
	case errors.Is(err, ErrOnlyRunningScenariosCanBePaused):
		status = http.StatusUnprocessableEntity
		msg = OnlyRunningScenariosCanBePaused
	case errors.Is(err, ErrOnlyPausedScenariosCanBeResume):
		status = http.StatusUnprocessableEntity
		msg = OnlyPausedScenariosCanBeReStarted
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
	case errors.Is(err, ErrStartingTestNotImplemented):
		status = http.StatusInternalServerError
		msg = StartingTestNotImplemented
	case errors.Is(err, ErrTestServiceConfigNotFound):
		status = http.StatusNotFound
		msg = TestServiceConfigNotFound
	case errors.Is(err, ErrNumStepsNotSet):
		status = http.StatusUnprocessableEntity
		msg = NumStepsNotSet
	case errors.Is(err, ErrNumStepsShouldBeOne):
		status = http.StatusUnprocessableEntity
		msg = NumStepsShouldBeOne
	case errors.Is(err, ErrIncreaseAgentNumNotBeNegative):
		status = http.StatusUnprocessableEntity
		msg = IncreaseAgentNumNotBeNegative
	case errors.Is(err, ErrExecNumMultiAgentShouldBePositive):
		status = http.StatusUnprocessableEntity
		msg = ExecNumMultiAgentShouldBePositive
	case errors.Is(err, ErrIncreaseFixedInputNotBeNegative):
		status = http.StatusUnprocessableEntity
		msg = IncreaseFixedNumberNotBeNegative
	case errors.Is(err, ErrExecNumMultiFixedInputShouldBePositive):
		status = http.StatusUnprocessableEntity
		msg = ExecNumMultiFixedInputShouldBePositive
	case errors.Is(err, ErrMultiFixedInputConfigNotTrue):
		status = http.StatusUnprocessableEntity
		msg = MultiFixedInputConfigNotTrue
	case errors.Is(err, ErrMultiAgentConfigNotTrue):
		status = http.StatusUnprocessableEntity
		msg = MultiAgentConfigNotTrue

	default:
		status = http.StatusInternalServerError
		msg = InternalServerErrorMessage
	}

	return &HTTPError{status, msg}
}

var (
	ErrFailedToGetTablesOfDatabase      = errors.New("failed to get tables of database")
	ErrFailedToGetDatabases             = errors.New("failed to get databases")
	ErrNotImplemented                   = errors.New("NOT IMPLEMENTED")
	ErrInternalServerError              = errors.New("internal server error")
	ErrBadRequest                       = errors.New("bad request")
	ErrPageNotFound                     = errors.New("404 page not found")
	ErrFailedToSendEventData            = errors.New("failed to send event data")
	ErrFailedToConsumeData              = errors.New("failed to consume data")
	ErrFailedToUnmarshalEventData       = errors.New("failed to unmarshal factorial event")
	ErrFailedToLoadConfig               = errors.New("failed to load config from env")
	ErrFailedToCreateMotherService      = errors.New("failed to create mother service item")
	ErrMotherServiceNotFound            = errors.New("mother service not found")
	ErrFailedToGetMotherService         = errors.New("failed to get mother service instance")
	ErrFailedToGetTestScenariosByStatus = errors.New("failed to get test scenarios by status")
	ErrFailedToGetMotherServices        = errors.New("failed to get mother service instances")
	ErrMotherServiceAlreadyExist        = errors.New("mother service already exist")
	ErrNegativePageOrPerPageNotAllowed  = errors.New("negative value for page or per page are not allowed")
	ErrTestScenarioNotFound             = errors.New("test scenario not found")
	ErrFailedToGetTestScenario          = errors.New("failed to get test scenario")
	ErrFailedToGetTestScenarios         = errors.New("failed to get test scenarios")
	ErrFailedToCreateTestScenario       = errors.New("failed to create test scenario")
	ErrTestCategoryNotFound             = errors.New("test category not found")
	ErrFailedToGetTestCategory          = errors.New("failed to get test category record")
	ErrFailedToValidateTestSvcCfg       = errors.New("failed to validate test service config data")
	ErrInvalidIDInParams                = errors.New("invalid id in params")
	ErrFailedToDeployMotherService      = errors.New("failed to deploy mother service")
	ErrFailedToDeployTestService        = errors.New("failed to deploy test service")
	ErrFailedToDeProvisionMotherService = errors.New("failed to deprovision mother service")
	ErrFailedToDeProvisionTestService   = errors.New("failed to deprovision test service")
	ErrFailedToUpdateTestScenario       = errors.New("failed to update test scenario")
	ErrFailedToGetReadyForTesting       = errors.New("failed to get ready for testing response")
	ErrFailedToPauseTestService         = errors.New("failed to pause test service")
	ErrFailedToResumeTestService        = errors.New("failed to resume test service")
	ErrFailedToStopTestService          = errors.New("failed to stop test services")
	ErrFailedToAbortTestService         = errors.New("failed to abort test services")
	ErrFailedToExecuteSingleScenario    = errors.New("failed to execute single scenario")

	// Validation errors.
	ErrFailedToGetTestServiceConfig                       = errors.New("failed to get test service config by id")
	ErrTestServiceConfigNotFound                          = errors.New("test service config not found")
	ErrInvalidMotherServiceName                           = errors.New("mother service name is required")
	ErrMotherServiceIsNil                                 = errors.New("mother service is nil")
	ErrTestScenarioServiceIsNil                           = errors.New("test scenario service is nil")
	ErrTestServiceConfigIsNil                             = errors.New("test service config is nil")
	ErrInvalidExceptionRate                               = errors.New("exception rate must be between 0 and 100")
	ErrInvalidResponseDelayRate                           = errors.New("response delay rate must be between 0 and 100")
	ErrInvalidDelayConfiguration                          = errors.New("invalid delay configuration: must be either no delay, fixed delay, or random delay")
	ErrInvalidRandomDelayRange                            = errors.New("random delay min must be less than max")
	ErrInvalidDatabaseName                                = errors.New("database name is required")
	ErrInvalidDatabaseTableName                           = errors.New("database table name is required")
	ErrFailedToGetTestCategoriesFromRepository            = errors.New("failed to get test categories from repository")
	ErrFailedToGetTestCategoryFromRepository              = errors.New("failed to get test category from repository")
	ErrExecutionDurationLessThanOne                       = errors.New("execution duration should be more than one")
	ErrMaxTestServiceCountLessThanOne                     = errors.New("max test service count should be more than 1")
	ErrAutoStepChangeRateLessThanOne                      = errors.New("auto step change rate should be more than one")
	ErrMaxTestServiceCountNotSet                          = errors.New("for this test scenario max test service count should be set")
	ErrNoNeedMaxTestServiceCount                          = errors.New("no need to set max service count")
	ErrExecutionDurationNotSet                            = errors.New("for this test scenario execution duration should be set")
	ErrNoNeedExecutionDuration                            = errors.New("no need to set execution duration")
	ErrAutoStepChangeNotSet                               = errors.New("for this test scenario auto step change should be set")
	ErrNoNeedAutoStepChange                               = errors.New("no need to set auto step change rate")
	ErrFailedToSendProvisioningEvent                      = errors.New("failed to send provisioning event")
	ErrInvalidMaxRequest                                  = errors.New("invalid max request number, couldn't be negative")
	ErrInvalidMaxDuration                                 = errors.New("invalid max duration time, couldn't be negative")
	ErrInvalidRequestDelayDuration                        = errors.New("invalid request delay duration")
	ErrInvalidRequestDelayDurationConfig                  = errors.New("invalid request delay duration configuration")
	ErrMinDelayDurationMoreThanMax                        = errors.New("min request random delay duration not be more than max request delay duration")
	ErrInvalidFixedTestNumber                             = errors.New("fixed test number couldn't be negative")
	ErrInvalidFixedTestNumberConfig                       = errors.New("invalid fixed test number configuration")
	ErrMinRandomTestNumberMoreThanMax                     = errors.New("min random test number couldn't be more than max random test number")
	ErrInvalidMinOrMaxRandomTestNumber                    = errors.New("min random test number or max random test number couldn't be negative")
	ErrInvalidBadValueRate                                = errors.New("bad value rate should be between 0 and 100")
	ErrInvalidNegativeValueRate                           = errors.New("negative value rate should be between 0 and 100")
	ErrInvalidRealValueRate                               = errors.New("real value rate should be between 0 and 100")
	ErrInvalidZeroValueRate                               = errors.New("zero value rate should be between 0 and 100")
	ErrInvalidStringValueRate                             = errors.New("string value rate should be between 0 and 100")
	ErrInvalidLongStringValueRate                         = errors.New("long string value rate should be between 0 and 100")
	ErrInvalidNullValueRate                               = errors.New("null value rate should be between 0 and 100")
	ErrInvalidTestNumberConfig                            = errors.New("invalid test number configuration")
	ErrInvalidZeroSumOfBadValues                          = errors.New("sum of all bad values rate should be 0 if bad value rate field is zero")
	ErrInvalid100SumOfBadValues                           = errors.New("sum of all bad values should be 100 if bad value field has value")
	ErrTestServiceConfigIsRequired                        = errors.New("test service config is required")
	ErrOnlyPendingScenariosCanBeStarted                   = errors.New("only pending scenarios can be started")
	ErrOnlyRunningScenariosCanBePaused                    = errors.New("only running scenarios can be paused")
	ErrOnlyPausedScenariosCanBeResume                     = errors.New("only pause scenarios can be resume")
	ErrOnlyRunAndPauseScenariosCanBeStop                  = errors.New("only run and pause can not be stoped")
	ErrAbortedScenariosCanBeAbort                         = errors.New("aborted can not be abort")
	ErrFailedToSetScenarioStatus                          = errors.New("failed to set scenario status")
	ErrGettingRunningTestServicesByScenario               = errors.New("failed to get running test services by scenario")
	ErrFailedToDeprovisionTestServices                    = errors.New("failed to deprovision test services")
	ErrInt32OutOfRange                                    = errors.New("out of int32 range")
	ErrFailedToAddScenarioToExecutionManager              = errors.New("failed to add scenario to execution manager")
	ErrFailedToRunScenarioInExecutionManager              = errors.New("failed to run scenario in execution manager")
	ErrFailedToPauseScenarioToExecutionManager            = errors.New("failed to pause scenario in execution manager")
	ErrFailedToResumeScenarioToExecutionManager           = errors.New("failed to resume scenario in execution manager")
	ErrFailedToStopScenarioToExecutionManager             = errors.New("failed to stop scenario in execution manager")
	ErrFailedToAbortScenarioToExecutionManager            = errors.New("failed to abort scenario in execution manager")
	ErrStartingTestNotImplemented                         = errors.New("starting test not implemented")
	ErrPausingTestNotImplemented                         = errors.New("pausing test not implemented")
	ErrResumingTestNotImplemented                        = errors.New("resuming test not implemented")
	ErrStoppingTestNotImplemented                          = errors.New("stopping test not implemented")
	ErrAbortTestNotImplemented                            = errors.New("abort test not implemented")
	ErrFailedToGetHealthCheck                             = errors.New("failed to get health check of test service")
	ErrFailedToGetMetrics                                 = errors.New("failed to get metric of test service")
	ErrFailedToGetLiveCheck                               = errors.New("failed to get live check")
	ErrRunAlreadyInProgress                               = errors.New("test service is already in progress")
	ErrFailedToRunTestService                             = errors.New("failed to execute run function of test service")
	ErrFailedToProvisionTestService                       = errors.New("failed to provision test service")
	ErrFailedToRunAgentControllerDueToProvisioningFailure = errors.New("unable to run test agent controller due to provisioning test service failure")
	ErrNumStepsNotSet                                     = errors.New("num steps not set")
	ErrNumStepsShouldBeOne                                = errors.New("num steps should be one")
	ErrIncreaseAgentNumNotBeNegative                      = errors.New("increase agent number couldn't be negative")
	ErrExecNumMultiAgentShouldBePositive                  = errors.New("execution number scenario in multi agent should be positive")
	ErrIncreaseFixedInputNotBeNegative                    = errors.New("increase fixed number couldn't be negative")
	ErrExecNumMultiFixedInputShouldBePositive             = errors.New("execution number of multi fixed input should be positive")
	ErrScenarioIsNotRunning                               = errors.New("scenario is not running")
	ErrMultiFixedInputConfigNotTrue                       = errors.New("both or none of execution number of multi fixed input and increase fixed input should be zero")
	ErrMultiAgentConfigNotTrue                            = errors.New("both or none of execution number of multi agent and increase agent should be zero")
	ErrInvalidDatabaseConfig                              = errors.New("invalid config applied to postgres database initializer")
)
