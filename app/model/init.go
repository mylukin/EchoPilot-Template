package model

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/mylukin/EchoPilot/storage/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type JSON map[string]interface{}

// Error
func (body JSON) Error() string {
	data, err := json.Marshal(body)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// Error is error
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

// Result is send message
type Result struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

// result error
func (body Result) Error() string {
	data, err := json.Marshal(body)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // ID 对应 request ID
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`      // 创建时间
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`      // 更新时间
}

// Error is error
func (m BaseModel) Error() string {
	data, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func GetCollection(tableName string) *mongo.Collection {
	return mongo.C(tableName)
}
