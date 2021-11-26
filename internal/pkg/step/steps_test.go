package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	steps := []Step{}
	err := Execute(steps)
	assert.NoError(t, err)

	counter := 0
	increment := func() error {
		counter++
		return nil
	}

	fail := func() error {
		return assert.AnError
	}

	steps = []Step{
		increment,
	}
	err = Execute(steps)
	assert.NoError(t, err)
	assert.Equal(t, 1, counter)

	counter = 0
	steps = []Step{
		increment,
		increment,
		increment,
	}
	err = Execute(steps)
	assert.NoError(t, err)
	assert.Equal(t, 3, counter)

	counter = 0
	steps = []Step{
		increment,
		fail,
		increment,
	}
	err = Execute(steps)
	assert.Error(t, err)
	assert.Equal(t, 1, counter)
}
