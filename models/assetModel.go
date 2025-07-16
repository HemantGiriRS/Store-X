package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type BaseAsset struct {
	Brand         string     `json:"brand" db:"brand"`
	Model         string     `json:"model" db:"model"`
	Type          string     `json:"type" db:"type"`
	SerialNo      string     `json:"serial_no" db:"serial_no"`
	OwnedBy       string     `json:"owned_by" db:"owned_by"`
	PurchaseDate  *time.Time `json:"purchase_date,omitempty" db:"purchase_date"`
	WarrantyStart *time.Time `json:"warranty_start,omitempty" db:"warranty_start"`
	WarrantyEnd   *time.Time `json:"warranty_end,omitempty" db:"warranty_end"`
}

type CreateAssetRequest struct {
	BaseAsset
	Specs json.RawMessage `json:"specs"`
}

type LaptopSpecs struct {
	Processor          string `json:"processor" db:"processor"`
	RAM_GB             int    `json:"ram_gb" db:"ram_gb"`
	Storage_GB         int    `json:"storage_gb" db:"storage_gb"`
	OS                 string `json:"os" db:"os"`
	ScreenSizeInch     int    `json:"screen_size_inch" db:"screen_size_inch"`
	BatteryBackupHours int    `json:"battery_backup_hours" db:"battery_backup_hours"`
}

type MouseSpecs struct {
	ConnectionType  string `json:"connection_type" db:"connection_type"`
	DPI             int    `json:"dpi" db:"dpi"`
	NumberOfButtons int    `json:"number_of_buttons" db:"number_of_buttons"`
}

type MonitorSpecs struct {
	SizeInch      int    `json:"size_inch" db:"size_inch"`
	Resolution    string `json:"resolution" db:"resolution"`
	RefreshRateHZ int    `json:"refresh_rate_hz" db:"refresh_rate_hz"`
	PanelType     string `json:"panel_type" db:"panel_type"`
}

type HardDiskSpecs struct {
	CapacityGB int    `json:"capacity_gb" db:"capacity_gb"`
	Type       string `json:"type" db:"type"`
	Interface  string `json:"interface" db:"interface"`
}

type PenDriveSpecs struct {
	CapacityGB int    `json:"capacity_gb" db:"capacity_gb"`
	USBType    string `json:"usb_type" db:"usb_type"`
}

type MobileSpecs struct {
	IMEI               string `json:"imei" db:"imei"`
	RAM_GB             int    `json:"ram_gb" db:"ram_gb"`
	Storage_GB         int    `json:"storage_gb" db:"storage_gb"`
	OS                 string `json:"os" db:"os"`
	ScreenSizeInch     int    `json:"screen_size_inch" db:"screen_size_inch"`
	BatteryCapacityMAH int    `json:"battery_capacity_mah" db:"battery_capacity_mah"`
}

type SimSpecs struct {
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Operator    string `json:"operator" db:"operator"`
	SimType     string `json:"sim_type" db:"sim_type"`
	PlanDetails string `json:"plan_details" db:"plan_details"`
}

type AccessoriesSpecs struct {
	Name          string `json:"name" db:"name"`
	Description   string `json:"description" db:"description"`
	Compatibility string `json:"compatibility" db:"compatibility"`
}

type AssignAssetRequest struct {
	EmployeeID string `json:"employee_id"`
}

type AssetDetail struct {
	ID              string         `json:"id" db:"id"`
	Brand           string         `json:"brand" db:"brand"`
	Model           string         `json:"model" db:"model"`
	Type            string         `json:"type" db:"type"`
	SerialNo        string         `json:"serial_no" db:"serial_no"`
	Status          string         `json:"status" db:"status"`
	OwnedBy         string         `json:"owned_by" db:"owned_by"`
	PurchaseDate    *time.Time     `json:"purchase_date" db:"purchase_date"`
	WarrantyEnd     *time.Time     `json:"warranty_end" db:"warranty_end"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	AssignedToID    sql.NullString `json:"assigned_to_id" db:"assigned_to_id"`
	AssignedToName  sql.NullString `json:"assigned_to_name" db:"assigned_to_name"`
	AssignedToEmail sql.NullString `json:"assigned_to_email" db:"assigned_to_email"`
}

type GetAssetsResponse struct {
	Data       []AssetDetail  `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

type UnassignAssetRequest struct {
	Reason string `json:"reason_of_retrieval"`
}

type AssetAssignmentInfo struct {
	Status     string `db:"status"`
	AssignedTo string `db:"assigned_to"`
}

type DeleteAssetRequest struct {
	Reason string `json:"reason"`
}

type SendForServiceRequest struct {
	AssignedTo  string `json:"assigned_to"`
	Price       string `json:"price"`
	Description string `json:"description"`
}

type ReceiveFromServiceRequest struct {
	Notes string `json:"notes"`
}

type TimelineEvent struct {
	ID          string    `json:"id" db:"id"`
	EventType   string    `json:"event_type" db:"event_type"`
	EventDate   time.Time `json:"event_date" db:"event_date"`
	Description string    `json:"description" db:"description"`
	ActorName   string    `json:"actor_name" db:"actor_name"`
}

type GetAssetTimelineResponse struct {
	Data       []TimelineEvent `json:"data"`
	Pagination PaginationInfo  `json:"pagination"`
}
