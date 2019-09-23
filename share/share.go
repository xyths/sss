package share

import (
	"context"
	"encoding/csv"
	"github.com/xyths/sss/cmd/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func AppendShare(config *utils.AppendConfig, csvfile string) error {
	file, err := os.Open(csvfile)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	log.Print(records)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Mongo.URI()))
	defer client.Disconnect(ctx)
	db := client.Database(config.Mongo.Database)
	coll := db.Collection("share")

	appendInMongo(coll, records)
	return nil
}

func appendInMongo(coll *mongo.Collection, records [][]string) error {
	log.Println("append records to mongo ...")
	return nil
}
