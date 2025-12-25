DROP DATABASE IF EXISTS mydatabase;
CREATE DATABASE mydatabase;

DO
$do$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = 'myuser') THEN
      CREATE ROLE myuser LOGIN PASSWORD 'mypassword';
   END IF;
END
$do$;

GRANT ALL PRIVILEGES ON DATABASE mydatabase TO myuser;