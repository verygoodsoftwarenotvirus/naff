-- name: GetMealPlanTaskNotificationContext :one
SELECT
	recipe_prep_tasks.name as prep_task_name,
	meal_plan_tasks.creation_explanation,
	meal_plan_events.meal_name,
	meal_plan_events.starts_at
FROM meal_plan_tasks
	JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option = meal_plan_options.id
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event = meal_plan_events.id
	JOIN recipe_prep_tasks ON meal_plan_tasks.belongs_to_recipe_prep_task = recipe_prep_tasks.id
WHERE meal_plan_options.archived_at IS NULL
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_tasks.id = sqlc.arg(meal_plan_task_id);
