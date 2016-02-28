package receive

import "fmt"

// Config describes the configuration of the client
type Config struct {
	Address  string
	Insecure bool
	Secret   string
	Force    bool
}

// Run starts a receiving client
func Run(c *Config, token string, outputDir string) {
	fmt.Printf("Receiving on %s/%t/%s T:%s O:%s\n", c.Address, c.Insecure, c.Secret, token, outputDir)
}
