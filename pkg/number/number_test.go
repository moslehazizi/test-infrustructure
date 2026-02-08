package number

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcRange(t *testing.T) {
	t.Run("const is nil , min and max are not nil", func(t *testing.T) {
		var (
			cntPtr      *int
			minVal      = 10
			maxVal      = 100
			minPtr      = &minVal
			maxPtr      = &maxVal
			expectedMin = 10
			expectedMax = 100
		)

		actualMin, actualMax := CalcRange(cntPtr, minPtr, maxPtr)

		assert.Equal(t, expectedMin, actualMin)
		assert.Equal(t, expectedMax, actualMax)
	})

	t.Run("const is not nil min and max are nil", func(t *testing.T) {
		var (
			cntVal      = 10
			cntPtr      = &cntVal
			minPtr      *int
			maxPtr      *int
			expectedMin = 10
			expectedMax = 10
		)

		actualMin, actualMax := CalcRange(cntPtr, minPtr, maxPtr)

		assert.Equal(t, expectedMin, actualMin)
		assert.Equal(t, expectedMax, actualMax)
	})

	t.Run("const is nil , min is nil and max is not nil", func(t *testing.T) {
		var (
			cntPtr      *int
			minPtr      *int
			maxVal      = 100
			maxPtr      = &maxVal
			expectedMin = 0
			expectedMax = 100
		)

		actualMin, actualMax := CalcRange(cntPtr, minPtr, maxPtr)

		assert.Equal(t, expectedMin, actualMin)
		assert.Equal(t, expectedMax, actualMax)
	})

	t.Run("const is nil , min is not nil and max is nil", func(t *testing.T) {
		var (
			cntPtr      *int
			minVal      = 10
			minPtr      = &minVal
			maxPtr      *int
			expectedMin = 0
			expectedMax = 10
		)

		actualMin, actualMax := CalcRange(cntPtr, minPtr, maxPtr)

		assert.Equal(t, expectedMin, actualMin)
		assert.Equal(t, expectedMax, actualMax)
	})

	t.Run("const is not nil min is not nil and max is nil", func(t *testing.T) {
		var (
			cntVal      = 10
			cntPtr      = &cntVal
			minVal      = 10
			minPtr      = &minVal
			maxPtr      *int
			expectedMin = 10
			expectedMax = 10
		)

		actualMin, actualMax := CalcRange(cntPtr, minPtr, maxPtr)

		assert.Equal(t, expectedMin, actualMin)
		assert.Equal(t, expectedMax, actualMax)
	})

	t.Run("const is not nil min is nil and max is not nil", func(t *testing.T) {
		var (
			cntVal      = 10
			cntPtr      = &cntVal
			minPtr      *int
			maxVal      = 100
			maxPtr      = &maxVal
			expectedMin = 10
			expectedMax = 10
		)

		actualMin, actualMax := CalcRange(cntPtr, minPtr, maxPtr)

		assert.Equal(t, expectedMin, actualMin)
		assert.Equal(t, expectedMax, actualMax)
	})

	t.Run("const , min . max are nil", func(t *testing.T) {
		var (
			cntPtr      *int
			minPtr      *int
			maxPtr      *int
			expectedMin = 0
			expectedMax = 0
		)

		actualMin, actualMax := CalcRange(cntPtr, minPtr, maxPtr)

		assert.Equal(t, expectedMin, actualMin)
		assert.Equal(t, expectedMax, actualMax)
	})

	t.Run("const , min . max are not nil", func(t *testing.T) {
		var (
			cntVal      = 5
			cntPtr      = &cntVal
			minVal      = 10
			minPtr      = &minVal
			maxVal      = 100
			maxPtr      = &maxVal
			expectedMin = 5
			expectedMax = 5
		)

		actualMin, actualMax := CalcRange(cntPtr, minPtr, maxPtr)

		assert.Equal(t, expectedMin, actualMin)
		assert.Equal(t, expectedMax, actualMax)
	})

}
