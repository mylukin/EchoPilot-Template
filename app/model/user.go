package model

import (
	"time"

	"github.com/labstack/gommon/log"
	"github.com/mylukin/EchoPilot/helper"
	"github.com/mylukin/EchoPilot/storage/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User
type User struct {
	BaseModel   `bson:",inline"`
	MainAddress string `bson:"mainAddress" json:"mainAddress"`               // main address，根据助记词生成ETH地址
	Timezone    string `bson:"timezone" json:"timezone"`                     // timezone
	Status      string `bson:"status" json:"status"`                         // status, NORMAL=正常, BAN=屏蔽, DELETED=删除
	Settings    bson.M `bson:"settings,omitempty" json:"settings,omitempty"` // settings
}

func (d *User) CollectionName() string {
	return "user"
}

func (d *User) GetCollection() *mongo.Collection {
	return GetCollection(d.CollectionName())
}

// new user object
func UserTable() *User {
	return &User{}
}

func init() {
	_, err := UserTable().GetCollection().Index(
		bson.M{
			"unique": true,
			"keys": bson.D{
				{"mainAddress", 1},
			},
		},
	)
	if err != nil {
		log.Warn(err)
	}
}

// create user
func (d *User) Create(mainAddress string, timezone string, status string) (*User, error) {
	doc := &User{}
	doc.MainAddress = mainAddress
	doc.Timezone = timezone
	doc.Status = status

	if doc.Timezone == "" {
		doc.Timezone = helper.Config("TZ")
	}
	if doc.Status == "" {
		doc.Status = "NORMAL"
	}
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = time.Now()
	}
	if doc.UpdatedAt.IsZero() {
		doc.UpdatedAt = time.Now()
	}
	res, err := d.GetCollection().Insert(doc)
	if err != nil {
		return doc, err
	}
	doc.ID = res.InsertedID.(primitive.ObjectID)
	return doc, nil
}

// get user by id
func (d *User) GetByID(id primitive.ObjectID) (*User, error) {
	doc := &User{}
	err := d.GetCollection().FindByID(id, doc)
	return doc, err
}

// get user by main address
func (d *User) GetByMainAddress(mainAddress string) (*User, error) {
	doc := &User{}
	err := d.GetCollection().Where(bson.D{{"mainAddress", mainAddress}}).Find(doc)
	return doc, err
}

// check main address is available
func (d *User) CheckMainAddressAvailable(mainAddress string) bool {
	_, err := d.GetByMainAddress(mainAddress)
	return err == nil
}

// get user settings by id
func (d *User) GetSettingsByID(id primitive.ObjectID) (bson.M, error) {
	doc := &User{}
	err := d.GetCollection().FindByID(id, doc)
	return doc.Settings, err
}

// save user settings by id
func (d *User) SaveSettingsByID(id primitive.ObjectID, settings bson.M) (*mongo.UpdateResult, error) {
	return d.GetCollection().Where(bson.D{{"_id", id}}).UpdateOne(bson.M{"$set": bson.M{
		"settings":  settings,
		"updatedAt": time.Now(),
	}})
}
