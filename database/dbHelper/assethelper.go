package dbHelper

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"storex/database"
	"storex/models"
	"strings"
)

func CreateAsset(tx *sqlx.Tx, asset models.BaseAsset, createdByID string) (string, error) {
	SQL := `INSERT INTO asset_table 
                (brand, model, type, serial_no, owned_by, created_by, purchase_date, warranty_start, warranty_end) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
            RETURNING id`

	var assetID string

	err := tx.Get(&assetID, SQL, asset.Brand, asset.Model, asset.Type, asset.SerialNo, asset.OwnedBy, createdByID, asset.PurchaseDate, asset.WarrantyStart, asset.WarrantyEnd)
	if err != nil {
		return "", fmt.Errorf("failed to create asset in asset_table: %w", err)
	}
	return assetID, nil
}

func CreateLaptopSpec(tx *sqlx.Tx, assetID string, spec models.LaptopSpecs) error {
	SQL := `INSERT INTO laptop_specs (asset_id, processor, ram_gb, storage_gb, os, screen_size_inch, battery_backup_hours) 
            VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := tx.Exec(SQL, assetID, spec.Processor, spec.RAM_GB, spec.Storage_GB, spec.OS, spec.ScreenSizeInch, spec.BatteryBackupHours)
	if err != nil {
		return fmt.Errorf("failed to create laptop specs: %w", err)
	}
	return nil
}

func CreateMouseSpec(tx *sqlx.Tx, assetID string, spec models.MouseSpecs) error {
	SQL := `INSERT INTO mouse_specs (asset_id, connection_type, dpi, number_of_buttons) 
            VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(SQL, assetID, spec.ConnectionType, spec.DPI, spec.NumberOfButtons)
	if err != nil {
		return fmt.Errorf("failed to create mouse specs: %w", err)
	}
	return nil
}

func CreateMonitorSpec(tx *sqlx.Tx, assetID string, spec models.MonitorSpecs) error {
	SQL := `INSERT INTO monitor_specs (asset_id, size_inch, resolution, refresh_rate_hz, panel_type) 
            VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(SQL, assetID, spec.SizeInch, spec.Resolution, spec.RefreshRateHZ, spec.PanelType)
	if err != nil {
		return fmt.Errorf("failed to create monitor specs: %w", err)
	}
	return nil
}

func CreateHardDiskSpec(tx *sqlx.Tx, assetID string, spec models.HardDiskSpecs) error {
	SQL := `INSERT INTO hard_disk_specs (asset_id, capacity_gb, type, interface) 
            VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(SQL, assetID, spec.CapacityGB, spec.Type, spec.Interface)
	if err != nil {
		return fmt.Errorf("failed to create hard disk specs: %w", err)
	}
	return nil
}

func CreatePenDriveSpec(tx *sqlx.Tx, assetID string, spec models.PenDriveSpecs) error {
	SQL := `INSERT INTO pen_drive_specs (asset_id, capacity_gb, usb_type) 
            VALUES ($1, $2, $3)`
	_, err := tx.Exec(SQL, assetID, spec.CapacityGB, spec.USBType)
	if err != nil {
		return fmt.Errorf("failed to create pen drive specs: %w", err)
	}
	return nil
}

func CreateMobileSpec(tx *sqlx.Tx, assetID string, spec models.MobileSpecs) error {
	SQL := `INSERT INTO mobile_specs (asset_id, imei, ram_gb, storage_gb, os, screen_size_inch, battery_capacity_mah) 
            VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := tx.Exec(SQL, assetID, spec.IMEI, spec.RAM_GB, spec.Storage_GB, spec.OS, spec.ScreenSizeInch, spec.BatteryCapacityMAH)
	if err != nil {
		return fmt.Errorf("failed to create mobile specs: %w", err)
	}
	return nil
}

func CreateSimSpec(tx *sqlx.Tx, assetID string, spec models.SimSpecs) error {
	SQL := `INSERT INTO sim_specs (asset_id, phone_number, operator, sim_type, plan_details) 
            VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(SQL, assetID, spec.PhoneNumber, spec.Operator, spec.SimType, spec.PlanDetails)
	if err != nil {
		return fmt.Errorf("failed to create sim specs: %w", err)
	}
	return nil
}

