
  INSERT INTO public."plans"
  (id, "name", "key", created_at, updated_at, deleted_at)
  VALUES('7ccd4553-bb6b-412d-a4f3-89fb863a2499'::uuid, 'Free', 'free', '2026-10-17 15:51:02.861', '2026-10-17 15:51:02.861', NULL);
  INSERT INTO public."plans"
  (id, "name", "key", created_at, updated_at, deleted_at)
  VALUES('94ab0ca3-e6c9-4027-b4ca-c59ac5868c08'::uuid, 'Premium', 'premium.monthly', '2026-10-17 15:51:02.861', '2026-10-17 15:51:02.861', NULL);
  INSERT INTO public."plans"
  (id, "name", "key", created_at, updated_at, deleted_at)
  VALUES('ef6c4898-6608-49e4-841e-f3c64c7d9234'::uuid, 'Premium', 'premium.yearly', '2026-10-17 15:51:02.861', '2026-10-17 15:51:02.861', NULL);


  INSERT INTO public.plan_versions
  (id, "version", "cost", discount_percentage, currency, "interval", perks, created_at, updated_at, deleted_at, plan_id)
  VALUES('6f7d5b29-9b19-48ac-89dc-e8602d12e387'::uuid, 1, 0, 0, 'usd', 30,
  '{"limits":{"seats":20,"workspaces":1,"projects":3,"nodes":25000,"storage_bytes":2147483648},"features":{"custom_roles":false,"custom_workflows":false,"audit_log":false}}'::jsonb,
  '2026-10-17 15:52:01.153', '2026-10-17 15:52:01.153', NULL, '7ccd4553-bb6b-412d-a4f3-89fb863a2499'::uuid);

  INSERT INTO public.plan_versions
  (id, "version", "cost", discount_percentage, currency, "interval", perks, created_at, updated_at, deleted_at, plan_id)
  VALUES('057063f2-f6c5-4318-994b-39023fea716c'::uuid, 1, 8.0, 0, 'usd', 30,
  '{"limits":{"seats":-1,"workspaces":-1,"projects":-1,"nodes":1000000,"storage_bytes":107374182400},"features":{"custom_roles":true,"custom_workflows":true,"audit_log":true}}'::jsonb,
  '2026-10-17 15:52:01.163', '2026-10-17 15:52:01.163', NULL, '94ab0ca3-e6c9-4027-b4ca-c59ac5868c08'::uuid);

  INSERT INTO public.plan_versions
  (id, "version", "cost", discount_percentage, currency, "interval", perks, created_at, updated_at, deleted_at, plan_id)
  VALUES('c56bbe99-4279-435a-a994-7c80f4719da6'::uuid, 1, 8.0, 15, 'usd', 365,
  '{"limits":{"seats":-1,"workspaces":-1,"projects":-1,"nodes":1000000,"storage_bytes":107374182400},"features":{"custom_roles":true,"custom_workflows":true,"audit_log":true}}'::jsonb,
  '2026-10-17 15:52:01.173', '2026-10-17 15:52:01.173', NULL, 'ef6c4898-6608-49e4-841e-f3c64c7d9234'::uuid);


