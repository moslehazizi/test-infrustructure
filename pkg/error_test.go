package pkg

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToHTTPError(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.err.Error(), func(t *testing.T) {
			got := ToHTTPError(tt.err)
			assert.Equal(t, got.msg, tt.wanted.msg)
			assert.Equal(t, got.status, tt.wanted.status)
		})
	}
}
