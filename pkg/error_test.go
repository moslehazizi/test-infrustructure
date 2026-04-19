package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toHTTPError(t *testing.T) {
	tests := []struct {
		err    error
		wanted HTTPError
	}{
		// INTERNAL
		{
			err:    errors.New("unknown error"),
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTestScenariosByStatus,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrInternalServerError,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToSendEventData,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToConsumeData,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToUnmarshalEventData,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToLoadConfig,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToCreateMotherService,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetMotherService,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetMotherServices,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrNegativePageOrPerPageNotAllowed,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTestScenario,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTestScenarios,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToCreateTestScenario,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTestCategory,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTestCategoriesFromRepository,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTestCategoryFromRepository,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToSendProvisioningEvent,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToCreateMotherService,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTestCategoriesFromRepository,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTestCategoriesFromRepository,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToCreateMotherService,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTestCategoryFromRepository,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTablesOfDatabase,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetDatabases,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToUpdateTestScenario,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrStartingTestNotImplemented,
			wanted: HTTPError{http.StatusInternalServerError, StartingTestNotImplemented},
		},
		{
			err:    ErrFailedToGetHealthCheck,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetMetrics,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetLiveCheck,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToRunTestService,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetTestServiceConfig,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToGetReadyForTesting,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToPauseTestService,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToResumeTestService,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},
		{
			err:    ErrFailedToStopTestService,
			wanted: HTTPError{http.StatusInternalServerError, InternalServerErrorMessage},
		},

		// Other
		{
			err:    ErrBadRequest,
			wanted: HTTPError{http.StatusBadRequest, InvalidReqBody},
		},
		{
			err:    ErrMotherServiceNotFound,
			wanted: HTTPError{http.StatusNotFound, MotherServiceNotFound},
		},
		{
			err:    ErrMotherServiceAlreadyExist,
			wanted: HTTPError{http.StatusConflict, MotherServiceAlreadyExist},
		},
		{
			err:    ErrTestScenarioNotFound,
			wanted: HTTPError{http.StatusNotFound, TestScenarioNotFound},
		},
		{
			err:    ErrTestCategoryNotFound,
			wanted: HTTPError{http.StatusNotFound, TestCategoryNotFound},
		},
		{
			err:    ErrInvalidIDInParams,
			wanted: HTTPError{http.StatusBadRequest, InvalidIDInParams},
		},
		{
			err:    ErrPageNotFound,
			wanted: HTTPError{http.StatusNotFound, PageNotFound},
		},
		{
			err:    ErrRunAlreadyInProgress,
			wanted: HTTPError{http.StatusConflict, TestServiceAlreadyInProgress},
		},
		{
			err:    ErrTestServiceConfigNotFound,
			wanted: HTTPError{http.StatusNotFound, TestServiceConfigNotFound},
		},

		// 422 validation errors
		{
			err:    ErrInvalidMotherServiceName,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidMotherServiceName},
		},
		{
			err:    ErrInvalidExceptionRate,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidExceptionRate},
		},
		{
			err:    ErrInvalidResponseDelayRate,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidResponseDelayRate},
		},
		{
			err:    ErrInvalidDelayConfiguration,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidDelayConfiguration},
		},
		{
			err:    ErrInvalidRandomDelayRange,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidRandomDelayRange},
		},
		{
			err:    ErrInvalidDatabaseName,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidDatabaseName},
		},
		{
			err:    ErrInvalidDatabaseTableName,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidDatabaseTableName},
		},

		{
			err:    ErrExecutionDurationLessThanOne,
			wanted: HTTPError{http.StatusUnprocessableEntity, ExecutionDurationLessThanOne},
		},
		{
			err:    ErrMaxTestServiceCountLessThanOne,
			wanted: HTTPError{http.StatusUnprocessableEntity, MaxTestServiceCountLessThanOne},
		},
		{
			err:    ErrAutoStepChangeRateLessThanOne,
			wanted: HTTPError{http.StatusUnprocessableEntity, AutoStepChangeRateLessThanOne},
		},
		{
			err:    ErrMaxTestServiceCountNotSet,
			wanted: HTTPError{http.StatusUnprocessableEntity, MaxTestServiceCountNotSet},
		},
		{
			err:    ErrNoNeedMaxTestServiceCount,
			wanted: HTTPError{http.StatusUnprocessableEntity, NoNeedMaxTestServiceCount},
		},
		{
			err:    ErrExecutionDurationNotSet,
			wanted: HTTPError{http.StatusUnprocessableEntity, ExecutionDurationNotSet},
		},
		{
			err:    ErrNoNeedExecutionDuration,
			wanted: HTTPError{http.StatusUnprocessableEntity, NoNeedExecutionDuration},
		},
		{
			err:    ErrAutoStepChangeNotSet,
			wanted: HTTPError{http.StatusUnprocessableEntity, AutoStepChangeNotSet},
		},
		{
			err:    ErrNoNeedAutoStepChange,
			wanted: HTTPError{http.StatusUnprocessableEntity, NoNeedAutoStepChange},
		},
		{
			err:    ErrInvalidMaxRequest,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidMaxRequest},
		},
		{
			err:    ErrInvalidMaxDuration,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidMaxDuration},
		},
		{
			err:    ErrInvalidBadValueRate,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidBadValueRate},
		},
		{
			err:    ErrInvalidNegativeValueRate,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidNegativeValueRate},
		},
		{
			err:    ErrInvalidRealValueRate,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidRealValueRate},
		},
		{
			err:    ErrInvalidZeroValueRate,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidZeroValueRate},
		},
		{
			err:    ErrInvalidStringValueRate,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidStringValueRate},
		},
		{
			err:    ErrInvalidLongStringValueRate,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidLongStringValueRate},
		},
		{
			err:    ErrInvalidNullValueRate,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidNullValueRate},
		},
		{
			err:    ErrInvalidRequestDelayDurationConfig,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidRequestDelayDurationConfig},
		},
		{
			err:    ErrInvalidRequestDelayDuration,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidRequestDelayDuration},
		},
		{
			err:    ErrMinDelayDurationMoreThanMax,
			wanted: HTTPError{http.StatusUnprocessableEntity, MinDelayDurationMoreThanMax},
		},
		{
			err:    ErrInvalidTestNumberConfig,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidTestNumberConfig},
		},
		{
			err:    ErrInvalidFixedTestNumber,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidFixedTestNumber},
		},
		{
			err:    ErrInvalidFixedTestNumberConfig,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidFixedTestNumberConfig},
		},
		{
			err:    ErrInvalidMinOrMaxRandomTestNumber,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidMinOrMaxRandomTestNumber},
		},
		{
			err:    ErrMinRandomTestNumberMoreThanMax,
			wanted: HTTPError{http.StatusUnprocessableEntity, MinRandomTestNumberMoreThanMax},
		},
		{
			err:    ErrInvalidZeroSumOfBadValues,
			wanted: HTTPError{http.StatusUnprocessableEntity, InvalidZeroSumOfBadValues},
		},
		{
			err:    ErrInvalid100SumOfBadValues,
			wanted: HTTPError{http.StatusUnprocessableEntity, Invalid100SumOfBadValues},
		},
		{
			err:    ErrTestServiceConfigIsRequired,
			wanted: HTTPError{http.StatusUnprocessableEntity, TestServiceConfigIsRequired},
		},
		{
			err:    ErrOnlyReadyScenariosCanBeStarted,
			wanted: HTTPError{http.StatusUnprocessableEntity, OnlyReadyScenariosCanBeStarted},
		},
		{
			err:    ErrOnlyRunningScenariosCanBePaused,
			wanted: HTTPError{http.StatusUnprocessableEntity, OnlyRunningScenariosCanBePaused},
		},
		{
			err:    ErrOnlyPausedScenariosCanBeResume,
			wanted: HTTPError{http.StatusUnprocessableEntity, OnlyPausedScenariosCanBeReStarted},
		},
		{
			err:    ErrNumStepsShouldBeOne,
			wanted: HTTPError{http.StatusUnprocessableEntity, NumStepsShouldBeOne},
		},
		{
			err:    ErrIncreaseAgentNumNotBeNegative,
			wanted: HTTPError{http.StatusUnprocessableEntity, IncreaseAgentNumNotBeNegative},
		},
		{
			err:    ErrExecNumMultiAgentShouldBePositive,
			wanted: HTTPError{http.StatusUnprocessableEntity, ExecNumMultiAgentShouldBePositive},
		},
		{
			err:    ErrIncreaseFixedInputNotBeNegative,
			wanted: HTTPError{http.StatusUnprocessableEntity, IncreaseFixedNumberNotBeNegative},
		},
		{
			err:    ErrExecNumMultiFixedInputShouldBePositive,
			wanted: HTTPError{http.StatusUnprocessableEntity, ExecNumMultiFixedInputShouldBePositive},
		},
		{
			err:    ErrScenariosCanNotBeDelete,
			wanted: HTTPError{http.StatusUnprocessableEntity, ScenariosCanNotBeDelete},
		},
		{
			err:    ErrScenariosCanNotBeUpdated,
			wanted: HTTPError{http.StatusUnprocessableEntity, ScenariosCanNotBeUpdated},
		},
	}
	for _, tt := range tests {
		t.Run(tt.err.Error(), func(t *testing.T) {
			got := toHTTPError(tt.err)
			assert.Equal(t, got.msg, tt.wanted.msg)
			assert.Equal(t, got.status, tt.wanted.status)
		})
	}
}

