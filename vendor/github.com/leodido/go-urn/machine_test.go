package urn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	for ii, tt := range tests {
		urn, err := NewMachine().Parse([]byte(tt.in))
		ok := err == nil

		if ok {
			assert.True(t, tt.ok, herror(ii, tt))
			assert.Equal(t, tt.obj.prefix, urn.prefix, herror(ii, tt))
			assert.Equal(t, tt.obj.ID, urn.ID, herror(ii, tt))
			assert.Equal(t, tt.obj.SS, urn.SS, herror(ii, tt))
			assert.Equal(t, tt.str, urn.String(), herror(ii, tt))
			assert.Equal(t, tt.norm, urn.Normalize().String(), herror(ii, tt))
		} else {
			assert.False(t, tt.ok, herror(ii, tt))
			assert.Empty(t, urn, herror(ii, tt))
			assert.Equal(t, tt.estr, err.Error())
		}
	}
}
