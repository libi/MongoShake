package utils

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

const (
	testCollection = "test"
)

func TestCommonFunctions(t *testing.T) {
	var nr int

	InitialLogger("", "", "debug", true, 1)

	{
		fmt.Printf("TestCommonFunctions case %d.\n", nr)
		nr++

		conn, err := NewMongoCommunityConn(testMongoAddress, VarMongoConnectModeSecondaryPreferred, true,
			ReadWriteConcernDefault, ReadWriteConcernDefault, "")
		assert.Equal(t, err, nil, "")

		ok, err := GetAndCompareVersion(conn, "3.4.0", "")
		assert.Equal(t, err, nil, "")
		assert.Equal(t, ok, true, "")

		conn.Close()
	}

	{
		fmt.Printf("TestCommonFunctions case %d.\n", nr)
		nr++

		conn, err := NewMongoCommunityConn(testMongoAddress, VarMongoConnectModeSecondaryPreferred, true,
			ReadWriteConcernDefault, ReadWriteConcernDefault, "")
		assert.Equal(t, err, nil, "")

		ok := conn.IsGood()
		assert.Equal(t, ok, true, "")

		ok = conn.HasOplogNs()
		assert.Equal(t, ok, true, "")

		name := conn.AcquireReplicaSetName()
		fmt.Printf("ReplicaSetName:%v\n", name)
		assert.NotEqual(t, name, nil, "")

		databases, err := conn.Client.ListDatabaseNames(nil, bson.M{})
		assert.Equal(t, nil, err, "should be equal")
		for _, db := range databases {
			if db != "admin" && db != "local" && db != "config" {
				fmt.Printf("delete database:%v\n", db)
				err = conn.Client.Database(db).Drop(nil)
				assert.Equal(t, nil, err, "should be equal")

			}
		}

		uniqueOk := conn.HasUniqueIndex()
		assert.Equal(t, uniqueOk, false, "")

		_, err = conn.Client.Database(testDb).Collection(testCollection).InsertOne(context.Background(),
			bson.D{{"x", 11}, {"y", 12}})
		assert.Equal(t, err, nil, "")
		_, err = conn.Client.Database(testDb).Collection(testCollection).InsertOne(context.Background(),
			bson.D{{"x", 21}, {"y", 22}})
		assert.Equal(t, err, nil, "")

		// create unique index
		indexOptions := options.Index().SetUnique(true)
		str, err := conn.Client.Database(testDb).Collection(testCollection).Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys:    bson.D{{"x", 1}},
				Options: indexOptions,
			})
		assert.Equal(t, err, nil, "")
		fmt.Printf("Create index:%v\n", str)

		uniqueOk = conn.HasUniqueIndex()
		assert.Equal(t, uniqueOk, true, "")
	}
}
