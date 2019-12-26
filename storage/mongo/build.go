package mongo

import "github.com/simplesurance/baur/storage"

// GetLatestBuildByDigest returns the build id of a build for the application
// with the passed digest. If multiple builds exist, the one with the latest
// stop_timestamp is returned.
// Inputs are not fetched from the database.
// If no builds exist storage.ErrNotExist is returned
func (c *Client) GetLatestBuildByDigest(appName, totalInputDigest string) (*storage.BuildWithDuration, error) {
	return &storage.BuildWithDuration{}, nil
}

// GetBuildOutputs returns build outputs
func (c *Client) GetBuildOutputs(buildID int) ([]*storage.Output, error) {
	return nil, nil
}

// BuildExist returns true if the build with the given ID exist.
func (c *Client) BuildExist(id int) (bool, error) {
	return false, nil
}

// GetBuildWithoutInputsOutputs returns a single build, if no build with the ID
// exist ErrNotExist is returned
func (c *Client) GetBuildWithoutInputsOutputs(id int) (*storage.BuildWithDuration, error) {
	return nil, nil
}

// GetSameTotalInputDigestsForAppBuilds finds TotalInputDigests that are the
// same for builds of an app with a build start time not before startTs
// If not builds with the same totalInputDigest is found, an empty slice is
// returned.
func (c *Client) GetBuildsWithoutInputsOutputs(filters []*storage.Filter, sorters []*storage.Sorter) ([]*storage.BuildWithDuration, error) {
	return nil, nil
}
