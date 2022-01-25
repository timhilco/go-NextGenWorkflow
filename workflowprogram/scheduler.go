package workflowprogram

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/timhilco/go-NextGenWorkflow/databases"
	"github.com/timhilco/go-NextGenWorkflow/datastructures"
	"github.com/timhilco/go-NextGenWorkflow/domain"
	"github.com/timhilco/go-NextGenWorkflow/util/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//Config Values
	C_NumberOfSchedulers = "NumberOfSchedulers"
	//Defaults
	DefaultNumberOfSchedulers = 2
)

type MessageProcessor interface {
	ProcessInboundMessage(ctx context.Context, replayChannel chan string, jsonMessage []byte) (string, error)
	ActivateBusinessProcessTemplateForPrincipal(ctx context.Context, aPerson domain.Person, aBusinessProcessID string, aBusinessProcessTemplate interface{},
		anEffectiveDate string) (*domain.PersonBusinessProcess, error)
}
type SchedulerBroker struct {
	BrokerID                 string
	SchedulerMap             map[string]*SchedulerTask
	NumberOfSchedulers       int
	Status                   string
	ServerContext            context.Context
	ProcessingContext        ProcessingContext
	PersonBusinessProcessMap map[string]*domain.PersonBusinessProcess
	mongoDatabaseConnection  databases.PersonBusinessProcessDB
}

func CreateSchedulerManager(ctx context.Context, properties map[string]string) (*SchedulerBroker, error) {

	numberOfSchedulers := DefaultNumberOfSchedulers
	if properties[C_NumberOfSchedulers] != "" {
		numberOfSchedulers, _ = strconv.Atoi(properties[C_NumberOfSchedulers])
	}
	sm := &SchedulerBroker{
		BrokerID:           uuid.New().String(),
		NumberOfSchedulers: numberOfSchedulers,
		Status:             "Created",
	}
	sm.PersonBusinessProcessMap = make(map[string]*domain.PersonBusinessProcess)
	sm.SchedulerMap = make(map[string]*SchedulerTask)
	processingContext := ctx.Value(C_ProcessingContext).(ProcessingContext)
	logger := processingContext.Logger
	databaseContext := databases.DatabaseContext{
		URL:    "mongodb://localhost:27017",
		Logger: logger,
	}
	connection := databases.CreatePersonBusinessDB(ctx, databaseContext)
	sm.mongoDatabaseConnection = connection

	return sm, nil
}
func (sm *SchedulerBroker) Start(ctx context.Context) {

	for i := 0; i < sm.NumberOfSchedulers; i++ {
		SchedulerTask := &SchedulerTask{}
		SchedulerTask.Initialize(ctx, sm)
		key := uuid.New().String()
		SchedulerTask.SchedulerID = key
		sm.SchedulerMap[key] = SchedulerTask
		go SchedulerTask.Start()
	}
	sm.Status = "Started"

}
func (sm *SchedulerBroker) assignScheduler() *SchedulerTask {
	var s *SchedulerTask
	for _, value := range sm.SchedulerMap {
		s = value
	}
	return s
}
func (sm *SchedulerBroker) Register(processingContext ProcessingContext) (*SchedulerTask, error) {

	SchedulerTask := sm.assignScheduler()
	return SchedulerTask, nil
}
func (s *SchedulerBroker) ProcessInboundMessage(ctx context.Context, replayChannel chan string, jsonMessage []byte) (string, error) {
	schedulerTask := s.assignScheduler()
	return schedulerTask.ProcessInboundMessage(ctx, replayChannel, jsonMessage)

}
func (s *SchedulerBroker) ActivateBusinessProcessTemplateForPrincipal(ctx context.Context,

	principal domain.Person,
	aBusinessProcessID string,
	aBusinessProcessTemplate interface{},
	anEffectiveDate string) (*domain.PersonBusinessProcess, error) {
	schedulerTask := s.assignScheduler()
	return schedulerTask.ActivateBusinessProcessTemplateForPrincipal(ctx, principal, aBusinessProcessID, aBusinessProcessTemplate, anEffectiveDate)
}

