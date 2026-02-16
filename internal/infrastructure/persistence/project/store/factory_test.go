package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoProjectreFactory(t *testing.T) {

	var factory *Factory
	var storer *LocalStorage

	factory = NewFactory()

	storer = NewLocalStorage(nil, "", nil)
	factory.Register("storer", storer)

	t.Run("Get storer", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing get storer")
		f := factory.Get("storer")
		assert.NotNil(t, f)
		assert.Equal(t, storer, f)
	})

	t.Run("Get storer not registered", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing get storer not registered")
		f := factory.Get("storer_not_registered")
		assert.Nil(t, f)
	})
}
