-- name: CreateRootNode :one
INSERT INTO org_nodes (id, tenant_id, parent_id, type, path, name, status_id, created_by, created_at, updated_at)
VALUES (
  sqlc.arg(id)::uuid,
  sqlc.arg(id)::uuid,
  NULL,
  'tenant',
  text2ltree(sqlc.arg(id)::text),
  sqlc.arg(name)::text,
  sqlc.arg(status_id)::uuid,
  sqlc.arg(created_by)::uuid,
  now(),
  now()
)
RETURNING *;

-- name: AppendNode :one
INSERT INTO org_nodes (id, tenant_id, parent_id, type, path, name, status_id, created_by, created_at, updated_at)
SELECT
  sqlc.arg(id)::uuid,
  p.tenant_id,
  p.id,
  sqlc.arg(type)::text,
  p.path || text2ltree(sqlc.arg(id)::text),
  sqlc.arg(name)::text,
  sqlc.arg(status_id)::uuid,
  sqlc.arg(created_by)::uuid,
  now(),
  now()
FROM org_nodes p
WHERE p.id = sqlc.arg(parent_id)::uuid
  AND p.deleted_at IS NULL
  AND (
    (p.type = 'tenant' AND sqlc.arg(type)::text = 'workspace')
    OR (p.type = 'workspace' AND sqlc.arg(type)::text = 'project')
  )
RETURNING *;

-- name: GetAncestors :many
SELECT a.*
FROM org_nodes n
JOIN org_nodes a ON a.path @> n.path AND a.id <> n.id
WHERE n.id = sqlc.arg(id)::uuid
ORDER BY nlevel(a.path);

-- name: GetDescendants :many
SELECT d.*
FROM org_nodes n
JOIN org_nodes d ON d.path <@ n.path AND d.id <> n.id
WHERE n.id = sqlc.arg(id)::uuid
  AND d.deleted_at IS NULL
  AND (sqlc.narg(type)::text IS NULL OR d.type = sqlc.narg(type)::text)
ORDER BY d.path;
