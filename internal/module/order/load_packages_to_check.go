package module

import (
	"context"
	"github.com/opentracing/opentracing-go"

	"HomeWork_1/internal/model"
)

func (m *Module) LoadPackagesToCheck(ctx context.Context) ([]model.Package, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "module.LoadPackagesToCheck")
	defer span.Finish()

	packages, err := m.storage.LoadPackages(ctx)
	if err != nil {
		return nil, err
	}

	return packages, nil
}
