package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	client   *mongo.Client
	database string
	context  context.Context
	cancel   context.CancelFunc
}

func CreateMongoClient(atlasURI string) (MongoDBClient, error) {
	clientOptions := options.Client().
		ApplyURI(atlasURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return MongoDBClient{}, err
	}

	return MongoDBClient{client: client, database: "myFirstDatabase", context: ctx, cancel: cancel}, nil
}

func (m MongoDBClient) Disconnect() {
	m.client.Disconnect(m.context)
	defer m.cancel()
	log.Println("Disconnecting Mongodb")
}

func (m MongoDBClient) SetValidators() error {

	names, err := m.client.Database(m.database).ListCollectionNames(nil, nil, nil)

	mapOfCollection := sliceToMap(names)

	if mapOfCollection["users"] != "" {
		err = m.client.Database(m.database).CreateCollection(m.context, "users")
		if err != nil {
			println(err.Error())
			return err
		}
	}

	if mapOfCollection["menuItem"] != "" {
		err = m.client.Database(m.database).CreateCollection(m.context, "menuItem")
		if err != nil {
			println(err.Error())
			return err
		}
	}

	return nil
}

func sliceToMap(slice []string) map[string]string {
	// initialize map
	elementMap := make(map[string]string)

	// put slice values into map
	for _, s := range slice {
		elementMap[s] = ""
	}

	return elementMap
}