func TestToHTTPError(t *testing.T) {
	t.Run("success_case_we_have_only_one_error_and_having_no_error_wrapping", func(t *testing.T) {
		want := &HTTPError{
			status: http.StatusInternalServerError,
			msg:    InternalServerErrorMessage,
		}
		got := ToHTTPError(ErrFailedToGetTestCategoryFromRepository)
		assert.Equal(t, want, got)
	})
	t.Run("success_case_we_have_2_wrapped_error_500_and_400_and_should_get_400", func(t *testing.T) {
		want := &HTTPError{
			status: http.StatusBadRequest,
			msg:    InvalidReqBody,
		}
		got := ToHTTPError(fmt.Errorf("%w:%w", ErrFailedToGetTestCategoryFromRepository, ErrBadRequest))
		assert.Equal(t, want, got)
	})
	t.Run("success_case_we_have_2_wrapped_error_400_and_500_and_should_get_400", func(t *testing.T) {
		want := &HTTPError{
			status: http.StatusBadRequest,
			msg:    InvalidReqBody,
		}
		got := ToHTTPError(fmt.Errorf("%w:%w", ErrBadRequest, ErrFailedToGetTestCategoryFromRepository))
		assert.Equal(t, want, got)
	})
	t.Run("success_case_we_have_3_wrapped_error_400,_422,_and_500_and_should_get_400", func(t *testing.T) {
		want := &HTTPError{
			status: http.StatusBadRequest,
			msg:    InvalidReqBody,
		}
		got := ToHTTPError(fmt.Errorf("%w:%w:%w", ErrBadRequest, ErrFailedToGetTestCategoryFromRepository, ErrTestServiceConfigIsRequired))
		assert.Equal(t, want, got)
	})
	t.Run("success_case_we_have_1_wrapped_error_400_and_should_get_400", func(t *testing.T) {
		want := &HTTPError{
			status: http.StatusBadRequest,
			msg:    InvalidReqBody,
		}
		got := ToHTTPError(fmt.Errorf("%w", ErrBadRequest))
		assert.Equal(t, want, got)
	})
	t.Run("success_case_we_do_not_have_a_wrapped_error", func(t *testing.T) {
		want := &HTTPError{
			status: http.StatusInternalServerError,
			msg:    InternalServerErrorMessage,
		}
		got := ToHTTPError(errors.New("something went wrong"))
		assert.Equal(t, want, got)
	})
}
