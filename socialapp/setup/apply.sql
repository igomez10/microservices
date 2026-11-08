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
--   schema_bootstrap       the checksum means an edited schema.sql is
--                          recognised as new work and re-applied; schema.sql
--                          must therefore stay re-runnable (IF NOT EXISTS).
--                          db/setup/seed.sql is NOT re-runnable and must never
--                          touch a live database again, so it runs only when
--                          schema_bootstrap is empty: the very first apply.
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
))::text AS needs_apply,
(NOT EXISTS (SELECT 1 FROM schema_bootstrap))::text AS first_apply \gset

\if :needs_apply
    \echo 'migrate: applying schema'
    \i /schema.sql
    \if :first_apply
        \echo 'migrate: first bootstrap, applying seed'
        \i /seed.sql
    \endif
    INSERT INTO schema_bootstrap (checksum) VALUES (:'checksum');
\else
    \echo 'migrate: schema already applied, nothing to do'
\endif

COMMIT;
