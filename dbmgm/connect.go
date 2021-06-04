package dbmgm

import (
	"os"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect a helper to connect tp Mongo DB
func Connect() error {
	uri := os.Getenv("MONGODB_URI_MGM")
	dbname := os.Getenv("MONGODB_NAME_MGM")

	uri = uri + "/" + dbname
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, dbname, options.Client().ApplyURI(uri))
	return err
}
