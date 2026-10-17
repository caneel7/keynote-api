WITH kind AS (
  INSERT INTO public.entity_kinds (id, "name", "key", created_at, updated_at, deleted_at)
  VALUES ('01a149bb-b200-7c3a-9e41-6f2d8b05a1c7'::uuid, 'Subscription', 'subscription', now(), now(), NULL)
  RETURNING id
)
INSERT INTO public.state_definitions (id, "key", "name", created_at, updated_at, deleted_at, entity_id)
VALUES
  ('01a149bb-b201-7d4b-8a52-7a3e9c16b2d8'::uuid, 'trialing',   'Trialing',   now(), now(), NULL, (SELECT id FROM kind)),
  ('01a149bb-b202-7e5c-b963-8b4fad27c3e9'::uuid, 'active',     'Active',     now(), now(), NULL, (SELECT id FROM kind)),
  ('01a149bb-b203-7f6d-a074-9c50be38d4fa'::uuid, 'past_due',   'Past Due',   now(), now(), NULL, (SELECT id FROM kind)),
  ('01a149bb-b204-7a7e-9185-ad61cf49e50b'::uuid, 'cancelled',  'Cancelled',  now(), now(), NULL, (SELECT id FROM kind)),
  ('01a149bb-b205-7b8f-8296-be72d05af61c'::uuid, 'unpaid',     'Unpaid',     now(), now(), NULL, (SELECT id FROM kind)),
  ('01a149bb-b206-7c90-b3a7-cf83e16b072d'::uuid, 'paused',     'Paused',     now(), now(), NULL, (SELECT id FROM kind)),
  ('01a149bb-b207-7da1-a4b8-d094f27c183e'::uuid, 'incomplete', 'Incomplete', now(), now(), NULL, (SELECT id FROM kind));