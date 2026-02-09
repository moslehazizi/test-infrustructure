package number

const (
	Zero = 0
)

func CalcRange(cntPtr, minPtr, maxPtr *int) (minVal, maxVal int) {
	if cntPtr != nil {
		return *cntPtr, *cntPtr
	}

	minVal = Zero
	maxVal = Zero

	if minPtr != nil {
		minVal = *minPtr
	}
	if maxPtr != nil {
		maxVal = *maxPtr
	}

	if minVal > maxVal {
		return maxVal, minVal
	}

	return minVal, maxVal
}

func NumValue(input *int) int {
	if input == nil {
		return Zero
	}
	return *input
}
