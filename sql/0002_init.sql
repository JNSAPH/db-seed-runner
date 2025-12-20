-- Create user if not exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'payload_oovolabs') THEN
        CREATE USER payload_oovolabs WITH PASSWORD 'CHANGEME';
    END IF;
END
$$;
