package services

import (
	"context"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"projekat/model"
	"projekat/repositories"
)

type IdempotencyService struct {
	repo   repositories.ConfigConsulRepository
	Tracer trace.Tracer
}

func NewIdempotencyService(repo repositories.ConfigConsulRepository, tracer trace.Tracer) IdempotencyService {
	return IdempotencyService{
		repo:   repo,
		Tracer: tracer,
	}
}

func (i IdempotencyService) Add(req *model.IdempotencyRequest, ctx context.Context) error {
	ctx, span := i.Tracer.Start(ctx, "IdempotencyService.Add")
	_, err := i.repo.AddIdempotencyRequest(req, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetStatus(codes.Ok, "Service-Ok")
	return nil
}

func (i IdempotencyService) Get(key string, ctx context.Context) (bool, error) {
	ctx, span := i.Tracer.Start(ctx, "IdempotencyService.Get")
	defer span.End()

	exists, err := i.repo.GetIdempotencyRequestByKey(key, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}

	span.SetStatus(codes.Ok, "Service-Ok")
	return exists, nil
}
