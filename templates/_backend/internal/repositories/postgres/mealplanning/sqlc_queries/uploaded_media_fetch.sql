-- name: GetUploadedMediaWithIDs :many
SELECT
	uploaded_media.id,
	uploaded_media.storage_path,
	uploaded_media.mime_type,
	uploaded_media.created_at,
	uploaded_media.last_updated_at,
	uploaded_media.archived_at,
	uploaded_media.created_by_user
FROM uploaded_media
WHERE uploaded_media.archived_at IS NULL
	AND uploaded_media.id = ANY(sqlc.arg(ids)::text[]);
