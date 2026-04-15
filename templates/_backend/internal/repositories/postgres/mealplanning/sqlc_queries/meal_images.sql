-- name: CreateMealImage :exec
INSERT INTO meal_images (
	id,
	belongs_to_meal,
	uploaded_media_id,
	uploaded_by_user
) VALUES (
	sqlc.arg(id),
	sqlc.arg(belongs_to_meal),
	sqlc.arg(uploaded_media_id),
	sqlc.arg(uploaded_by_user)
);