// ProcessTimeouts is a function to process timeouts
func (s *SchedulerBroker) ProcessTimeouts(ctx context.Context) error {
	processingContext := ctx.Value(C_ProcessingContext).(ProcessingContext)
	logger := processingContext.Logger
	logger.Info("Processing timeoutsSchedulerTask --> Processing Timeouts: ")
	dbcontext := databases.DatabaseContext{
		Logger: logger,
		URL:    "mongodb://localhost:27017",
	}
	db := databases.CreatePersonBusinessDB(ctx, dbcontext)

	findOptions := options.Find()
	findOptions.SetLimit(2)

	// Here's an array in which you can store the decoded documents
	//var results []*domain.PersonBusinessProcess

	collection := db.MongoClient.Database("personBusinessProcessDB").Collection("personBusinessProcess")
	// Passing bson.D{{}} as the filter matches all documents in the collection
	filter := bson.M{"state": bson.D{{Key: "$eq", Value: domain.C_Status_InProcess}}}

	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB -> Error: %s", err))
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var pbp *domain.PersonBusinessProcess
		err := cur.Decode(&pbp)
		if err != nil {
			logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB: Error ->%s", err))
		}
		doUpdate := false
		if pbp.IsDone() {
			if pbp.State == domain.C_Status_Complete {

			} else {
				pbp.State = domain.C_Status_Complete
				doUpdate = true
			}
		} else {
			schedulerEvent := domain.SchedulerEvent{
				EventID:       "SchedulerTask Timeout",
				EffectiveDate: "aDate",
			}
			pbp.ProcessTimeouts(schedulerEvent)
			doUpdate = true

		}
		if doUpdate {
			s.mongoDatabaseConnection.UpdatePersonBusinessProcessDocument(ctx, pbp.InternalID, pbp)
		}
	}

	if err := cur.Err(); err != nil {
		logger.Fatal(fmt.Sprintf("PersonBusinessProcessDB: Error ->%s", err))
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	return nil
}

/*
	const processTimeouts = function () {
		logger.info(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>");
		logger.info("SchedulerTask -> Processing Timeouts for " + date);
		logger.log('info', ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>");


		const schedulerEvent = module.exports.event({
			id: "SchedulerTask Timeout",
			effectiveDate: date
		});
		businessProcessDB.values().forEach(function (personBusinessProcess) {
			personBusinessProcess.processTimeouts(schedulerEvent);
		});
	};
	that.processTimeouts = processTimeouts;
	// End of Private methods
*/
// SchedulerTask object
type SchedulerTask struct {
	//ProcessingContext     ProcessingContext
	ProcessingDate        string
	Logger                logger.HilcoLogger
	ConsumersReplyChannel map[string]chan string
	RequestChannel        chan string
	//mongoDatabaseConnection  databases.PersonBusinessProcessDB

	SchedulerID string
	//BusinessProcessTemplate domain.BusinessProcessTemplate
	Broker *SchedulerBroker
}

//Register the consumer and reply with
func (s *SchedulerTask) Register(replyChannel chan string) (string, chan string) {
	key := uuid.New().String()
	s.ConsumersReplyChannel[key] = replyChannel
	return key, s.RequestChannel
}

//Initialize is
func (s *SchedulerTask) Initialize(ctx context.Context, broker *SchedulerBroker) {
	processingContext, _ := ctx.Value(C_ProcessingContext).(ProcessingContext)
	s.Logger = *processingContext.Logger
	text := fmt.Sprintf("SchedulerTask:%s ->Initializing the SchedulerTask", s.SchedulerID)
	s.Logger.Info(text)
	s.ConsumersReplyChannel = make(map[string]chan string)
	s.ProcessingDate = processingContext.ProcessingDate
	s.RequestChannel = make(chan string)
	s.Broker = broker

	//s.BusinessProcessTemplate = "Baking"
	/*
		databaseContext := databases.DatabaseContext{
			URL:    "mongodb://localhost:27017",
			Logger: context.Logger,
		}
		connection := databases.CreatePersonBusinessDB(databaseContext)
		s.mongoDatabaseConnection = connection
	*/

}

