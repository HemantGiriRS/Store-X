-- ENUM TYPES
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'asset_type') THEN
            CREATE TYPE asset_type AS ENUM ('laptop', 'mouse', 'monitor', 'hard-disk', 'pen-drive', 'mobile', 'sim', 'accessories');
        END IF;

        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'asset_status') THEN
            CREATE TYPE asset_status AS ENUM ('available', 'assigned', 'waitForRepair', 'service', 'damage', 'deleted');
        END IF;

        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'asset_owned_by') THEN
            CREATE TYPE asset_owned_by AS ENUM ('RemoteState', 'Client');
        END IF;

        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'employee_type') THEN
            CREATE TYPE employee_type AS ENUM ('full-time', 'intern', 'freelancer');
        END IF;

        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'employee_role') THEN
            CREATE TYPE employee_role AS ENUM ('admin', 'asset_manager', 'employee_manager', 'employee');
        END IF;
    END$$;