-- name: VerifyProjectHierarchy :one
SELECT EXISTS (
  SELECT 1
  FROM tenants t
  LEFT JOIN workspaces w
    ON w.tenant_id = t.id AND w.id = sqlc.narg(workspace_id)::uuid
  LEFT JOIN projects p
    ON p.workspace_id = w.id AND p.id = sqlc.narg(project_id)::uuid
  WHERE t.id = sqlc.arg(tenant_id)::uuid
    AND (sqlc.narg(workspace_id)::uuid IS NULL OR w.id IS NOT NULL)
    AND (sqlc.narg(project_id)::uuid IS NULL OR p.id IS NOT NULL)
);

-- name: VerifyPermission :one
SELECT EXISTS (
  SELECT 1
  FROM tenant_member tm
  
);