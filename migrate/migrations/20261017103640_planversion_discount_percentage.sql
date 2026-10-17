-- Modify "plan_versions" table
ALTER TABLE "public"."plan_versions" ADD COLUMN "discount_percentage" double precision NOT NULL DEFAULT 0;
-- Drop index "plans_name_key" from table: "plans"
DROP INDEX "public"."plans_name_key";
