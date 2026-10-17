INSERT INTO public.modules (id, "name", "key", "type", created_at, updated_at) VALUES
  (gen_random_uuid(), 'Tenant',           'tenant',           'tenant',  now(), now()),
  (gen_random_uuid(), 'Subscriptions',    'subscriptions',    'tenant',  now(), now()),
  (gen_random_uuid(), 'Workspaces',       'workspaces',       'tenant',  now(), now()),
  (gen_random_uuid(), 'Projects',         'projects',         'tenant',  now(), now()),
  (gen_random_uuid(), 'Users',            'users',            'tenant',  now(), now()),
  (gen_random_uuid(), 'Roles',            'roles',            'tenant',  now(), now()),
  (gen_random_uuid(), 'Audit Log',        'audit',            'tenant',  now(), now()),
  (gen_random_uuid(), 'Project Settings', 'project_settings', 'project', now(), now()),
  (gen_random_uuid(), 'Project Members',  'project_members',  'project', now(), now()),
  (gen_random_uuid(), 'Project Roles',    'project_roles',    'project', now(), now()),
  (gen_random_uuid(), 'Nodes',            'nodes',            'project', now(), now()),
  (gen_random_uuid(), 'Comments',         'comments',         'project', now(), now()),
  (gen_random_uuid(), 'Attachments',      'attachments',      'project', now(), now()),
  (gen_random_uuid(), 'Workflows',        'workflows',        'project', now(), now()),
  (gen_random_uuid(), 'Sprints',          'sprints',          'project', now(), now()),
  (gen_random_uuid(), 'Retros',           'retros',           'project', now(), now());

INSERT INTO public.permissions (id, module_id, "name", "key", created_at, updated_at)
SELECT gen_random_uuid(), m.id, p.name, p.key, now(), now()
FROM (VALUES
  ('tenant',           'View organization',          'tenant.read'),
  ('tenant',           'Edit organization',          'tenant.update'),
  ('tenant',           'Delete organization',        'tenant.delete'),
  ('subscriptions',    'View subscription',          'subscription.read'),
  ('subscriptions',    'Change plan',                'subscription.update'),
  ('subscriptions',    'Cancel subscription',        'subscription.cancel'),
  ('workspaces',       'Create workspace',           'workspace.create'),
  ('workspaces',       'View all workspaces',        'workspace.read'),
  ('workspaces',       'Edit workspace',             'workspace.update'),
  ('workspaces',       'Archive workspace',          'workspace.archive'),
  ('workspaces',       'Delete workspace',           'workspace.delete'),
  ('projects',         'Create project',             'project.create'),
  ('users',            'Invite users',               'user.invite'),
  ('users',            'View users',                 'user.read'),
  ('users',            'Change user role',           'user.update'),
  ('users',            'Remove users',               'user.remove'),
  ('roles',            'Create role',                'role.create'),
  ('roles',            'View roles',                 'role.read'),
  ('roles',            'Edit role',                  'role.update'),
  ('roles',            'Delete role',                'role.delete'),
  ('audit',            'View audit log',             'audit.read'),
  ('project_settings', 'View project',               'project.read'),
  ('project_settings', 'Edit project',               'project.update'),
  ('project_settings', 'Archive project',            'project.archive'),
  ('project_settings', 'Delete project',             'project.delete'),
  ('project_members',  'Add member',                 'member.add'),
  ('project_members',  'View members',               'member.read'),
  ('project_members',  'Change member role',         'member.update'),
  ('project_members',  'Remove member',              'member.remove'),
  ('project_roles',    'Create project role',        'project_role.create'),
  ('project_roles',    'View project roles',         'project_role.read'),
  ('project_roles',    'Edit project role',          'project_role.update'),
  ('project_roles',    'Delete project role',        'project_role.delete'),
  ('nodes',            'Create nodes',               'node.create'),
  ('nodes',            'View nodes',                 'node.read'),
  ('nodes',            'Edit nodes',                 'node.update'),
  ('nodes',            'Delete nodes',               'node.delete'),
  ('nodes',            'Assign nodes',               'node.assign'),
  ('nodes',            'Change node status',         'node.transition'),
  ('nodes',            'Move nodes',                 'node.move'),
  ('comments',         'Add comments',               'comment.create'),
  ('comments',         'View comments',              'comment.read'),
  ('comments',         'Edit any comment',           'comment.update'),
  ('comments',         'Delete any comment',         'comment.delete'),
  ('attachments',      'Upload attachments',         'attachment.create'),
  ('attachments',      'View attachments',           'attachment.read'),
  ('attachments',      'Delete any attachment',      'attachment.delete'),
  ('workflows',        'Create workflow',            'workflow.create'),
  ('workflows',        'View workflows',             'workflow.read'),
  ('workflows',        'Edit workflow',              'workflow.update'),
  ('workflows',        'Delete workflow',            'workflow.delete'),
  ('sprints',          'Create sprint',              'sprint.create'),
  ('sprints',          'View sprints',               'sprint.read'),
  ('sprints',          'Edit sprint',                'sprint.update'),
  ('sprints',          'Delete sprint',              'sprint.delete'),
  ('sprints',          'Start sprint',               'sprint.start'),
  ('sprints',          'Close sprint',               'sprint.close'),
  ('retros',           'Create retro',               'retro.create'),
  ('retros',           'View retros',                'retro.read'),
  ('retros',           'Edit retro',                 'retro.update'),
  ('retros',           'Delete retro',               'retro.delete'),
  ('retros',           'Close retro',                'retro.close')
) AS p(module_key, name, key)
JOIN public.modules m ON m.key = p.module_key;

