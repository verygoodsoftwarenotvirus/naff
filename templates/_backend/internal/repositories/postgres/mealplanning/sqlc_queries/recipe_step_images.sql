-- name: CreateRecipeStepImage :exec
INSERT INTO recipe_step_images (
	id,
	belongs_to_recipe_step,
	uploaded_media_id,
	uploaded_by_user
) VALUES (
	sqlc.arg(id),
	sqlc.arg(belongs_to_recipe_step),
	sqlc.arg(uploaded_media_id),
	sqlc.arg(uploaded_by_user)
);

-- name: GetRecipeStepImagesByStep :many
SELECT
	recipe_step_images.id,
	recipe_step_images.belongs_to_recipe_step,
	recipe_step_images.uploaded_media_id,
	recipe_step_images.uploaded_by_user,
	recipe_step_images.created_at,
	recipe_step_images.archived_at
FROM recipe_step_images
WHERE recipe_step_images.belongs_to_recipe_step = sqlc.arg(belongs_to_recipe_step)
	AND recipe_step_images.archived_at IS NULL
ORDER BY recipe_step_images.created_at;
