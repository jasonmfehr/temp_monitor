package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	startCallCount := 0

	lambdaStart = func(handler interface{}) {
		startCallCount++
	}

	main()

	assert.Equal(t, 1, startCallCount)
}
