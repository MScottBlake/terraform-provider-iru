package client

// Blueprint represents a Kandji Blueprint.
type Blueprint struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Icon           string `json:"icon,omitempty"`
	Color          string `json:"color,omitempty"`
	EnrollmentCode string `json:"enrollment_code,omitempty"`
}

// ADEIntegration represents a Kandji ADE Integration.
type ADEIntegration struct {
	ID          string     `json:"id,omitempty"`
	BlueprintID string     `json:"blueprint_id,omitempty"` // Used for Update input
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	Blueprint   *Blueprint `json:"blueprint,omitempty"` // From response
}

// Device represents a Kandji Device.
type Device struct {
	ID                string `json:"device_id,omitempty"` // Note: API might return device_id or id. Check docs. Docs said id. But commonly it's device_id in list? I'll use ID alias if needed.
	DeviceName        string `json:"device_name,omitempty"`
	AssetTag          string `json:"asset_tag,omitempty"`
	SerialNumber      string `json:"serial_number,omitempty"`
	Model             string `json:"model,omitempty"`
	OSVersion         string `json:"os_version,omitempty"`
	BlueprintID       string `json:"blueprint_id,omitempty"`
	UserID            string `json:"user_id,omitempty"`
	Platform          string `json:"platform,omitempty"`
	LastCheckIn       string `json:"last_check_in,omitempty"`
}

// Tag represents a Kandji Tag.
type Tag struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// CustomScript represents a Kandji Custom Script library item.
type CustomScript struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name"`
	Active             bool   `json:"active"`
	ExecutionFrequency string `json:"execution_frequency"`
	Restart            bool   `json:"restart"`
	Script             string `json:"script"`
	RemediationScript  string `json:"remediation_script,omitempty"`
	ShowInSelfService  bool   `json:"show_in_self_service"`
}

// CustomProfile represents a Kandji Custom Profile library item.
type CustomProfile struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name"`
	Active        bool   `json:"active"`
	Profile       string `json:"profile,omitempty"`
	MDMIdentifier string `json:"mdm_identifier,omitempty"`
	RunsOnMac     bool   `json:"runs_on_mac"`
	RunsOnIPhone  bool   `json:"runs_on_iphone"`
	RunsOnIPad    bool   `json:"runs_on_ipad"`
	RunsOnTV      bool   `json:"runs_on_tv"`
	RunsOnVision  bool   `json:"runs_on_vision"`
}

// User represents a Kandji User.
type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	IsArchived bool   `json:"is_archived"`
}

// PrismFileVault represents FileVault status information from Prism.
type PrismFileVault struct {
	DeviceID     string `json:"device_id"`
	DeviceName   string `json:"device__name"`
	SerialNumber string `json:"serial_number"`
	Status       bool   `json:"status"`
	KeyType      string `json:"key_type"`
	KeyEscrowed  bool   `json:"key_escrowed"`
}

// PrismAppFirewall represents Application Firewall status information from Prism.
type PrismAppFirewall struct {
	DeviceID               string `json:"device_id"`
	DeviceName             string `json:"device__name"`
	SerialNumber           string `json:"serial_number"`
	Status                 bool   `json:"status"`
	BlockAllIncoming       bool   `json:"block_all_incoming"`
	StealthMode            bool   `json:"stealth_mode"`
	AllowSignedApplications bool   `json:"allow_signed_applications"`
}

