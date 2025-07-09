package dbHelper

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"storex/database"
	"storex/models"
	"storex/utils"
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
