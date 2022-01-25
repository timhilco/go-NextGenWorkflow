package databases

import (
	"context"
	"fmt"

	"github.com/timhilco/go-NextGenWorkflow/domain"
	"github.com/timhilco/go-NextGenWorkflow/util/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//PersonBusinessProcessDB is
type PersonBusinessProcessDB struct {
	MongoClient *mongo.Client
	logger      *logger.HilcoLogger
}

//CreatePersonBusinessDB is
func CreatePersonBusinessDB(ctx context.Context, processingContext DatabaseContext) PersonBusinessProcessDB {
	pbp := PersonBusinessProcessDB{}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	logger := processingContext.Logger
	pbp.logger = logger
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	pbp.MongoClient = client

	if err != nil {
		logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB: %s", err))
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB -> Error: %s", err))
	}

	logger.Info("PersonBusinessProcessDB -> Connected to MongoDB!")
	return pbp
}

//InsertPersonBusinessProcessDocument inserts documents into the PersonBusinessProcess Collection
func (p *PersonBusinessProcessDB) InsertPersonBusinessProcessDocument(ctx context.Context, personBusinessProcess *domain.PersonBusinessProcess) error {
	p.logger.Info("PersonBusinessProcessDB -> Inserting document for: " + personBusinessProcess.PersonGlobalIdentifier)
	collection := p.MongoClient.Database("personBusinessProcessDB").Collection("personBusinessProcess")

	insertResult, err := collection.InsertOne(context.TODO(), personBusinessProcess)
	if err != nil {
		p.logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB -> Error: %s", err))
	}

	p.logger.Info(fmt.Sprintln("Inserted a single document: ", insertResult.InsertedID))
	return nil
}

//UpdatePersonBusinessProcessDocument inserts documents into the PersonBusinessProcess Collection
func (p *PersonBusinessProcessDB) UpdatePersonBusinessProcessDocument(ctx context.Context, key string, personBusinessProcess *domain.PersonBusinessProcess) error {
	p.logger.Info("PersonBusinessProcessDB -> Updating document for: " + key)
	filter := bson.M{"internalID": bson.D{{Key: "$eq", Value: key}}}

	collection := p.MongoClient.Database("personBusinessProcessDB").Collection("personBusinessProcess")
	updateResult, err := collection.ReplaceOne(context.TODO(), filter, personBusinessProcess)
	if err != nil {
		p.logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB -> Error: %s", err))
	}

	s := fmt.Sprintf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	p.logger.Info(s)
	return nil
}

//DeletePersonBusinessProcessDocument deletes documents into the PersonBusinessProcess Collection
func (p *PersonBusinessProcessDB) DeletePersonBusinessProcessDocument(ctx context.Context) error {
	p.logger.Info("PersonBusinessProcessDB -> Deleting document for:")
	collection := p.MongoClient.Database("personBusinessProcessDB").Collection("personBusinessProcess")
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		p.logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB -> Error: %s", err))
	}
	s := fmt.Sprintf("Deleted %v documents in the PersonBusinessProcess collection\n", deleteResult.DeletedCount)
	p.logger.Info(s)
	return nil
}

//DeleteAllPersonBusinessProcessDocument deletes documents into the PersonBusinessProcess Collection
func (p *PersonBusinessProcessDB) DeleteAllPersonBusinessProcessDocument(ctx context.Context) error {
	p.logger.Info("PersonBusinessProcessDB -> Deleting All documents")
	collection := p.MongoClient.Database("personBusinessProcessDB").Collection("personBusinessProcess")
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		p.logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB -> Error: %s", err))
	}
	s := fmt.Sprintf("Deleted %v documents in the PersonBusinessProcess collection\n", deleteResult.DeletedCount)
	p.logger.Info(s)
	return nil
}

//GetPersonBusinessProcessDocument gets documents into the PersonBusinessProcess Collection
func (p *PersonBusinessProcessDB) GetPersonBusinessProcessDocument(ctx context.Context, aBusinessPersonID string) (*domain.PersonBusinessProcess, error) {
	// Pass these options to the Find method
	p.logger.Info("PersonBusinessProcessDB -> Get document for: " + aBusinessPersonID)
	findOptions := options.Find()
	findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents
	var results []*domain.PersonBusinessProcess

	collection := p.MongoClient.Database("personBusinessProcessDB").Collection("personBusinessProcess")
	// Passing bson.D{{}} as the filter matches all documents in the collection
	filter := bson.M{"internalID": bson.D{{Key: "$eq", Value: aBusinessPersonID}}}

	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		p.logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB Get -> Filter Error: %s", err))
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem domain.PersonBusinessProcess
		err := cur.Decode(&elem)
		if err != nil {
			p.logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB Get -> Decode Error ->%s", err))
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		p.logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB Get -> Cur Error ->%s", err))
	}

	// Close the cursor once finished
	cur.Close(context.TODO())
	pbp := results[0]
	template := getBusinessProcessTemplate(pbp.BusinessProcessTemplateID)
	pbp.BusinessProcessTemplate = template
	we := make([]domain.EventExpectation, 0)
	for _, eeID := range pbp.WaitingExpectationsID {
		ees := template.EventExpectations
		for _, ee := range ees {
			if ee.ID == eeID {
				we = append(we, ee)
			}
		}
	}
	pbp.WaitingExpectations = we
	return pbp, nil
}

// Close closes the Mongo client
func (p *PersonBusinessProcessDB) Close(ctx context.Context) error {
	err := p.MongoClient.Disconnect(context.TODO())
	p.logger.Info("PersonBusinessProcessDB -> Connection to MongoDB closed.")
	return err
}
func getBusinessProcessTemplate(key string) domain.BusinessProcessTemplate {
	return domain.CreateBakeTemplate()
}
