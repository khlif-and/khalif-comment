CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

--SEPARATOR--

DROP TRIGGER IF EXISTS update_comments_modtime ON comments;

--SEPARATOR--

CREATE TRIGGER update_comments_modtime BEFORE UPDATE ON comments FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

--SEPARATOR--

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'readonly_user') THEN
        CREATE ROLE readonly_user WITH LOGIN PASSWORD 'readonly_password';
    END IF;
END
$$;

--SEPARATOR--

-- Grant permissions (Optional: sesuaikan nama DB jika perlu, default di script ini biasanya berjalan di DB aktif)
GRANT USAGE ON SCHEMA public TO readonly_user;

--SEPARATOR--

GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly_user;

--SEPARATOR--

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO readonly_user;