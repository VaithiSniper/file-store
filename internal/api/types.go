package api

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
