-- Modify "tenants" table
ALTER TABLE "public"."tenants" ALTER COLUMN "name" TYPE public.citext, ALTER COLUMN "slug" TYPE public.citext;
-- Modify "workspaces" table
ALTER TABLE "public"."workspaces" ALTER COLUMN "name" TYPE public.citext;
-- Create index "workspace_tenant_id_id" to table: "workspaces"
CREATE UNIQUE INDEX "workspace_tenant_id_id" ON "public"."workspaces" ("tenant_id", "id");
-- Create index "workspace_tenant_id_name" to table: "workspaces"
CREATE UNIQUE INDEX "workspace_tenant_id_name" ON "public"."workspaces" ("tenant_id", "name");
-- Modify "projects" table
ALTER TABLE "public"."projects" ALTER COLUMN "name" TYPE public.citext, ADD COLUMN "created_by" uuid NOT NULL, ADD CONSTRAINT "projects_users_projectcreator" FOREIGN KEY ("created_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "project_tenant_id_workspace_id" to table: "projects"
CREATE INDEX "project_tenant_id_workspace_id" ON "public"."projects" ("tenant_id", "workspace_id");
-- Create index "project_tenant_id_workspace_id_name" to table: "projects"
CREATE INDEX "project_tenant_id_workspace_id_name" ON "public"."projects" ("tenant_id", "workspace_id", "name");
