package model

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/mylukin/EchoPilot/storage/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // ID 对应 request ID
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`        // 创建时间
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`        // 更新时间
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
