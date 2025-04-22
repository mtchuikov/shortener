-- name: InsertShortenURL :one
INSERT INTO shorten_urls (shorten_id, original_url)
VALUES (@shorten_id::varchar, @original_url::varchar)
ON CONFLICT (original_url) DO UPDATE
  SET shorten_id = shorten_urls.shorten_id
RETURNING shorten_id;

-- name: GetOriginalURL :one
SELECT original_url 
FROM shorten_urls 
WHERE shorten_id = @shorten_id::varchar;