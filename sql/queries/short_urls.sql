-- name: CreateShortURL :one
INSERT INTO short_urls (user_id, slug, original_url)
VALUES (@user_id :: text, @slug :: text, @original_url :: text) ON CONFLICT (original_url) DO
UPDATE
SET original_url = short_urls.original_url RETURNING slug;

-- name: BatchCreateShortURLs :many
INSERT INTO short_urls (user_id, slug, original_url)
SELECT unnest(@user_ids::text[]) AS user_id,
       unnest(@slugs::text[]) AS slug,
       unnest(@original_urls::text[]) AS original_url ON CONFLICT (slug) DO
UPDATE
SET slug = short_urls.slug RETURNING slug;

-- name: GetOriginalURLBySlug :one
SELECT original_url,
       deleted
FROM short_urls
WHERE slug = @slug :: text;

-- name: ListShortURLsByUser :many
SELECT slug,
       original_url
FROM short_urls
WHERE user_id = @user_id :: text
  AND deleted = FALSE;

-- name: BatchUpdateShortURLStatus :many
UPDATE short_urls
SET deleted = @deleted :: bool
WHERE user_id = @user_id :: text
  AND slug = ANY(@slugs :: text[])
RETURNING slug;