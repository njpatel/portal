package send

import "fmt"

// Config describes the configuration of the client
type Config struct {
	Address  string
	Insecure bool
	Secret   string
}

// Run starts a new client for one-shot sending
func Run(c *Config, args []string) {
	fmt.Printf("Starting send on %s with %t and %s for %s\n", c.Address, c.Insecure, c.Secret, args)
}

// RunSync starts a new client for sync sending
func RunSync(c *Config, args []string) {
	fmt.Printf("Starting sync on %s with %t and %s for %s\n", c.Address, c.Insecure, c.Secret, args)
}
