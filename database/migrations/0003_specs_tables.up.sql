-- ALL SPECIFICATION TABLES
CREATE TABLE IF NOT EXISTS laptop_specs (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                            asset_id UUID UNIQUE NOT NULL REFERENCES asset_table(id),
                                            processor TEXT NOT NULL,
                                            ram_gb INT NOT NULL,
                                            storage_gb INT NOT NULL,
                                            os TEXT NOT NULL,
                                            screen_size_inch INT,
                                            battery_backup_hours INT,
                                            updated_at TIMESTAMP DEFAULT NOW(),
                                            archived_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS mouse_specs (
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                           asset_id UUID UNIQUE NOT NULL REFERENCES asset_table(id),
                                           connection_type TEXT NOT NULL CHECK (connection_type IN ('wired', 'wireless')),
                                           dpi INT,
                                           number_of_buttons INT,
                                           updated_at TIMESTAMP DEFAULT NOW(),
                                           archived_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS monitor_specs (
                                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                             asset_id UUID UNIQUE NOT NULL REFERENCES asset_table(id),
                                             size_inch INT NOT NULL,
                                             resolution TEXT NOT NULL,
                                             refresh_rate_hz INT,
                                             panel_type TEXT,
                                             updated_at TIMESTAMP DEFAULT NOW(),
                                             archived_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS hard_disk_specs (
                                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                               asset_id UUID UNIQUE NOT NULL REFERENCES asset_table(id),
                                               capacity_gb INT NOT NULL,
                                               type TEXT NOT NULL CHECK (type IN ('HDD', 'SSD')),
                                               interface TEXT,
                                               updated_at TIMESTAMP DEFAULT NOW(),
                                               archived_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pen_drive_specs (
                                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                               asset_id UUID UNIQUE NOT NULL REFERENCES asset_table(id),
                                               capacity_gb INT NOT NULL,
                                               usb_type TEXT NOT NULL CHECK (usb_type IN ('USB 2.0', 'USB 3.0', 'USB-C')),
                                               updated_at TIMESTAMP DEFAULT NOW(),
                                               archived_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS mobile_specs (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                            asset_id UUID UNIQUE NOT NULL REFERENCES asset_table(id),
                                            imei TEXT NOT NULL,
                                            ram_gb INT NOT NULL,
                                            storage_gb INT NOT NULL,
                                            os TEXT NOT NULL,
                                            screen_size_inch INT,
                                            battery_capacity_mah INT,
                                            updated_at TIMESTAMP DEFAULT NOW(),
                                            archived_at TIMESTAMP
);

-- SOLUTION: Partial index for mobile_specs IMEI to support soft-deleting
CREATE UNIQUE INDEX IF NOT EXISTS idx_mobile_specs_imei_unique_if_not_archived
    ON mobile_specs(imei)
    WHERE (archived_at IS NULL);


CREATE TABLE IF NOT EXISTS sim_specs(
                                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                        asset_id UUID UNIQUE NOT NULL REFERENCES asset_table(id),
                                        phone_number TEXT NOT NULL,
                                        operator TEXT NOT NULL,
                                        sim_type TEXT CHECK (sim_type IN ('physical', 'eSIM')),
                                        plan_details TEXT,
                                        updated_at TIMESTAMP DEFAULT NOW(),
                                        archived_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_sim_specs_phone_number_unique_if_not_archived
    ON sim_specs(phone_number)
    WHERE (archived_at IS NULL);

CREATE TABLE IF NOT EXISTS accessories_specs (
                                                 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                 asset_id UUID UNIQUE NOT NULL REFERENCES asset_table(id),
                                                 name TEXT NOT NULL,
                                                 description TEXT,
                                                 compatibility TEXT,
                                                 updated_at TIMESTAMP DEFAULT NOW(),
                                                 archived_at TIMESTAMP
);