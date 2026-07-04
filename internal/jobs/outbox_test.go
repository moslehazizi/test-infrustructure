package jobs

import (
	"context"
	"control-panel-service/internal/usecase/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestOutboxJob_Run(t *testing.T) {
	t.Run("success_case_processes_on_every_tick", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 220*time.Millisecond)
		defer cancel()

		mockProcessor := new(mocks.MockOutboxProcessor)
		job := &outboxJob{mockProcessor}

		mockProcessor.On("ProcessPending", mock.Anything).Return(nil)

		job.Run(ctx, 50*time.Millisecond)

		mockProcessor.AssertCalled(t, "ProcessPending", mock.Anything)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("failed_case_process_pending_error_does_not_stop_the_loop", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 220*time.Millisecond)
		defer cancel()

		mockProcessor := new(mocks.MockOutboxProcessor)
		job := &outboxJob{mockProcessor}

		mockProcessor.On("ProcessPending", mock.Anything).Return(errors.New("claim failed"))

		job.Run(ctx, 50*time.Millisecond)

		mockProcessor.AssertCalled(t, "ProcessPending", mock.Anything)
		mockProcessor.AssertExpectations(t)
	})

	t.Run("success_case_stops_immediately_when_context_already_cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		mockProcessor := new(mocks.MockOutboxProcessor)
		job := &outboxJob{mockProcessor}

		job.Run(ctx, 50*time.Millisecond)

		mockProcessor.AssertNotCalled(t, "ProcessPending", mock.Anything)
	})
}
