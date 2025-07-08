-- EMPLOYEE TABLE
CREATE TABLE IF NOT EXISTS employee_table (
                                              id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
                                              name TEXT NOT NULL,
                                              email TEXT NOT NULL UNIQUE,
                                              phone_no TEXT ,
                                              type employee_type DEFAULT 'full-time',
                                              asset_status INT DEFAULT 0,
                                              role employee_role DEFAULT 'employee' ,
                                              created_by UUID,
                                              created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                              updated_at TIMESTAMP DEFAULT NOW(),
                                              updated_by UUID,
                                              archived_at TIMESTAMP,
                                              archived_by UUID
);

-- ASSET TABLE
CREATE TABLE IF NOT EXISTS asset_table (
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                           brand TEXT NOT NULL,
                                           model TEXT NOT NULL,
                                           type asset_type NOT NULL,
                                           serial_no TEXT NOT NULL UNIQUE,
                                           status asset_status NOT NULL DEFAULT 'available',
                                           assigned_to TEXT,
                                           owned_by asset_owned_by NOT NULL,
                                           purchase_date TIMESTAMP,
                                           warranty_start TIMESTAMP,
                                           warranty_end TIMESTAMP,
                                           created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                           created_by UUID,
                                           updated_at TIMESTAMP DEFAULT NOW(),
                                           updated_by UUID,
                                           archived_at TIMESTAMP,
                                           archived_by UUID
);

-- SERVICE TABLE
CREATE TABLE IF NOT EXISTS service_table (
                                             id UUID PRIMARY KEY,
                                             assigned_to TEXT NOT NULL,
                                             assigned_by UUID NOT NULL REFERENCES employee_table(id),
                                             asset_id UUID NOT NULL REFERENCES asset_table(id) ON DELETE CASCADE,
                                             price TEXT,
                                             description TEXT,
                                             assigned_date TIMESTAMP NOT NULL,
                                             received_date TIMESTAMP
);

-- ASSIGNED LOG TABLE
CREATE TABLE IF NOT EXISTS assigned_log_table (
                                                  id UUID PRIMARY KEY,
                                                  asset_id UUID NOT NULL REFERENCES asset_table(id) ON DELETE CASCADE,
                                                  employee_id UUID NOT NULL REFERENCES employee_table(id) ON DELETE CASCADE,
                                                  reason_of_retrieval TEXT,
                                                  start_at TIMESTAMP NOT NULL,
                                                  end_at TIMESTAMP
);

-- INDEXES
-- For asset_table
CREATE INDEX IF NOT EXISTS idx_asset_status ON asset_table(status);
CREATE INDEX IF NOT EXISTS idx_asset_type ON asset_table(type);
CREATE INDEX IF NOT EXISTS idx_asset_serial_no ON asset_table(serial_no);
CREATE INDEX IF NOT EXISTS idx_asset_owned_by ON asset_table(owned_by);

-- For employee_table
CREATE INDEX IF NOT EXISTS idx_employee_email ON employee_table(email);
CREATE INDEX IF NOT EXISTS idx_employee_role ON employee_table(role);
CREATE INDEX IF NOT EXISTS idx_employee_type ON employee_table(type);

-- For service_table
CREATE INDEX IF NOT EXISTS idx_service_asset_id ON service_table(asset_id);
CREATE INDEX IF NOT EXISTS idx_service_assigned_to ON service_table(assigned_to);

-- For assigned_log_table
CREATE INDEX IF NOT EXISTS idx_log_asset_id ON assigned_log_table(asset_id);
CREATE INDEX IF NOT EXISTS idx_log_employee_id ON assigned_log_table(employee_id);