func CreateAccessoriesSpec(tx *sqlx.Tx, assetID string, spec models.AccessoriesSpecs) error {
	SQL := `INSERT INTO accessories_specs (asset_id, name, description, compatibility) 
            VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(SQL, assetID, spec.Name, spec.Description, spec.Compatibility)
	if err != nil {
		return fmt.Errorf("failed to create accessories specs: %w", err)
	}
	return nil
}

func GetAssetStatus(tx *sqlx.Tx, assetID string) (string, error) {
	SQL := `SELECT status FROM asset_table WHERE id = $1 FOR UPDATE`
	var status string

	err := tx.Get(&status, SQL, assetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("asset with ID %s not found", assetID)
		}
		return "", fmt.Errorf("failed to get asset status: %w", err)
	}
	// The "FOR UPDATE" clause locks the selected row to prevent another transaction
	// from assigning this asset at the exact same time (a race condition).
	// The lock is automatically released when the transaction ends (commit or rollback).
	return status, nil
}

func UpdateAssetForAssignment(tx *sqlx.Tx, assetID, employeeID, updatedByID string) error {
	SQL := `UPDATE asset_table 
            SET status = 'assigned', assigned_to = $1, updated_by = $2, updated_at = NOW() 
            WHERE id = $3`

	_, err := tx.Exec(SQL, employeeID, updatedByID, assetID)
	if err != nil {
		return fmt.Errorf("failed to update asset_table for assignment: %w", err)
	}

	return nil
}

func CreateAssignmentLog(tx *sqlx.Tx, assetID, employeeID, assignedByID string) error {
	SQL := `INSERT INTO assigned_log_table (asset_id, employee_id, assigned_by, start_at) 
            VALUES ($1, $2, $3, NOW())`

	_, err := tx.Exec(SQL, assetID, employeeID, assignedByID)
	if err != nil {
		return fmt.Errorf("failed to create assignment log: %w", err)
	}
	return nil
}
func GetAssetsFiltered(filters map[string]string, page, pageSize int) ([]models.AssetDetail, int64, error) {
	var assets []models.AssetDetail
	var totalRecords int64

	baseSelect := `
        SELECT 
            a.id, a.brand, a.model, a.type, a.serial_no, a.status, 
            a.owned_by, a.purchase_date, a.warranty_end, a.created_at,
            a.assigned_to as assigned_to_id,
            e.name as assigned_to_name,
            e.email as assigned_to_email
        FROM 
            asset_table a
        LEFT JOIN 
            employee_table e ON a.assigned_to::uuid = e.id AND e.archived_at IS NULL
    `
	baseCount := `SELECT COUNT(a.id) FROM asset_table a`
	whereClauses := []string{"a.archived_at IS NULL"}
	args := []interface{}{}

	if query, ok := filters["q"]; ok && query != "" {
		whereClauses = append(whereClauses, "(a.brand ILIKE ? OR a.model ILIKE ? OR a.serial_no ILIKE ?)")
		likeQuery := "%" + query + "%"
		args = append(args, likeQuery, likeQuery, likeQuery)
	}
	if assetType, ok := filters["type"]; ok && assetType != "" {
		whereClauses = append(whereClauses, "a.type = ?")
		args = append(args, assetType)
	}
	if status, ok := filters["status"]; ok && status != "" {
		whereClauses = append(whereClauses, "a.status = ?")
		args = append(args, status)
	}
	if ownedBy, ok := filters["owned_by"]; ok && ownedBy != "" {
		whereClauses = append(whereClauses, "a.owned_by = ?")
		args = append(args, ownedBy)
	}

	whereStatement := "WHERE " + strings.Join(whereClauses, " AND ")

	countQueryString := fmt.Sprintf("%s %s", baseCount, whereStatement)
	countQuery := database.SX.Rebind(countQueryString)
	err := database.SX.Get(&totalRecords, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute asset count query: %w", err)
	}

	if totalRecords == 0 {
		return []models.AssetDetail{}, 0, nil
	}

	pagedArgs := append(args, pageSize, (page-1)*pageSize)
	selectQueryString := fmt.Sprintf("%s %s ORDER BY a.created_at DESC LIMIT ? OFFSET ?", baseSelect, whereStatement)

	selectQuery := database.SX.Rebind(selectQueryString)
	err = database.SX.Select(&assets, selectQuery, pagedArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute asset select query: %w", err)
	}

	return assets, totalRecords, nil
}

func GetAssetForUnassignment(tx *sqlx.Tx, assetID string) (*models.AssetAssignmentInfo, error) {
	var info models.AssetAssignmentInfo
	SQL := `
				SELECT 
				    status, assigned_to 
				FROM 
				    asset_table 
				WHERE 
				    id = $1 
				FOR UPDATE
				`
	err := tx.Get(&info, SQL, assetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset with ID %s not found", assetID)
		}
		return nil, fmt.Errorf("failed to get asset for unassignment: %w", err)
	}
	return &info, nil
}

func UpdateAssetForUnassignment(tx *sqlx.Tx, assetID, updatedByID string) error {
	SQL := `UPDATE asset_table 
            SET status = 'available', assigned_to = NULL, updated_by = ?, updated_at = NOW() 
            WHERE id = ?`

	query := database.SX.Rebind(SQL)
	_, err := tx.Exec(query, updatedByID, assetID)
	if err != nil {
		return fmt.Errorf("failed to update asset for unassignment: %w", err)
	}
	return nil
}

func CloseAssignmentLog(tx *sqlx.Tx, assetID, reason, retrievedByID string) error {
	SQL := `UPDATE assigned_log_table 
            SET end_at = NOW(), reason_of_retrieval = ?, retrieved_by = ?
            WHERE asset_id = ? AND end_at IS NULL`

	query := database.SX.Rebind(SQL)
	// Add retrievedByID to the arguments
	result, err := tx.Exec(query, reason, retrievedByID, assetID)
	if err != nil {
		return fmt.Errorf("failed to close assignment log: %w", err)
	}

	// ... (rest of the function is the same) ...
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify log update: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no active assignment log found for asset ID %s", assetID)
	}

	return nil
}
