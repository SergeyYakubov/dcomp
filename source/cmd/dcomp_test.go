package main

import (
	//      "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	//      "utils"
)

func TestDummy(t *testing.T) {
	assert.Equal(t, Dummy(), 2, "This cannot be")
}
