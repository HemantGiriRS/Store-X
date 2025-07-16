package dbHelper

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"storex/database"
	"storex/models"
	"storex/utils"
	"strings"
)

func FindUserByEmail(email string) (*models.UserContext, error) {
	SQL := `SELECT id , role FROM employee_table WHERE email = $1`
	args := models.UserContext{
		ID:   "",
		Role: "",
	}
	err := database.SX.Get(&args, SQL, email)
	if err != nil {
		//fmt.Println(err)
		return &args, err
	}
	return &args, nil
}

func RegisterUser(email string) (*models.UserContext, error) {
	SQL := `INSERT INTO employee_table (name ,email) VALUES ($1, $2) RETURNING id , role`
	args := models.UserContext{
		ID:   "",
		Role: "",
	}
	name := utils.GetName(email)
	err := database.SX.Get(&args, SQL, name, email)
	if err != nil {
		return &args, err
	}
	return &args, nil
}

func CreateUser(user models.SignUpRequest, role string, createdBy uuid.UUID) (string, error) {
	name := utils.GetName(user.Email)

	SQL := `INSERT INTO employee_table (name, email, phone_no, type, role, created_by) 
            VALUES ($1, $2, $3, $4, $5, $6) 
            RETURNING id`

	var id string
	err := database.SX.Get(&id, SQL, name, user.Email, user.PhoneNo, user.Type, role, createdBy)
	if err != nil {
		return "", err
	}

	return id, nil
}

func IncrementEmployeeAssetCount(tx *sqlx.Tx, employeeID string) error {
	SQL := `UPDATE employee_table SET asset_status = asset_status + 1 WHERE id = $1`
	_, err := tx.Exec(SQL, employeeID)
	if err != nil {
		return fmt.Errorf("failed to increment asset count for employee %s: %w", employeeID, err)
	}
	return nil
}

func UpdateEmployeeRole(employeeID, newRole, actorID string) error {
	SQL := `UPDATE employee_table SET role = $1, updated_at = NOW() , updated_by =$2 WHERE id = $3`

	result, err := database.SX.Exec(SQL, newRole, actorID, employeeID)
	if err != nil {
		return fmt.Errorf("database error during role update: %w", err)
	}

	// Check if any row was actually updated.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check affected rows: %w", err)
	}

	// If RowsAffected is 0, it means the WHERE clause (id = $2) did not find a match.
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func GetEmployeesFiltered(filters map[string]string, page, pageSize int) ([]models.EmployeeDetail, int64, error) {
	var employees []models.EmployeeDetail
	var totalRecords int64

	whereClauses := []string{"archived_at IS NULL"}
	args := []interface{}{}

	if query, ok := filters["q"]; ok && query != "" {
		whereClauses = append(whereClauses, "(name ILIKE ? OR email ILIKE ? OR phone_no ILIKE ?)")
		likeQuery := "%" + query + "%"
		args = append(args, likeQuery, likeQuery, likeQuery)
	}

	if role, ok := filters["role"]; ok && role != "" {
		whereClauses = append(whereClauses, "role = ?")
		args = append(args, role)
	}

	if empType, ok := filters["type"]; ok && empType != "" {
		whereClauses = append(whereClauses, "type = ?")
		args = append(args, empType)
	}

	if status, ok := filters["assignment_status"]; ok && status != "" {
		if status == "assigned" {
			whereClauses = append(whereClauses, "asset_status > 0")
		} else if status == "not_assigned" {
			whereClauses = append(whereClauses, "asset_status = 0")
		}
	}

	whereStatement := ""
	if len(whereClauses) > 0 {
		whereStatement = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// --- 1. Execute COUNT query ---
	countQueryString := fmt.Sprintf("SELECT COUNT(*) FROM employee_table %s", whereStatement)
	countQuery := database.SX.Rebind(countQueryString)
	err := database.SX.Get(&totalRecords, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute count query: %w", err)
	}

	if totalRecords == 0 {
		return []models.EmployeeDetail{}, 0, nil
	}

	// --- 2. Execute SELECT query ---
	selectQueryString := fmt.Sprintf(`
        SELECT 
            id, name, email, phone_no, type, role, asset_status, created_at, updated_at 
        FROM 
            employee_table
        %s
        ORDER BY 
            created_at DESC 
        LIMIT ? OFFSET ?`, whereStatement)

	offset := (page - 1) * pageSize
	pagedArgs := append(args, pageSize, offset)

	selectQuery := database.SX.Rebind(selectQueryString)
	err = database.SX.Select(&employees, selectQuery, pagedArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute select query: %w", err)
	}

	return employees, totalRecords, nil
}

func DecrementEmployeeAssetCount(tx *sqlx.Tx, employeeID string) error {
	// We ensure the count never goes below zero.
	SQL := `UPDATE employee_table SET asset_status = asset_status - 1 WHERE id = ? AND asset_status > 0`
	query := database.SX.Rebind(SQL)
	_, err := tx.Exec(query, employeeID)
	if err != nil {
		return fmt.Errorf("failed to decrement asset count for employee %s: %w", employeeID, err)
	}
	return nil
}

func GetEmployeeStateForDeletion(tx *sqlx.Tx, employeeID string) (*models.EmployeeState, error) {
	var state models.EmployeeState
	query := "SELECT asset_status, archived_at FROM employee_table WHERE id = $1 FOR UPDATE"
	err := tx.Get(&state, query, employeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employee with ID %s not found", employeeID)
		}
		return nil, fmt.Errorf("failed to get employee state for deletion: %w", err)
	}
	return &state, nil
}

