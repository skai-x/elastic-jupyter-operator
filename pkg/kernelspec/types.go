package kernelspec

type kernelConfig struct {
	Language    string   `json:"language,omitempty"`
	DisplayName string   `json:"display_name,omitempty"`
	Metadata    metadata `json:"metadata,omitempty"`
	Argv        []string `json:"argv,omitempty"`
}

type metadata struct {
	ProcessProxy processProxy `json:"process_proxy,omitempty"`
	Config       config       `json:"config,omitempty"`
}

type processProxy struct {
	ClassName string `json:"class_name,omitempty"`
}

type config struct {
	ImageName string `json:"image_name,omitempty"`
}
