-- Applies the schema exactly once, atomically, even when several replicas
-- start at the same moment.
--
-- Everything is in ONE transaction and ONE session, which is what makes it
-- safe:
--
--   pg_advisory_xact_lock  every replica runs this. Without the lock, three
--                          sessions read "not applied" simultaneously and all
--                          three try to insert the seed rows. The lock is held
--                          until COMMIT and released automatically.
--
--   schema_bootstrap       db/setup/schema.sql is NOT idempotent: it ends with
--                          seed INSERTs that have no ON CONFLICT clause, so a
--                          second run fails on users_pkey. Statement-level
--                          guards (IF NOT EXISTS) cover the CREATEs only, so
--                          the guard has to be at the file level. The checksum
--                          means an edited schema is recognised as new work.
--
--   BEGIN/COMMIT           a pod killed mid-apply rolls back completely rather
--                          than leaving a half-built schema behind.
\set ON_ERROR_STOP on

BEGIN;

SELECT pg_advisory_xact_lock(4021371);

CREATE TABLE IF NOT EXISTS schema_bootstrap (
    checksum   text PRIMARY KEY,
    applied_at timestamptz NOT NULL DEFAULT now()
);

SELECT (NOT EXISTS (
    SELECT 1 FROM schema_bootstrap WHERE checksum = :'checksum'
))::text AS needs_apply \gset

\if :needs_apply
    \echo 'migrate: applying schema'
    \i /schema.sql
    INSERT INTO schema_bootstrap (checksum) VALUES (:'checksum');
\else
    \echo 'migrate: schema already applied, nothing to do'
\endif

COMMIT;
