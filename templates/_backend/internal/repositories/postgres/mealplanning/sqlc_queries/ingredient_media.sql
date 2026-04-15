-- name: CreateIngredientMedia :exec
INSERT INTO ingredient_media (
	id,
	valid_ingredient_id,
	uploaded_media_id,
	index
) VALUES (
	sqlc.arg(id),
	sqlc.arg(valid_ingredient_id),
	sqlc.arg(uploaded_media_id),
	sqlc.arg(index)
);

-- name: GetIngredientMediaByIngredient :many
SELECT
	ingredient_media.id,
	ingredient_media.valid_ingredient_id,
	ingredient_media.uploaded_media_id,
	ingredient_media.index,
	ingredient_media.created_at,
	ingredient_media.archived_at
FROM ingredient_media
WHERE ingredient_media.valid_ingredient_id = sqlc.arg(valid_ingredient_id)
	AND ingredient_media.archived_at IS NULL
ORDER BY ingredient_media.index, ingredient_media.id;
