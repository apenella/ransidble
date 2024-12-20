package unpack

import (
	"testing"

	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFetchFactory(t *testing.T) {

	var factory *Factory

	factory = NewFactory()

	unpacker := NewPlainFormat(afero.NewMemMapFs(), logger.NewFakeLogger())
	factory.Register("unpacker", unpacker)

	t.Run("Get unpacker", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing get fetcher")
		u := factory.Get("unpacker")
		assert.NotNil(t, u)
		assert.Equal(t, unpacker, u)
	})

	t.Run("Get unpacker not registered", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing get fetcher not registered")
		u := factory.Get("fetcher_not_registered")
		assert.Nil(t, u)
	})
}
