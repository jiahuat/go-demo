package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCreateClusterReq(t *testing.T) {
	r1 := &CreateClusterReq{}
	err := r1.Validate()
	assert.NotNil(t, err)

	r2 := &CreateClusterReq{Name: "name"}
	err = r2.Validate()
	assert.Nil(t, err)
}
