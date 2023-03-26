package shopify

type Configuration struct {
	// Name is the name of the Shopify store as shown in
	// <store-name>.myshopify.com.
	Name string `json:"name"`
	// APIKey is the API key for the Shopify store.
	APIKey string `json:"apiKey"`
	// Password is the password associated with the API key.
	Password string `json:"password"`
}

type Configurations []Configuration

// Get returns the configuration for the given name. The first configuration is
// returned if name is an empty string. Nil is returned if no configuration can
// be found for the given name.
func (configs Configurations) Get(name string) *Configuration {
	if len(configs) < 1 {
		return nil
	}
	if name == "" {
		return &configs[0]
	}
	for _, c := range configs {
		if c.Name == name {
			return &c
		}
	}
	return nil
}
