package receive

// Config describes the configuration of the client
type Config struct {
	Address  string
	Insecure bool
	Secret   string
	Force    bool
}