//Start is
func (s *SchedulerTask) Start() {
	text := fmt.Sprintf("SchedulerTask:%s ->Starting the Event Loop", s.SchedulerID)
	s.Logger.Info(text)
	doTerm := false
	ticker := time.NewTicker(5000 * time.Millisecond)
	for !doTerm {
		select {

		case <-ticker.C:
			t := time.Now()
			fTime := t.Format(time.RFC3339)
			text := fmt.Sprintf("SchedulerTask:%s -> Hit Ticker @ %s", s.SchedulerID, fTime)
			s.Logger.Info(text)
			/*
				case <-s.Broker.ProcessingContext.TermChan:
					doTerm = true
			*/
		case message := <-s.RequestChannel:
			jsonMessage := []byte(message)
			f := make(map[string]string)
			json.Unmarshal(jsonMessage, &f)
			code := f["schedulerCode"]
			command := f["schedulerMessageJSONString"]
			replyChannel := s.ConsumersReplyChannel[code]
			internalReferenceNumber, _ := s.ProcessInboundMessage(s.Broker.ServerContext, replyChannel, []byte(command))
			fmt.Println("number: " + internalReferenceNumber)
			//replyChannel <- internalReferenceNumber

		default:
			/*
				switch e := ev.(type) {
				case *kafka.Message:
					// Process ingress car event message
					kafkaLogger.Info(fmt.Sprintf("Consumer Processing Inbound Message: %sn", ev))
					processInboundMessage(e)
				case kafka.Error:
					// Errors are generally just informational.
					kafkaLogger.Info(fmt.Sprintf("Consumer error: %sn", ev))
				default:
					kafkaLogger.Info(fmt.Sprintf("Consumer event: %s: ignored", ev))
				}
			*/
		}
	}
}
func (s *SchedulerTask) ProcessInboundMessage(ctx context.Context, replayChannel chan string, jsonMessage []byte) (string, error) {
	//s.Logger.Info("SchedulerTask -> Message  ----->")
	//s.Logger.Info(message)
	//s.Logger.Info("SchedulerTask -> <--------------")

	var e datastructures.SchedulerInboundMessage
	err := json.Unmarshal(jsonMessage, &e)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("SchedulerTask:%s --> processInboundMessage-> Bad Message: %s", s.SchedulerID, err))
	}
	jsonMessage2 := []byte(e.SchedulerMessageJSONString)
	var m datastructures.Message
	err = json.Unmarshal(jsonMessage2, &m)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("SchedulerTask:%s --> processInboundMessage-> Bad Message: %s", s.SchedulerID, err))
	}
	messageType := m.Message.Header.MessageType
	switch messageType {
	case "Event":
		s.Logger.Info(fmt.Sprintf("SchedulerTask:%s --> processInboundMessage -> Processing Event", s.SchedulerID))
		s.ProcessEvent(ctx, m)
	case "Command":
		s.Logger.Info(fmt.Sprintf("SchedulerTask:%s --> processInboundMessage -> Processing Command", s.SchedulerID))

		person := domain.Person{
			ExternalID: m.Message.Header.GlobalPersonIdentfier,
			InternalID: "Tim",
			LastName:   "Sample",
			FirstName:  "Tim",
		}
		aBusinessProcessID := m.Message.Body.ClientReferenceNumber
		aBusinessProcessTemplate := "Baking"
		anEffectiveDate := m.Message.Body.EffectiveDate
		personBusinessProcess, err := s.ActivateBusinessProcessTemplateForPrincipal(
			ctx,
			person,
			aBusinessProcessID,
			aBusinessProcessTemplate,
			anEffectiveDate)
		return personBusinessProcess.InternalID, err

	default:

	}
	return "", nil
}

/*

	// Define private methods
	function doneWorkflow() {
		// ^businessProcessDB select: [:each | each isDone]!
		logger.debug("SchedulerTask --> SchedulerTask:doneWorkflow");
		let arr = businessProcessDB.filter(function (entry) {
			// printObject(entry);
			let bool = entry.isDone();
			return bool;
		});
		return arr;
	}
*/

/*
	function personBusinessProcessForEventOwner(anOwner) {
		logger.debug('SchedulerTask-- > Getting personBusinessProcess for : ' + anOwner);
		return businessProcessDB.get(anOwner);
	}
*/

// UnscheduleDoneWorkflows cleans up completed workflows
func (s *SchedulerTask) UnscheduleDoneWorkflows() (string error) {
	fmt.Println("SchedulerTask --> UnscheduleDoneWorkflows: ")
	return nil
}

/*
	function unscheduleDoneWorkflows() {
		// self doneWorkflows do: [:each | self unschedule: each]! !
		logger.debug('SchedulerTask-- > unscheduleDoneWorkflows');
		let aArray = doneWorkflow();
		aArray.forEach(function (anEntry) {
			logger.log('debug', anEntry);
			let key = anEntry.id;
			unscheduleAnWorkflow(key);
		});
	}
*/

/*
	 const scheduleAWorkflowForObject = async (anWorkflow, anObject) =>{
		let text = "SchedulerTask --> scheduleAWorkflowForObject enter - Adding Workflow:\n--- " + anWorkflow.toString() + " \n---for: " + anObject;
		logger.debug(text);
		businessProcessDB.set(anObject, anWorkflow);
		const mongoObject = anWorkflow.toMongoDocument();
		await businessProcessDBMongo.insert(anObject, mongoObject)
		logger.debug(businessProcessDB.entries());
		logger.debug("SchedulerTask --> scheduleAWorkflowForObject exit ------------------------------------------------------------------");

	}
*/

