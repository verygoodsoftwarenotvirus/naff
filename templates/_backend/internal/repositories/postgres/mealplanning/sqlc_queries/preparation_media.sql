-- name: CreatePreparationMedia :exec
INSERT INTO preparation_media (
	id,
	valid_preparation_id,
	for_ingredient_id,
	uploaded_media_id,
	index
) VALUES (
	sqlc.arg(id),
	sqlc.arg(valid_preparation_id),
	sqlc.arg(for_ingredient_id),
	sqlc.arg(uploaded_media_id),
	sqlc.arg(index)
);

-- name: GetPreparationMediaByPreparation :many
SELECT
	preparation_media.id,
	preparation_media.valid_preparation_id,
	preparation_media.for_ingredient_id,
	preparation_media.uploaded_media_id,
	preparation_media.index,
	preparation_media.created_at,
	preparation_media.archived_at
FROM preparation_media
WHERE preparation_media.valid_preparation_id = sqlc.arg(valid_preparation_id)
	AND preparation_media.archived_at IS NULL
ORDER BY preparation_media.index, preparation_media.id;

-- name: GetPreparationMediaByPreparationAndIngredient :many
-- Returns media for a prep+ingredient: ingredient-specific media OR general (for_ingredient_id IS NULL)
SELECT
	preparation_media.id,
	preparation_media.valid_preparation_id,
	preparation_media.for_ingredient_id,
	preparation_media.uploaded_media_id,
	preparation_media.index,
	preparation_media.created_at,
	preparation_media.archived_at
FROM preparation_media
WHERE preparation_media.valid_preparation_id = sqlc.arg(valid_preparation_id)
	AND (preparation_media.for_ingredient_id = sqlc.arg(for_ingredient_id) OR preparation_media.for_ingredient_id IS NULL)
	AND preparation_media.archived_at IS NULL
ORDER BY preparation_media.index, preparation_media.id;
