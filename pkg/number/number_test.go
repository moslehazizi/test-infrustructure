package number

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcRange(t *testing.T) {
	t.Run("const_is_nil_min_and_max_are_not_nil", func(t *testing.T) {
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

	t.Run("const_is_not_nil_min_and_max_are_nil", func(t *testing.T) {
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

	t.Run("const_is_nil_min_is_nil_and_max_is_not_nil", func(t *testing.T) {
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

	t.Run("const_is_nil_min_is_not_nil_and_max_is_nil", func(t *testing.T) {
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

	t.Run("const_is_not_nil_min_is_not_nil_and_max_is_nil", func(t *testing.T) {
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

	t.Run("const_is_not_nil_min_is_nil_and_max_is_not_nil", func(t *testing.T) {
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

	t.Run("const_min_max_are_nil", func(t *testing.T) {
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

	t.Run("const_min_max_are_not_nil", func(t *testing.T) {
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

func TestNumValue(t *testing.T) {
	t.Run("input_is_nil", func(t *testing.T) {
		var (
			inputPtr  *int
			excpected int = 0
		)

		actual := NumValue(inputPtr)
		assert.Equal(t, excpected, actual)
	})

	t.Run("input_is_not_nil", func(t *testing.T) {
		var (
			inputVal      = 10
			inputPtr      = &inputVal
			excpected int = 10
		)

		actual := NumValue(inputPtr)
		assert.Equal(t, excpected, actual)
	})
}
