package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"math"
	"net/http"
	"storex/database"
	"storex/database/dbHelper"
	"storex/middlewares"
	"storex/models"
	"storex/utils"
	"strconv"
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
	if err := utils.EncodeResponse(w, http.StatusOK, response); err != nil {
		logrus.Error(err)
	}
}

func GetAssets(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	filters := make(map[string]string)

	// 1. Parse and validate filters
	if q := queryParams.Get("q"); q != "" {
		filters["q"] = q
	}

	if assetType := strings.ToLower(queryParams.Get("type")); assetType != "" {
		if utils.IsValidAssetType(assetType) {
			filters["type"] = assetType
		}
	}

	if status := strings.ToLower(queryParams.Get("status")); status != "" {
		if utils.IsValidAssetStatus(status) {
			filters["status"] = status
		}
	}

	// Correctly handle the case-sensitive 'owned_by' filter
	if ownedByInput := strings.ToLower(queryParams.Get("owned_by")); ownedByInput != "" {
		if dbOwnerValue, isValid := utils.GetValidAssetOwner(ownedByInput); isValid {
			filters["owned_by"] = dbOwnerValue
		}
	}

	// 2. Sanitize pagination parameters
	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(queryParams.Get("pageSize"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 3. Call database helper
	assets, totalRecords, err := dbHelper.GetAssetsFiltered(filters, page, pageSize)
	if err != nil {
		logrus.Errorf("Failed to get assets from database: %v", err)
		http.Error(w, "Failed to retrieve asset data", http.StatusInternalServerError)
		return
	}

	// 4. Calculate pagination metadata
	totalPages := 0
	if totalRecords > 0 {
		totalPages = int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	}

	// 5. Assemble and send response
	response := models.GetAssetsResponse{
		Data: assets,
		Pagination: models.PaginationInfo{
			CurrentPage:  page,
			PageSize:     pageSize,
			TotalRecords: totalRecords,
			TotalPages:   totalPages,
		},
	}
	if err := utils.EncodeResponse(w, http.StatusOK, response); err != nil {
		logrus.Error(err)
	}
}

func UnassignAsset(w http.ResponseWriter, r *http.Request) {
	// 1. Get IDs from context and URL.
	vars := mux.Vars(r)
	assetID := vars["asset_id"]

	retrieverIDVal := r.Context().Value(middlewares.UserIDKey)
	if retrieverIDVal == nil {
		http.Error(w, "Unauthorized: could not identify user", http.StatusUnauthorized)
		return
	}
	retrieverID, _ := retrieverIDVal.(uuid.UUID)

	// 2. Decode request body.
	var req models.UnassignAssetRequest
	if err := utils.DecodeRequest(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Reason == "" {
		http.Error(w, "reason_of_retrieval is a required field", http.StatusBadRequest)
		return
	}

	// 3. Start the database transaction.
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		// A. Check if the asset is actually assigned and get the employee ID.
		assetInfo, err := dbHelper.GetAssetForUnassignment(tx, assetID)
		if err != nil {
			return err // Will be handled as 404 or 500 below.
		}
		if assetInfo.Status != "assigned" {
			return fmt.Errorf("asset is not currently assigned (status: %s)", assetInfo.Status)
		}
		if assetInfo.AssignedTo == "" {
			return fmt.Errorf("asset is marked as assigned but has no employee ID linked")
		}

		// B. Close the active assignment log entry, now passing the retriever's ID.
		err = dbHelper.CloseAssignmentLog(tx, assetID, req.Reason, retrieverID.String())
		if err != nil {
			return err
		}

		// C. Decrement the employee's asset count.
		err = dbHelper.DecrementEmployeeAssetCount(tx, assetInfo.AssignedTo)
		if err != nil {
			return err
		}

		// D. Finally, update the asset itself to be available.
		err = dbHelper.UpdateAssetForUnassignment(tx, assetID, retrieverID.String())
		if err != nil {
			return err
		}

		return nil // Success!
	})

	// 4. Handle the final result of the transaction.
	if txErr != nil {
		logrus.Errorf("Transaction failed for un-assigning asset: %v", txErr)
		if strings.Contains(txErr.Error(), "not found") {
			http.Error(w, txErr.Error(), http.StatusNotFound)
		} else if strings.Contains(txErr.Error(), "not currently assigned") || strings.Contains(txErr.Error(), "no active assignment log") {
			http.Error(w, txErr.Error(), http.StatusConflict) // 409 Conflict is appropriate here.
		} else {
			http.Error(w, "Failed to un-assign asset. The operation was rolled back.", http.StatusInternalServerError)
		}
		return
	}

	// 5. Send success response.
	response := map[string]string{
		"message": "Asset un-assigned successfully and is now available",
	}
	if err := utils.EncodeResponse(w, http.StatusOK, response); err != nil {
		logrus.Error(err)
	}
}
