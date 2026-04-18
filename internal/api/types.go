package api

import "file-store/internal/file"

type StatusResponse struct {
	API   string `json:"api"`
	Store string `json:"store"`
	DB    string `json:"db"`
}

type SelfConfigResponse struct {
	Id                  string   `json:"id"`
	Name                string   `json:"name"`
	ListenAddress       string   `json:"listen_address"`
	ApiUrl              string   `json:"api_url"`
	Status              string   `json:"status"`
	Uptime              string   `json:"uptime"`
	Latency             string   `json:"latency"`
	DataIn              string   `json:"data_in"`
	DataOut             string   `json:"data_out"`
	PathTransformFunc   string   `json:"path_transform_func"`
	MessageFormat       string   `json:"message_format"`
	BaseStorageLocation string   `json:"base_storage_location"`
	BootstrapNodes      []string `json:"bootstrap_nodes"`
}

type FileType struct {
	BasePath    string          `json:"base_path"`
	KeyPath     string          `json:"key_path"`
	FullPath    string          `json:"full_path"`
	Permissions string          `json:"file_mode"`
	Size        string          `json:"size"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	SyncStatus  file.SyncStatus `json:"sync_status"`
}

type FileTypeMap map[string]FileType

type FileInfoResponse struct {
	FileMap    FileTypeMap `json:"file_map"`
	TotalFiles int         `json:"total_files"`
}
