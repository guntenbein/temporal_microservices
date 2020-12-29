package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatchCount(t *testing.T) {
	assert.EqualValues(t, 1, batchCount(3, 10))
	assert.EqualValues(t, 10, batchCount(99, 10))
	assert.EqualValues(t, 0, batchCount(0, 10))
	assert.EqualValues(t, 0, batchCount(99, 0))
}

func TestBatch(t *testing.T) {
	pp, err := makeParallelepipeds(6, true)
	assert.NoError(t, err)

	outFigures := makeBatch(pp, 2, 2)
	assert.EqualValues(t, 2, len(outFigures))

	outFigures = makeBatch(pp, 1, 10)
	assert.EqualValues(t, 0, len(outFigures))

	outFigures = makeBatch(pp, -2, 3)
	assert.EqualValues(t, 0, len(outFigures))

	outFigures = makeBatch(pp, 2, 0)
	assert.EqualValues(t, 0, len(outFigures))

	outFigures = makeBatch(pp, 1, 5)
	assert.EqualValues(t, 1, len(outFigures))
}
