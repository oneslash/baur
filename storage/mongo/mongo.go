package mongo

import (
	"compress/zlib"
	"context"
	"github.com/simplesurance/baur/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"time"
)

const connectionTimeout = 5 * time.Second

// Client is a mongodb (can be used for aws documentdb) storage client
type Client struct {
	Db *mongo.Database
}

// New establishes a connection a mongodb/documentdb
func New(url string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(url),
		options.Client().SetRetryWrites(true),
		options.Client().SetReadPreference(readpref.SecondaryPreferred()),
		options.Client().SetZlibLevel(zlib.BestSpeed))
	if err != nil {
		return nil, err
	}

	// skipping error since it has already been parsed in mongo.Connect
	connStr, _ := connstring.Parse(url)

	return &Client{
		Db: conn.Database(connStr.Database),
	}, nil
}

// GetApps returns all application records ordered by Name
func (c *Client) GetApps() ([]*storage.Application, error) {
	return nil, nil
}

// GetSameTotalInputDigestsForAppBuilds finds TotalInputDigests that are the
// same for builds of an app with a build start time not before startTs
// If not builds with the same totalInputDigest is found, an empty slice is
// returned.
func (c *Client) GetSameTotalInputDigestsForAppBuilds(appName string, startTs time.Time) (map[string][]int, error) {
	return nil, nil
}
