package uuidutils

import (
	"context"

	"github.com/google/uuid"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

func Generate(ctx context.Context) string {
	uuid := uuid.New().String()

	logging.LogInfoByCtxf(ctx, "Generated new UUID: '%s'.", uuid)

	return uuid
}
