package mongodb

import (
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
	"mongoTest.io/server"
)

func (m MongoDBClient) CreateMenuItem(menuItem server.MenuItem) (string, error) {
	menuItemCollection := m.client.Database(m.database).Collection("menuItem")

	result, err := menuItemCollection.InsertOne(nil, menuItem)
	if err != nil {
		return "", err
	}

	oid, err := insertOneID(result)
	if err != nil {
		return oid, err
	}

	return oid, nil
}

func insertOneID(fuc *mongo.InsertOneResult) (string, error) {
	if oid, ok := fuc.InsertedID.(primitive.ObjectID); ok {
		return oid.String(), nil
	}

	return "", errors.New("can't Get ObjectID of menu Item")
}

func (m MongoDBClient) AllMenuItems() ([]server.MenuItem, error) {
	menuItemCollection := m.client.Database(m.database).Collection("menuItem")
	cursor, err := menuItemCollection.Find(nil, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var bsonItem []bson.M
	if err = cursor.All(nil, &bsonItem); err != nil {
		log.Fatal(err)
	}

	menuItems := m.menuBSONToMenuItem(bsonItem)

	return menuItems, nil
}

func (m MongoDBClient) menuBSONToMenuItem(bsonItem []bson.M) []server.MenuItem {
	menuItems := []server.MenuItem{}
	for _, item := range bsonItem {
		menuItem := server.MenuItem{}
		menuItem.ImageSource = item["imagesource"].(string)
		menuItem.ItemDescription = item["itemdescription"].(string)
		menuItem.ItemName = item["itemname"].(string)
		menuItem.ItemIngredients = item["itemingredients"].(string)
		menuItem.ItemPrice = item["itemprice"].(string)
		menuItem.ID = item["_id"].(primitive.ObjectID).Hex()

		menuItems = append(menuItems, menuItem)
	}
	return menuItems
}

func (m MongoDBClient) MenuItemFromID(id string) (server.MenuItem, error) {
	var bsonItem bson.M
	menuItemCollection := m.client.Database(m.database).Collection("menuItem")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return server.MenuItem{}, err
	}

	err = menuItemCollection.FindOne(nil, bson.M{"_id": docID}).Decode(&bsonItem)
	if err != nil {
		return server.MenuItem{}, err
	}

	bsonItems := []bson.M{}
	bsonItems = append(bsonItems, bsonItem)

	menuItem := m.menuBSONToMenuItem(bsonItems)

	return menuItem[0], nil
}
