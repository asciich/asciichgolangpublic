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

func IsUuid(input string) bool {
	_, err := uuid.Parse(input)
	return err == nil 
}
