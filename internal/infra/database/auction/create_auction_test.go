package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCloseExpiredAuctions_WithRealMongo(t *testing.T) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/?authSource=admin"))
	assert.NoError(t, err)

	db := client.Database("auctions_test")
	collection := db.Collection("auctions")

	err = collection.Drop(ctx)
	assert.NoError(t, err)

	repo := &AuctionRepository{Collection: collection}

	// Cria um leilão expirado
	expiredAuction := AuctionEntityMongo{
		Id:          "auction123",
		ProductName: "Produto Teste",
		Category:    "Categoria",
		Description: "Descrição do produto",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now().Add(-1 * time.Minute).Unix(),
	}

	_, err = collection.InsertOne(ctx, expiredAuction)
	assert.NoError(t, err)

	err = repo.CloseExpiredAuctions(ctx)
	assert.NoError(t, err)

	var updated AuctionEntityMongo
	err = collection.FindOne(ctx, bson.M{"_id": "auction123"}).Decode(&updated)
	assert.NoError(t, err)
	assert.Equal(t, auction_entity.Completed, updated.Status)
}