// ActivateBusinessProcessTemplateForPrincipal create a new business process
func (s *SchedulerTask) ActivateBusinessProcessTemplateForPrincipal(ctx context.Context,
	aPerson domain.Person,
	aBusinessProcessID string,
	aBusinessProcessTemplate interface{},
	anEffectiveDate string) (*domain.PersonBusinessProcess, error) {

	s.Logger.Info(fmt.Sprintf("SchedulerTask:%s --> activateBusinessProcessTemplateForPrincipal: %s", s.SchedulerID, aPerson.LastName))
	businessProcessTemplate := domain.CreateBakeTemplate()
	we := businessProcessTemplate.InitialExpectations()
	personBusinessProcess := domain.PersonBusinessProcess{
		InternalID:                     aBusinessProcessID,
		BusinessProcessReferenceNumber: aBusinessProcessID,
		BusinessProcessTemplate:        businessProcessTemplate,
		PersonGlobalIdentifier:         aPerson.InternalID,
		EffectiveDate:                  anEffectiveDate,
		State:                          "Activated",
		WaitingExpectations:            we,
	}
	err := s.Broker.mongoDatabaseConnection.InsertPersonBusinessProcessDocument(ctx, &personBusinessProcess)
	s.Broker.PersonBusinessProcessMap[aBusinessProcessID] = &personBusinessProcess

	if err != nil {
		s.Logger.Fatal("SchedulerTask: ActivateBusinessProcessTemplateForPrincipal Error")
	}

	return &personBusinessProcess, nil
}

/*
	const activateBusinessProcessTemplateForPrincipal = function (aPerson, anBusinessProcessId, aBusinessProcessTemplate, anEffectiveDate) {
		logger.debug('SchedulerTask --> activateBusinessProcessTemplateForObject:\n ' +
			aBusinessProcessTemplate.toString() + " \n---Person: " + aPerson);
		const spec = {
			"businessProcessTemplate": aBusinessProcessTemplate,
			"person": aPerson,
			"effectiveDate": anEffectiveDate,
			"businessProcessId": anBusinessProcessId
		};
		const aWorkflow = personBusinessProcessClass.personBusinessProcess(spec);
		logger.debug(aWorkflow.toString());
		scheduleAWorkflowForObject(aWorkflow, anBusinessProcessId);
		logger.debug("SchedulerTask --> activateBusinessProcessTemplateForPrincipal exit ------------------------------------------------------------------");
		return aWorkflow;
	};
	that.activateBusinessProcessTemplateForPrincipal = activateBusinessProcessTemplateForPrincipal;
	//
*/

// ProcessEvent is a method to process an event acgainst a business process
func (s *SchedulerTask) ProcessEvent(ctx context.Context, jsonEvent datastructures.Message) (string, error) {
	eventName := jsonEvent.Message.Header.EventName
	s.Logger.Info(fmt.Sprintf("SchedulerTask:%s Process Event --> Processing Event: %s", s.SchedulerID, eventName))

	/*
		jsonMessage := []byte(jsonEvent)
		var e datastructures.Message
		err := json.Unmarshal(jsonMessage, &e)
		if err != nil {
			s.Logger.Info(fmt.Sprintf("SchedulerTask -> Bad Message: %s", err))
		}
	*/
	aBusinessProcessID := jsonEvent.Message.Body.ClientReferenceNumber
	/*
		aBusinessProcessTemplate := domain.CreateBakeTemplate()
		anEffectiveDate := jsonEvent.Message.Body.EffectiveDate
		personBusinessProcess := domain.PersonBusinessProcess{
			InternalID:                     aBusinessProcessID,
			BusinessProcessReferenceNumber: aBusinessProcessID,
			BusinessProcessTemplate:        aBusinessProcessTemplate,
			PersonGlobalIdentifier:         "aPerson.InternalID",
			EffectiveDate:                  anEffectiveDate,
			State:                          "Activated",
		}*/
	var personBusinessProcess *domain.PersonBusinessProcess
	//personBusinessProcess = s.Broker.PersonBusinessProcessMap[aBusinessProcessID]
	personBusinessProcess, _ = s.Broker.mongoDatabaseConnection.GetPersonBusinessProcessDocument(ctx, aBusinessProcessID)
	personBusinessProcess.Execute(jsonEvent)
	s.Broker.mongoDatabaseConnection.UpdatePersonBusinessProcessDocument(ctx, aBusinessProcessID, personBusinessProcess)
	return "AnEvent", nil
}

/*
	const event = function (anEvent) {
		logger.debug('SchedulerTask --> event:\n' + anEvent.toString());
		const owner = anEvent.businessProcessId;
		const personBusinessProcess = personBusinessProcessForEventOwner(owner);
		personBusinessProcess.execute(anEvent);
		logger.debug("Event Execute completed - Workflow now:");
		logger.debug(personBusinessProcess.toString());
		unscheduleDoneWorkflows();
		logger.debug("SchedulerTask --> event exit ------------------------------------------------------------------");
	};
	that.event = event;
*/

/*
	const setProcessingDate = function (aDate) {
		date = aDate;
	};
	that.setProcessingDate = setProcessingDate;
	//
*/
