package resolver

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func CursorClose(cur *mongo.Cursor, ctx context.Context) {
	err := cur.Close(ctx)
	if err != nil {
		log.Print(err)
	}
}
