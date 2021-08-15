package mongodb

import (
	"mongoTest.io/server"
)

type ObjectID [12]byte

func (m MongoDBClient) CreateUser(user server.User) (string, error) {
	userCollection := m.client.Database(m.database).Collection("users")

	userResult, err := userCollection.InsertOne(nil, user)
	if err != nil {
		return "", err
	}

	oid, err := insertOneID(userResult)
	if err != nil {
		return oid, err
	}

	return oid, nil
}
