package repository

import (
	"testing"

	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository/local"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestProjectRepositoryFactory(t *testing.T) {

	var factory *Factory

	factory = NewFactory()

	database := local.NewDatabaseDriver(afero.NewMemMapFs(), "repository", logger.NewFakeLogger())
	factory.Register("local", database)

	t.Run("Get repository", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing get repository")
		f := factory.Get("local")
		assert.NotNil(t, f)
		assert.Equal(t, database, f)
	})

	t.Run("Get repository not registered", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing get repository not registered")
		f := factory.Get("repository_not_registered")
		assert.Nil(t, f)
	})
}
