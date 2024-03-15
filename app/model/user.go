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
	BaseModel `bson:",inline"`
	Username  string `bson:"username" json:"username"` // username
	Timezone  string `bson:"timezone" json:"timezone"` // timezone
	Status    string `bson:"status" json:"status"`     // status, NORMAL=正常, BAN=屏蔽, DELETED=删除
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
				{"username", 1},
			},
		},
	)
	if err != nil {
		log.Warn(err)
	}
}

// create user
func (d *User) Create(username string, timezone string) (*User, error) {
	doc := &User{}
	doc.Username = username
	doc.Timezone = timezone

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

// get user by username
func (d *User) GetByUsername(username string) (*User, error) {
	doc := &User{}
	err := d.GetCollection().Where(bson.D{{"username", username}}).Find(doc)
	return doc, err
}

// check username is available
func (d *User) CheckUsernameAvailable(username string) bool {
	_, err := d.GetByUsername(username)
	return err == nil
}