func SoftDeleteEmployee(tx *sqlx.Tx, employeeID, archiverID string) error {
	SQL := `
        UPDATE employee_table 
        SET 
            archived_at = NOW(), 
            archived_by = ?,
            updated_at = NOW(),
            updated_by = ?
        WHERE id = ?`

	query := database.SX.Rebind(SQL)
	_, err := tx.Exec(query, archiverID, archiverID, employeeID)
	if err != nil {
		return fmt.Errorf("failed to soft delete employee: %w", err)
	}
	return nil
}

func GetEmployeeTimeline(employeeID string, page, pageSize int) ([]models.TimelineEvent, int64, error) {
	var timelineEvents []models.TimelineEvent
	var totalRecords int64

	// --- 1. Count all assignment and retrieval events for the employee ---
	countQueryString := `
        SELECT
            (SELECT COUNT(*) FROM assigned_log_table WHERE employee_id = $1) +
            (SELECT COUNT(*) FROM assigned_log_table WHERE employee_id = $1 AND end_at IS NOT NULL)
    `
	err := database.SX.Get(&totalRecords, countQueryString, employeeID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute employee timeline count query: %w", err)
	}

	if totalRecords == 0 {
		return []models.TimelineEvent{}, 0, nil
	}

	// --- 2. Fetch the unified, sorted, and paginated data ---
	selectQueryString := `
        SELECT * FROM (
            -- Event Type 1: Asset Assignment
            SELECT 
                al.id::text,
                'Assignment' AS event_type, 
                al.start_at AS event_date,
                'Asset assigned: ' || a.brand || ' ' || a.model || ' (S/N: ' || a.serial_no || ')' AS description,
                actor.name AS actor_name
            FROM assigned_log_table al
            JOIN asset_table a ON al.asset_id = a.id
            JOIN employee_table actor ON al.assigned_by = actor.id
            WHERE al.employee_id = ?

            UNION ALL

            -- Event Type 2: Asset Retrieval
            SELECT 
                al.id::text || '-ret' AS id,
                'Retrieval' AS event_type, 
                al.end_at AS event_date,
                'Asset retrieved: ' || a.brand || ' ' || a.model || '. Reason: ' || COALESCE(al.reason_of_retrieval, 'N/A') AS description,
                retriever.name AS actor_name
            FROM assigned_log_table al
            JOIN asset_table a ON al.asset_id = a.id
            JOIN employee_table retriever ON al.retrieved_by = retriever.id
            WHERE al.employee_id = ? AND al.end_at IS NOT NULL

        ) AS employee_timeline
        ORDER BY event_date DESC
        LIMIT ? OFFSET ?`

	offset := (page - 1) * pageSize
	pagedArgs := []interface{}{employeeID, employeeID, pageSize, offset}

	selectQuery := database.SX.Rebind(selectQueryString)
	err = database.SX.Select(&timelineEvents, selectQuery, pagedArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute employee timeline select query: %w", err)
	}

	return timelineEvents, totalRecords, nil
}