// PrismApp represents an application installed on a device from Prism.
type PrismApp struct {
	DeviceID     string `json:"device_id"`
	DeviceName   string `json:"device__name"`
	SerialNumber string `json:"serial_number"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	BundleID     string `json:"bundle_id"`
	Path         string `json:"path"`
}

// PrismEntry represents a generic entry from a Prism endpoint.
type PrismEntry map[string]interface{}

// Vulnerability represents a vulnerability from Vulnerability Management.
type Vulnerability struct {
	CVEID              string   `json:"cve_id"`
	Severity           string   `json:"severity"`
	CVSSScore          float64  `json:"cvss_score"`
	FirstDetectionDate string   `json:"first_detection_date"`
	DeviceCount        int      `json:"device_count"`
	Status             string   `json:"status"`
	Software           []string `json:"software"`
}

// CustomApp represents a Kandji Custom App library item.
type CustomApp struct {
	ID                     string `json:"id,omitempty"`
	Name                   string `json:"name"`
	FileKey                string `json:"file_key"`
	InstallType            string `json:"install_type"`
	InstallEnforcement     string `json:"install_enforcement"`
	UnzipLocation          string `json:"unzip_location,omitempty"`
	AuditScript            string `json:"audit_script,omitempty"`
	PreinstallScript       string `json:"preinstall_script,omitempty"`
	PostinstallScript      string `json:"postinstall_script,omitempty"`
	ShowInSelfService      bool   `json:"show_in_self_service"`
	SelfServiceCategoryID  string `json:"self_service_category_id,omitempty"`
	SelfServiceRecommended bool   `json:"self_service_recommended"`
	Active                 bool   `json:"active"`
	Restart                bool   `json:"restart"`
}

// InHouseApp represents a Kandji In-House App library item (.ipa).
type InHouseApp struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name"`
	FileKey      string `json:"file_key"`
	RunsOnIPhone bool   `json:"runs_on_iphone"`
	RunsOnIPad   bool   `json:"runs_on_ipad"`
	RunsOnTV     bool   `json:"runs_on_tv"`
	Active       bool   `json:"active"`
}

// AuditEvent represents an audit log event.
type AuditEvent struct {
	ID              string      `json:"id"`
	Action          string      `json:"action"`
	OccurredAt      string      `json:"occurred_at"`
	ActorID         string      `json:"actor_id"`
	ActorType       string      `json:"actor_type"`
	TargetID        string      `json:"target_id"`
	TargetType      string      `json:"target_type"`
	TargetComponent string      `json:"target_component"`
	NewState        interface{} `json:"new_state"`
}

// Licensing represents tenant licensing information.
type Licensing struct {
	Counts struct {
		ComputersCount int `json:"computers_count"`
		IOSCount       int `json:"ios_count"`
		IPadOSCount    int `json:"ipados_count"`
		MacOSCount     int `json:"macos_count"`
		TVOSCount      int `json:"tvos_count"`
	} `json:"counts"`
	Limits struct {
		PlanType    string `json:"plan_type"`
		MaxDevices  int    `json:"max_devices"`
	} `json:"limits"`
	TenantOverLicenseLimit bool `json:"tenantOverLicenseLimit"`
}

// Threat represents a detected malware/pup threat.
type Threat struct {
	ThreatName         string `json:"threat_name"`
	Classification     string `json:"classification"`
	Status             string `json:"status"`
	DeviceName         string `json:"device_name"`
	DeviceID           string `json:"device_id"`
	DetectionDate      string `json:"detection_date"`
	FilePath           string `json:"file_path"`
	FileHash           string `json:"file_hash"`
	DeviceSerialNumber string `json:"device_serial_number"`
}

// BehavioralDetection represents a behavioral detection event.
type BehavioralDetection struct {
	ID             string `json:"id"`
	ThreatID       string `json:"threat_id"`
	Description    string `json:"description"`
	Classification string `json:"classification"`
	DetectionDate  string `json:"detection_date"`
	ThreatStatus   string `json:"threat_status"`
	DeviceInfo     struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		SerialNumber string `json:"serial_number"`
	} `json:"device_info"`
}

// DeviceSecretsALBC represents Activation Lock Bypass Code response.
type DeviceSecretsALBC struct {
	UserBasedALBC   string `json:"user_based_albc"`
	DeviceBasedALBC string `json:"device_based_albc"`
}

// DeviceSecretsFileVault represents FileVault Recovery Key response.
type DeviceSecretsFileVault struct {
	Key string `json:"key"`
}

// DeviceSecretsUnlockPin represents Unlock Pin response.
type DeviceSecretsUnlockPin struct {
	Pin string `json:"pin"`
}

// DeviceSecretsRecoveryLock represents Recovery Lock Password response.
type DeviceSecretsRecoveryLock struct {
	RecoveryPassword string `json:"recovery_password"`
}

// BlueprintLibraryItem represents a library item assigned to a blueprint.
type BlueprintLibraryItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
