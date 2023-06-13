-- name: ListEndpointByOwner :many
SELECT ep.id AS endpoint_id,ep.name AS endpoint_name, ep.api_id, ep.template_name, 
       ak.key, a.name AS api_name, a.type, a.icon, a.image, a.document_url
  FROM endpoints AS ep
  INNER JOIN apikeys AS ak
    ON ep.api_id = ak.api_id
  INNER JOIN apis AS a
    ON ep.api_id = a.id
  WHERE ak.owner = $1
    AND ep.deleted_at IS NULL
    AND ak.deleted_at IS NULL
    AND a.deleted_at IS NULL
    AND a.type = 'source'
  ORDER BY a.name, ep.id;

-- name: ListEndpointByAPIID :many
SELECT name, api_id, template_name
  FROM endpoints
 WHERE api_id = ANY(@api_id::int[]) 
   AND deleted_at IS NULL;

-- name: CreateEndpoint :one
INSERT INTO endpoints (
    name, api_id, template_name
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: DeleteEndpoint :execrows
DELETE FROM endpoints
 WHERE id = $1;