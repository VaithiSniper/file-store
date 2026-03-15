package api

type SelfConfigResponse struct {
	ListenAddress       string   `json:"listen_address"`
	PathTransformFunc   string   `json:"path_transform_func"`
	MessageFormat       string   `json:"message_format"`
	BaseStorageLocation string   `json:"base_storage_location"`
	BootstrapNodes      []string `json:"bootstrap_nodes"`
}
