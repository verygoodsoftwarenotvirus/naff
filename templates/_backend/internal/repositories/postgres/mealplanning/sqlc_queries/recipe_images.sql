-- name: CreateRecipeImage :exec
INSERT INTO recipe_images (
	id,
	belongs_to_recipe,
	uploaded_media_id,
	uploaded_by_user
) VALUES (
	sqlc.arg(id),
	sqlc.arg(belongs_to_recipe),
	sqlc.arg(uploaded_media_id),
	sqlc.arg(uploaded_by_user)
);
