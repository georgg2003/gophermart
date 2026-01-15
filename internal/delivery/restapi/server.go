package restapi

import (
	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/pkg/logging"
	"github.com/georgg2003/gophermart/internal/usecase"
)

//go:generate go tool oapi-codegen --config=./gen/server.yaml ../../../api/swagger.yaml
//go:generate go tool oapi-codegen --config=./gen/models.yaml ../../../api/swagger.yaml

type server struct {
	cfg    *config.Config
	logger *logging.Logger
	uc     usecase.UseCase
}

func NewServer(
	cfg *config.Config,
	logger *logging.Logger,
	uc usecase.UseCase,
) ServerInterface {
	return &server{
		cfg:    cfg,
		logger: logger,
		uc:     uc,
	}
}
