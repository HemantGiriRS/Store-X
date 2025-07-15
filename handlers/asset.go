package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"storex/database"
	"storex/database/dbHelper"
	"storex/middlewares"
	"storex/models"
	"storex/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SpecHandler func(tx *sqlx.Tx, assetID string, specsData json.RawMessage) error

// SpecHandlers is our map (registry) that maps an asset type string to its handler function.
var SpecHandlers = map[string]SpecHandler{
	"laptop":      handleLaptopSpecs,
	"mouse":       handleMouseSpecs,
	"monitor":     handleMonitorSpecs,
	"hard-disk":   handleHardDiskSpecs,
	"pen-drive":   handlePenDriveSpecs,
	"mobile":      handleMobileSpecs,
	"sim":         handleSimSpecs,
	"accessories": handleAccessoriesSpecs,
}

var ErrUnknownAssetType = errors.New("unknown or unsupported asset type")

func handleLaptopSpecs(tx *sqlx.Tx, assetID string, specsData json.RawMessage) error {
	var spec models.LaptopSpecs
	if err := json.Unmarshal(specsData, &spec); err != nil {
		return fmt.Errorf("invalid laptop specs format: %w", err)
	}
	return dbHelper.CreateLaptopSpec(tx, assetID, spec)
}

func handleMouseSpecs(tx *sqlx.Tx, assetID string, specsData json.RawMessage) error {
	var spec models.MouseSpecs
	if err := json.Unmarshal(specsData, &spec); err != nil {
		return fmt.Errorf("invalid mouse specs format: %w", err)
	}
	return dbHelper.CreateMouseSpec(tx, assetID, spec)
}

func handleMonitorSpecs(tx *sqlx.Tx, assetID string, specsData json.RawMessage) error {
	var spec models.MonitorSpecs
	if err := json.Unmarshal(specsData, &spec); err != nil {
		return fmt.Errorf("invalid monitor specs format: %w", err)
	}
	return dbHelper.CreateMonitorSpec(tx, assetID, spec)
}

func handleHardDiskSpecs(tx *sqlx.Tx, assetID string, specsData json.RawMessage) error {
	var spec models.HardDiskSpecs
	if err := json.Unmarshal(specsData, &spec); err != nil {
		return fmt.Errorf("invalid hard-disk specs format: %w", err)
	}
	return dbHelper.CreateHardDiskSpec(tx, assetID, spec)
}

func handlePenDriveSpecs(tx *sqlx.Tx, assetID string, specsData json.RawMessage) error {
	var spec models.PenDriveSpecs
	if err := json.Unmarshal(specsData, &spec); err != nil {
		return fmt.Errorf("invalid pen-drive specs format: %w", err)
	}
	return dbHelper.CreatePenDriveSpec(tx, assetID, spec)
}

func handleMobileSpecs(tx *sqlx.Tx, assetID string, specsData json.RawMessage) error {
	var spec models.MobileSpecs
	if err := json.Unmarshal(specsData, &spec); err != nil {
		return fmt.Errorf("invalid mobile specs format: %w", err)
	}
	return dbHelper.CreateMobileSpec(tx, assetID, spec)
}

func handleSimSpecs(tx *sqlx.Tx, assetID string, specsData json.RawMessage) error {
	var spec models.SimSpecs
	if err := json.Unmarshal(specsData, &spec); err != nil {
		return fmt.Errorf("invalid sim specs format: %w", err)
	}
	return dbHelper.CreateSimSpec(tx, assetID, spec)
}

func handleAccessoriesSpecs(tx *sqlx.Tx, assetID string, specsData json.RawMessage) error {
	var spec models.AccessoriesSpecs
	if err := json.Unmarshal(specsData, &spec); err != nil {
		return fmt.Errorf("invalid accessories specs format: %w", err)
	}
	return dbHelper.CreateAccessoriesSpec(tx, assetID, spec)
}

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	creatorIDVal := r.Context().Value(middlewares.UserIDKey)
	if creatorIDVal == nil {
		http.Error(w, "unauthorized: could not identify user", http.StatusUnauthorized)
		return
	}
	creatorID, ok := creatorIDVal.(uuid.UUID)
	if !ok {
		http.Error(w, "internal server error: user ID has invalid type", http.StatusInternalServerError)
		return
	}

	var req models.CreateAssetRequest
	if err := utils.DecodeRequest(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	//fmt.Println(req)

	req.Type = strings.ToLower(req.Type)
	var assetID string // To store the new asset's ID.

	if _, found := SpecHandlers[req.Type]; !found {
		http.Error(w, ErrUnknownAssetType.Error(), http.StatusBadRequest)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		newID, err := dbHelper.CreateAsset(tx, req.BaseAsset, creatorID.String())
		if err != nil {
			return err
		}
		assetID = newID

		handler := SpecHandlers[req.Type]

		// Execute the specific handler for the asset's specs.
		return handler(tx, assetID, req.Specs)
	})

	if txErr != nil {
		logrus.Errorf("transaction failed for creating asset: %v", txErr)
		http.Error(w, "failed to create asset. the operation was rolled back.", http.StatusInternalServerError)
		return
	}

	// If successful, use your existing EncodeResponse function.
	res := map[string]string{
		"message":  "Asset created successfully",
		"asset_id": assetID,
	}
	if err := utils.EncodeResponse(w, http.StatusCreated, res); err != nil {
		logrus.Error(err)
		http.Error(w, "failed in sending response", http.StatusInternalServerError)
	}
}

func AssignAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["asset_id"]

	assignerIDVal := r.Context().Value(middlewares.UserIDKey)
	if assignerIDVal == nil {
		http.Error(w, "unauthorized: could not identify assigner", http.StatusUnauthorized)
		return
	}
	assignerID, _ := assignerIDVal.(uuid.UUID)

	var req models.AssignAssetRequest
	if err := utils.DecodeRequest(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(assetID); err != nil {
		http.Error(w, "Invalid asset_id format", http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(req.EmployeeID); err != nil {
		http.Error(w, "Invalid employee_id format", http.StatusBadRequest)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		status, err := dbHelper.GetAssetStatus(tx, assetID)
		if err != nil {
			return err
		}
		if status != "available" {
			return fmt.Errorf("asset is not available for assignment, current status: %s", status)
		}

		err = dbHelper.UpdateAssetForAssignment(tx, assetID, req.EmployeeID, assignerID.String())
		if err != nil {
			return err
		}

		err = dbHelper.CreateAssignmentLog(tx, assetID, req.EmployeeID, assignerID.String())
		if err != nil {
			return err
		}

		err = dbHelper.IncrementEmployeeAssetCount(tx, req.EmployeeID)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		logrus.Errorf("Transaction failed for assigning asset: %v", txErr)
		// Provide specific HTTP status codes based on the error.
		if strings.Contains(txErr.Error(), "not found") {
			http.Error(w, txErr.Error(), http.StatusNotFound)
		} else if strings.Contains(txErr.Error(), "not available") {
			http.Error(w, txErr.Error(), http.StatusConflict) // 409 Conflict is perfect for this.
		} else {
			http.Error(w, "Failed to assign asset. The operation was rolled back.", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"message": "Asset assigned successfully",
	}
	utils.EncodeResponse(w, http.StatusOK, response)
}
