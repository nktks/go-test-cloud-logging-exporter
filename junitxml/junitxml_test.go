package junitxml

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/junit.xml
var b []byte

func TestParse(t *testing.T) {
	j, err := Unmarshal(b)
	require.NoError(t, err)
	require.NotNil(t, j)
}
func TestParse_Fail(t *testing.T) {
	s := "test"
	j, err := Unmarshal([]byte(s))
	require.Error(t, err)
	require.Nil(t, j)
}
