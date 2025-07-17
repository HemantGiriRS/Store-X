-- Add the column to the assignment log table
ALTER TABLE assigned_log_table
    ADD COLUMN IF NOT EXISTS retrieved_by UUID REFERENCES employee_table(id);

-- Add the column to the service table
ALTER TABLE service_table
    ADD COLUMN IF NOT EXISTS retrieved_by UUID REFERENCES employee_table(id);