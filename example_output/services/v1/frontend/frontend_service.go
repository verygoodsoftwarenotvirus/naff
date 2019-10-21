package frontend

import (
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/config"
)

var serviceName = "frontend_service"

type Service struct {
	logger logging.Logger
	config config.FrontendSettings
}

// ProvideFrontendService provides the frontend service to dependency injection
func ProvideFrontendService(logger logging.Logger, cfg config.FrontendSettings) *Service {
	svc := &Service{
		config: cfg,
		logger: logger.WithName(serviceName),
	}
	return svc
}
