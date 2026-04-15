package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/pflag"
)

var (
	checkOnlyFlag = pflag.Bool("check", false, "only check if files match")
	databaseFlag  = pflag.String("database", postgres, "what database to use")
)

const (
	postgres = "postgres"
)

func main() {
	pflag.Parse()

	runErrors := &multierror.Error{}
	databaseToUse := *databaseFlag

	queryOutput := map[string][]*Query{
		"internal/repositories/postgres/internalops/sqlc_queries/internalops":                                   buildMaintenanceQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_ingredients":                            buildValidIngredientsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_instruments":                            buildValidInstrumentsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_preparations":                           buildValidPreparationsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_measurement_units":                      buildValidMeasurementUnitsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_ingredient_states":                      buildValidIngredientStatesQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_vessels":                                buildValidVesselsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_ingredient_groups":                      buildValidIngredientGroupsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_ingredient_preparations":                buildValidIngredientPreparationsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_prep_task_configs":                      buildValidPrepTaskConfigsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_preparation_vessels":                    buildValidPreparationVesselsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_ingredient_measurement_units":           buildValidIngredientMeasurementUnitsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_measurement_unit_conversions":           buildValidMeasurementUnitConversionsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_ingredient_state_ingredients":           buildValidIngredientStateIngredientsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/valid_preparation_instruments":                buildValidPreparationInstrumentsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/account_instrument_ownerships":                buildAccountInstrumentOwnershipQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meal_components":                              buildMealComponentsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meal_plan_events":                             buildMealPlanEventsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_media":                                 buildRecipeMediaQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_prep_task_steps":                       buildRecipePrepTaskStepsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_ratings":                               buildRecipeRatingsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_step_completion_condition_ingredients": buildRecipeStepCompletionConditionIngredientsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_prep_tasks":                            buildRecipePrepTasksQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meals":                                        buildMealsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meal_plans":                                   buildMealPlansQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_step_completion_conditions":            buildRecipeStepCompletionConditionQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meal_plan_option_votes":                       buildMealPlanOptionVotesQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meal_plan_options":                            buildMealPlanOptionsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meal_plan_tasks":                              buildMealPlanTasksQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipes":                                      buildRecipesQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_lists":                                 buildRecipeListsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meal_lists":                                   buildMealListsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_step_ingredients":                      buildRecipeStepIngredientsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_step_instruments":                      buildRecipeStepInstrumentsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_step_products":                         buildRecipeStepProductsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_steps":                                 buildRecipeStepsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_step_vessels":                          buildRecipeStepVesselsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/user_ingredient_preferences":                  buildUserIngredientPreferencesQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meal_plan_grocery_list_items":                 buildMealPlanGroceryListItemsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meal_plan_recipe_option_selections":           buildMealPlanRecipeOptionSelectionsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/meal_list_items":                              buildMealListItemsQueries(databaseToUse),
		"internal/repositories/postgres/mealplanning/sqlc_queries/recipe_list_items":                            buildRecipeListItemsQueries(databaseToUse),
		"internal/repositories/postgres/oauth/sqlc_queries/oauth2_client_tokens":                                buildOAuth2ClientTokensQueries(databaseToUse),
		"internal/repositories/postgres/oauth/sqlc_queries/oauth2_clients":                                      buildOAuth2ClientsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/account_invitations":                              buildAccountInvitationsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/account_user_memberships":                         buildAccountUserMembershipsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/accounts":                                         buildAccountsQueries(databaseToUse),
		"internal/repositories/postgres/auditlogentries/sqlc_queries/audit_logs":                                buildAuditLogEntryQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/admin":                                            buildAdminQueries(databaseToUse),
		"internal/repositories/postgres/auth/sqlc_queries/password_reset_tokens":                                buildPasswordResetTokensQueries(databaseToUse),
		"internal/repositories/postgres/auth/sqlc_queries/user_sessions":                                        buildUserSessionsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/users":                                            buildUsersQueries(databaseToUse),
		"internal/repositories/postgres/settings/sqlc_queries/service_settings":                                 buildServiceSettingQueries(databaseToUse),
		"internal/repositories/postgres/settings/sqlc_queries/service_setting_configurations":                   buildServiceSettingConfigurationQueries(databaseToUse),
		"internal/repositories/postgres/webhooks/sqlc_queries/webhooks":                                         buildWebhooksQueries(databaseToUse),
		"internal/repositories/postgres/webhooks/sqlc_queries/webhook_trigger_events":                           buildWebhookTriggerEventsQueries(databaseToUse),
		"internal/repositories/postgres/webhooks/sqlc_queries/webhook_trigger_configs":                          buildWebhookTriggerConfigsQueries(databaseToUse),
		"internal/repositories/postgres/notifications/sqlc_queries/user_notifications":                          buildUserNotificationQueries(databaseToUse),
		"internal/repositories/postgres/waitlists/sqlc_queries/waitlists":                                       buildWaitlistsQueries(databaseToUse),
		"internal/repositories/postgres/waitlists/sqlc_queries/waitlist_signups":                                buildWaitlistSignupsQueries(databaseToUse),
		"internal/repositories/postgres/issuereports/sqlc_queries/issue_reports":                                buildIssueReportsQueries(databaseToUse),
		"internal/repositories/postgres/uploadedmedia/sqlc_queries/uploaded_media":                              buildUploadedMediaQueries(databaseToUse),
		"internal/repositories/postgres/dataprivacy/sqlc_queries/user_data_disclosures":                         buildUserDataDisclosuresQueries(databaseToUse),
		"internal/repositories/postgres/payments/sqlc_queries/products":                                         buildPaymentsProductsQueries(databaseToUse),
		"internal/repositories/postgres/payments/sqlc_queries/subscriptions":                                    buildPaymentsSubscriptionsQueries(databaseToUse),
		"internal/repositories/postgres/payments/sqlc_queries/purchases":                                        buildPaymentsPurchasesQueries(databaseToUse),
		"internal/repositories/postgres/payments/sqlc_queries/payment_transactions":                             buildPaymentsTransactionsQueries(databaseToUse),
		"internal/repositories/postgres/comments/sqlc_queries/comments":                                         buildCommentsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/user_roles":                                       buildUserRolesQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/permissions":                                      buildPermissionsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/user_role_permissions":                            buildUserRolePermissionsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/user_role_assignments":                            buildUserRoleAssignmentsQueries(databaseToUse),
		"internal/repositories/postgres/identity/sqlc_queries/user_role_hierarchy":                              buildUserRoleHierarchyQueries(databaseToUse),
	}

	checkOnly := *checkOnlyFlag

	for filePath, queries := range queryOutput {
		filePath += ".generated.sql"
		existingFile, err := os.ReadFile(filePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if _, err = os.Create(filePath); err != nil {
					log.Fatal(fmt.Errorf("creating file: %w", err))
				}
			}
			if err != nil {
				log.Fatal(fmt.Errorf("opening existing file: %w", err))
			}
		}

		var fileContent strings.Builder
		for i, query := range queries {
			if i != 0 {
				fileContent.WriteString("\n")
			}
			fileContent.WriteString(query.Render())
		}

		var fileOutput strings.Builder
		for line := range strings.SplitSeq(strings.TrimSpace(fileContent.String()), "\n") {
			fileOutput.WriteString(strings.TrimSuffix(line, " ") + "\n")
		}

		if string(existingFile) != fileOutput.String() && checkOnly {
			runErrors = multierror.Append(runErrors, fmt.Errorf("files don't match: %s", filePath))
		}

		if !checkOnly {
			if err = os.WriteFile(filePath, []byte(fileOutput.String()), 0o600); err != nil {
				runErrors = multierror.Append(runErrors, err)
			}
		}
	}

	if runErrors.ErrorOrNil() != nil {
		log.Fatal(runErrors)
	}
}
