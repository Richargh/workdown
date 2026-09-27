package kernel

type WorkItem struct {
	Remote   string            `json:"remote"`
	Provider string            `json:"provider"`
	ID       string            `json:"id"`
	Key      string            `json:"key"`
	URL      string            `json:"url"`
	Title    string            `json:"title"`
	Fields   []Field           `json:"fields,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type Field struct {
	Name        string `json:"name"`
	ProviderKey string `json:"providerKey"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Editable    bool   `json:"editable"`
}
