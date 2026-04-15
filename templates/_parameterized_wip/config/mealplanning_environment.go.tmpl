package config

// This file contains environment config rendering for meal planning domain workers.
// Domain: mealplanning — remove this file when swapping the domain.

import (
	"path"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	mpfConfigObservabilityServiceName   = "meal_plan_finalizer"
	mpgliConfigObservabilityServiceName = "meal_plan_grocery_list_initializer"
	mptcConfigObservabilityServiceName  = "meal_plan_task_creator"
)

// renderMealPlanningConfigs creates, configures, and returns the meal planning worker configs
// for environment rendering. Returns configs for validation and a file path -> content map for writing.
func (s *EnvironmentConfigSet) renderMealPlanningConfigs(outputDir string, pretty bool) (configs []validation.ValidatableWithContext, files map[string][]byte) {
	mpfConfig := &MealPlanFinalizerConfig{
		Observability: s.RootConfig.Observability,
		Events:        s.RootConfig.Events,
		Database:      databaseConfigForService(&s.RootConfig.Database, s.ServiceDatabaseUsers, mpfConfigObservabilityServiceName),
		Queues:        s.RootConfig.Queues,
	}
	mpfConfig.Observability.Tracing.ServiceName = mpfConfigObservabilityServiceName
	mpfConfig.Observability.Metrics.ServiceName = mpfConfigObservabilityServiceName
	mpfConfig.Observability.Logging.ServiceName = mpfConfigObservabilityServiceName
	mpfConfig.Observability.Profiling.ServiceName = mpfConfigObservabilityServiceName
	disableWorkerOtelMetrics(&mpfConfig.Observability)

	mpgliConfig := &MealPlanGroceryListInitializerConfig{
		Observability: s.RootConfig.Observability,
		Analytics:     s.RootConfig.Analytics,
		Events:        s.RootConfig.Events,
		Database:      databaseConfigForService(&s.RootConfig.Database, s.ServiceDatabaseUsers, mpgliConfigObservabilityServiceName),
		Queues:        s.RootConfig.Queues,
	}
	mpgliConfig.Observability.Tracing.ServiceName = mpgliConfigObservabilityServiceName
	mpgliConfig.Observability.Metrics.ServiceName = mpgliConfigObservabilityServiceName
	mpgliConfig.Observability.Logging.ServiceName = mpgliConfigObservabilityServiceName
	mpgliConfig.Observability.Profiling.ServiceName = mpgliConfigObservabilityServiceName
	disableWorkerOtelMetrics(&mpgliConfig.Observability)

	mptcConfig := &MealPlanTaskCreatorConfig{
		Observability: s.RootConfig.Observability,
		Analytics:     s.RootConfig.Analytics,
		Events:        s.RootConfig.Events,
		Database:      databaseConfigForService(&s.RootConfig.Database, s.ServiceDatabaseUsers, mptcConfigObservabilityServiceName),
		Queues:        s.RootConfig.Queues,
	}
	mptcConfig.Observability.Tracing.ServiceName = mptcConfigObservabilityServiceName
	mptcConfig.Observability.Metrics.ServiceName = mptcConfigObservabilityServiceName
	mptcConfig.Observability.Logging.ServiceName = mptcConfigObservabilityServiceName
	mptcConfig.Observability.Profiling.ServiceName = mptcConfigObservabilityServiceName
	disableWorkerOtelMetrics(&mptcConfig.Observability)

	configs = []validation.ValidatableWithContext{mpfConfig, mpgliConfig, mptcConfig}

	files = map[string][]byte{
		path.Join(outputDir, stringOrDefault(s.MealPlanFinalizerConfigPath, "job_meal_plan_finalizer_config.json")):                             renderJSON(mpfConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanGroceryListInitializerConfigPath, "job_meal_plan_grocery_list_initializer_config.json")): renderJSON(mpgliConfig, pretty),
		path.Join(outputDir, stringOrDefault(s.MealPlanTaskCreatorConfigPath, "job_meal_plan_task_creator_config.json")):                        renderJSON(mptcConfig, pretty),
	}

	return configs, files
}
