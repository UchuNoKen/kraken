package hdfsbackend

import (
	"fmt"
	"io"

	"code.uber.internal/infra/kraken/lib/backend/nameparse"
)

// DockerBlobClient is an HDFS client for uploading / download blobs to a docker
// registry.
type DockerBlobClient struct {
	config Config
	client *client
}

// NewDockerBlobClient creates a new DockerBlobClient.
func NewDockerBlobClient(config Config) (*DockerBlobClient, error) {
	config, err := config.applyDefaults()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}
	return &DockerBlobClient{client: newClient(config), config: config}, nil
}

func (c *DockerBlobClient) path(name string) (string, error) {
	return nameparse.ShardDigestPath(c.config.RootDirectory, name)
}

// Download downloads a blob for name into dst. name should be a sha256 digest
// of the desired blob.
func (c *DockerBlobClient) Download(name string, dst io.Writer) error {
	path, err := c.path(name)
	if err != nil {
		return err
	}
	return c.client.download(path, dst)
}

// Upload uploads src to name. name should be a sha256 digest of src.
func (c *DockerBlobClient) Upload(name string, src io.Reader) error {
	path, err := c.path(name)
	if err != nil {
		return err
	}
	return c.client.upload(path, src)
}
