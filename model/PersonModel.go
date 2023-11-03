package model

import (
	"GoAPI/initializer"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var mongoClient *mongo.Client
var collection *mongo.Collection

type Person struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name" form:"name"`
	Major      string             `json:"major" bson:"major" form:"major"`
	Appearance bool               `json:"appearance" bson:"appearance" form:"appearance"`
	ImageURL   string             `json:"image_url" bson:"image_url"`
	CreatedAt  *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt  *time.Time         `json:"updated_at" bson:"updated_at"`
	Level      string             `json:"level" bson:"level" form:"level"`
	Office     string             `json:"office" bson:"office" form:"office"`
}

type Office struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Name    string             `json:"name" bson:"name" form:"name"`
	Address string             `json:"address" bson:"address" form:"address"`
}

type Salary struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id"`
	Level string             `json:"level" bson:"level" form:"level"`
	Value string             `json:"value" bson:"value" form:"value"`
}

func ModelList(search string) []bson.M {
	collection := initializer.ConnectDB("person_info")
	defer initializer.DisconnectDB()
	filter := bson.M{}
	if search != "" {
		//searchWithoutDiacritics := unidecode.Unidecode(search)
		filter["$or"] = []bson.M{
			bson.M{"name": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"major": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	sort := bson.M{"created_at": -1}
	options := options.Find().SetSort(sort)
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	return results
}

func ModelGet(id string) *Person {
	collection := initializer.ConnectDB("person_info")
	defer initializer.DisconnectDB()
	var person Person
	personId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": personId}
	err := collection.FindOne(context.TODO(), filter).Decode(&person)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		panic(err)
	}
	return &person
}

func ModelCreate(person Person) string {
	collection := initializer.ConnectDB("person_info")
	defer initializer.DisconnectDB()
	_, err := collection.InsertOne(context.TODO(), person)
	if err != nil {
		panic(err)
	}
	return "Created successfully"
}

func ModelDelete(id string) string {
	collection := initializer.ConnectDB("person_info")
	defer initializer.DisconnectDB()
	personId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": personId}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return "Deleted successfully"
}

func ModelUpdate(person Person) string {
	collection := initializer.ConnectDB("person_info")
	defer initializer.DisconnectDB()
	update := bson.M{"$set": bson.M{"name": person.Name, "major": person.Major, "office": person.Office, "level": person.Level, "appearance": person.Appearance, "updated_at": person.UpdatedAt}}
	filter := bson.M{"_id": person.ID}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return "Updated successfully"
}

func ModelSalaryCreate(salary Salary) string {
	collection := initializer.ConnectDB("salary_info")
	defer initializer.DisconnectDB()
	_, err := collection.InsertOne(context.TODO(), salary)
	if err != nil {
		panic(err)
	}
	return "Created successfully"
}

func ModelSalaryList() []bson.M {
	collection := initializer.ConnectDB("salary_info")
	defer initializer.DisconnectDB()
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	return results
}

func ModelDeleteSalary(id string) string {
	collection := initializer.ConnectDB("salary_info")
	defer initializer.DisconnectDB()
	salaryID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": salaryID}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return "Deleted successfully"
}

func ModelUpdateSalary(salary Salary) string {
	collection := initializer.ConnectDB("salary_info")
	defer initializer.DisconnectDB()
	update := bson.M{"$set": bson.M{"level": salary.Level, "value": salary.Value + " $"}}
	filter := bson.M{"_id": salary.ID}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return "Updated successfully"
}

func ModelOfficeCreate(office Office) string {
	collection := initializer.ConnectDB("office_info")
	defer initializer.DisconnectDB()
	_, err := collection.InsertOne(context.TODO(), office)
	if err != nil {
		panic(err)
	}
	return "Created successfully"
}

func ModelOfficeList() []bson.M {
	collection := initializer.ConnectDB("office_info")
	defer initializer.DisconnectDB()
	filter := bson.M{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	return results
}

func ModelDeleteOffice(id string) string {
	collection := initializer.ConnectDB("office_info")
	defer initializer.DisconnectDB()
	officeID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": officeID}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return "Deleted successfully"
}

func ModelUpdateOffice(office Office) string {
	collection := initializer.ConnectDB("office_info")
	defer initializer.DisconnectDB()
	update := bson.M{"$set": bson.M{"name": office.Name, "address": office.Address}}
	filter := bson.M{"_id": office.ID}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return "Updated successfully"
}
