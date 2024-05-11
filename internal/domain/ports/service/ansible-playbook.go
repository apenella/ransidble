package service

import (
	"context"

	"github.com/apenella/ransidble/internal/domain/core/entity"
)

type AnsiblePlaybookServicer interface {
	GenerateID() string
	Run(ctx context.Context, task *entity.Task) error
}
