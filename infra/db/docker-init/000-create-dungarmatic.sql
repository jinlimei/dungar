DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_roles
        WHERE rolname = 'dungarmatic'
    ) THEN
        CREATE ROLE dungarmatic LOGIN PASSWORD 'dungarmatic';
    END IF;
END
$$;
