package terraform

// Client represents a Terraform client with a working directory.
type Client struct {
	workingDir string
}

// NewClient creates a new Terraform client with the specified working directory.
func NewClient(workingDir string) *Client {
	return &Client{workingDir: workingDir}
}