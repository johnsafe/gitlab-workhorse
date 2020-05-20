package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHoversRead(t *testing.T) {
	h := setupHovers(t)

	require.Equal(t, `[{"value":"hello"}]`, string(h.For("1")))

	require.NoError(t, os.Remove(h.File.Name()))
}

func setupHovers(t *testing.T) *Hovers {
	h, err := NewHovers("")
	require.NoError(t, err)

	require.NoError(t, h.Read("hoverResult", []byte(`{"id":"2","label":"hoverResult","result":{"contents": ["hello"]}}`)))
	require.NoError(t, h.Read("textDocument/hover", []byte(`{"id":"4","label":"textDocument/hover","outV":"3","inV":"2"}`)))
	require.NoError(t, h.Read("textDocument/references", []byte(`{"id":"3","label":"textDocument/references","outV":"3","inV":"1"}`)))

	return h
}