INSERT INTO public.roles (id, tenant_id, project_id, "name", "type", created_by, created_at, updated_at) VALUES
  ('87704fdb-b47b-4b0b-ad9c-e420f26bbe38'::uuid, NULL, NULL, 'Owner',           'tenant',  NULL, now(), now()),
  ('75532888-9966-48c3-a454-06ec9029c00a'::uuid, NULL, NULL, 'Admin',           'tenant',  NULL, now(), now()),
  ('081036dc-d7a9-4d59-9a06-37b0ca13e9a4'::uuid, NULL, NULL, 'Creator',         'tenant',  NULL, now(), now()),
  ('7a0ea870-2e2a-48f8-9058-afc5963c10f9'::uuid, NULL, NULL, 'Member',          'tenant',  NULL, now(), now()),
  ('2d4a3fd2-108c-47a3-9e27-1cc3c55e4542'::uuid, NULL, NULL, 'Project Admin',   'project', NULL, now(), now()),
  ('f62668f9-331a-40bb-b1ed-918bbdb2ac9a'::uuid, NULL, NULL, 'Project Manager', 'project', NULL, now(), now()),
  ('8079fa08-8a9c-42e1-9476-0282e4b7e728'::uuid, NULL, NULL, 'Developer',       'project', NULL, now(), now()),
  ('2178e0e4-ceca-4878-8335-9766720bf062'::uuid, NULL, NULL, 'Viewer',          'project', NULL, now(), now());

INSERT INTO public.role_permissions (id, role_id, permission_id, created_at, updated_at)
SELECT gen_random_uuid(), r.id, p.id, now(), now()
FROM public.roles r
JOIN public.modules m ON m.type::text = r.type::text
JOIN public.permissions p ON p.module_id = m.id
WHERE r.id IN (
  '87704fdb-b47b-4b0b-ad9c-e420f26bbe38'::uuid,
  '2d4a3fd2-108c-47a3-9e27-1cc3c55e4542'::uuid
);

INSERT INTO public.role_permissions (id, role_id, permission_id, created_at, updated_at)
SELECT gen_random_uuid(), r.id, p.id, now(), now()
FROM public.roles r
JOIN public.modules m ON m.type::text = r.type::text
JOIN public.permissions p ON p.module_id = m.id
WHERE r.id = '75532888-9966-48c3-a454-06ec9029c00a'::uuid
  AND p.key NOT IN ('tenant.delete', 'subscription.update', 'subscription.cancel');

INSERT INTO public.role_permissions (id, role_id, permission_id, created_at, updated_at)
SELECT gen_random_uuid(), r.id, p.id, now(), now()
FROM (VALUES
  ('081036dc-d7a9-4d59-9a06-37b0ca13e9a4'::uuid, ARRAY[
    'tenant.read', 'user.read', 'workspace.read', 'workspace.create', 'project.create']),
  ('7a0ea870-2e2a-48f8-9058-afc5963c10f9'::uuid, ARRAY[
    'tenant.read', 'user.read']),
  ('f62668f9-331a-40bb-b1ed-918bbdb2ac9a'::uuid, ARRAY[
    'project.read', 'member.read', 'project_role.read',
    'node.create', 'node.read', 'node.update', 'node.delete', 'node.assign', 'node.transition', 'node.move',
    'comment.create', 'comment.read', 'comment.update', 'comment.delete',
    'attachment.create', 'attachment.read', 'attachment.delete',
    'workflow.create', 'workflow.read', 'workflow.update', 'workflow.delete',
    'sprint.create', 'sprint.read', 'sprint.update', 'sprint.delete', 'sprint.start', 'sprint.close',
    'retro.create', 'retro.read', 'retro.update', 'retro.delete', 'retro.close']),
  ('8079fa08-8a9c-42e1-9476-0282e4b7e728'::uuid, ARRAY[
    'project.read', 'member.read', 'project_role.read',
    'node.create', 'node.read', 'node.update', 'node.assign', 'node.transition', 'node.move',
    'comment.create', 'comment.read',
    'attachment.create', 'attachment.read',
    'workflow.read', 'sprint.read', 'retro.create', 'retro.read']),
  ('2178e0e4-ceca-4878-8335-9766720bf062'::uuid, ARRAY[
    'project.read', 'member.read', 'project_role.read',
    'node.read', 'comment.read', 'attachment.read',
    'workflow.read', 'sprint.read', 'retro.read'])
) AS g(role_id, perm_keys)
CROSS JOIN LATERAL unnest(g.perm_keys) AS k(perm_key)
JOIN public.roles r ON r.id = g.role_id
JOIN public.permissions p ON p.key = k.perm_key
JOIN public.modules m ON m.id = p.module_id AND m.type::text = r.type::text;