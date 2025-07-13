package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func calculateAuctionEndTime() time.Time {
	durationStr := os.Getenv("AUCTION_DURATION_MINUTES")
	durationMinutes, err := strconv.Atoi(durationStr)
	if err != nil {
		log.Printf("Erro ao converter AUCTION_DURATION_MINUTES: %v. Usando valor padrão de 10 minutos.", err)
		durationMinutes = 10
	}
	return time.Now().Add(time.Duration(durationMinutes) * time.Minute)
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   calculateAuctionEndTime().Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

// Fecha todos os leilões com status "Active" e timestamp menor que o tempo atual
func (ar *AuctionRepository) CloseExpiredAuctions(ctx context.Context) error {
	now := time.Now().Unix()

	filter := bson.M{
		"status":    auction_entity.Active,
		"timestamp": bson.M{"$lt": now},
	}

	update := bson.M{
		"$set": bson.M{"status": auction_entity.Completed},
	}

	result, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Erro ao fechar leilões expirados", err)
		return err
	}

	log.Printf("Leilões fechados automaticamente: %d", result.ModifiedCount)
	return nil
}
