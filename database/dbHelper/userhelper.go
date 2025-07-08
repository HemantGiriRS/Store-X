package dbHelper

import (
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

func CreateUser(email string) (*models.UserContext, error) {
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
