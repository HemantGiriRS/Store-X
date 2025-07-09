DO $$
DECLARE
admin_uuid UUID := '11111111-1111-1111-1111-111111111111';
BEGIN

INSERT INTO employee_table (id,name,email,phone_no,type,role,created_by,created_at,updated_at)
VALUES (
             admin_uuid,
             'Default Admin',
             'admin@remotestate.com',
             '0000000000',
             'full-time',
             'admin',
             admin_uuid,
             NOW(),
             NOW()
         )
    ON CONFLICT (email) WHERE (archived_at IS NULL) DO NOTHING;
END $$;