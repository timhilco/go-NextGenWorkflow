package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/timhilco/go-NextGenWorkflow/databases"
	"github.com/timhilco/go-NextGenWorkflow/datastructures"
	"github.com/timhilco/go-NextGenWorkflow/domain"
	"github.com/timhilco/go-NextGenWorkflow/util"
	"github.com/timhilco/go-NextGenWorkflow/util/logger"
	"github.com/timhilco/go-NextGenWorkflow/workflowprogram"
)

func buildEventMessage(task string) []byte {
	_, event := util.BuildBakeMessage("Event", task, "id")
	return event
}
func buildBusinessProcessTemplate() domain.BusinessProcessTemplate {
	bpt := domain.CreateBakeTemplate()
	return bpt
}
func buildPersonBusinessProcess(businessProcessTemplate domain.BusinessProcessTemplate) domain.PersonBusinessProcess {
	we := businessProcessTemplate.InitialExpectations()
	personBusinessProcess := domain.PersonBusinessProcess{
		InternalID:                     "aBusinessProcessID",
		BusinessProcessReferenceNumber: "aBusinessProcessID",
		BusinessProcessTemplate:        businessProcessTemplate,
		PersonGlobalIdentifier:         "aPerson.InternalID",
		EffectiveDate:                  "anEffectiveDate",
		State:                          "Activated",
		WaitingExpectations:            we,
	}
	return personBusinessProcess
}

func setup(logger *logger.HilcoLogger) {
	context := databases.DatabaseContext{
		Logger: logger,
		URL:    "mongodb://localhost:27017",
	}
	database := databases.CreatePersonBusinessDB(context)
	database.DeleteAllPersonBusinessProcessDocument()
}

func TestBrokerServerStartUp(t *testing.T) {
	ctx := context.Background()
	configProperties := make(map[string]string)
	configProperties[workflowprogram.C_NumberOfSchedulers] = "2"
	now := time.Now()
	processingDate := now.Format(time.RFC3339)
	var logger = logger.NewMultiWithFile(false)
	setup(logger)
	var processingContext = workflowprogram.ProcessingContext{
		ProcessingDate: processingDate,
		Logger:         logger,
	}
	ctx2 := context.WithValue(ctx, workflowprogram.C_ProcessingContext, processingContext)
	ctx3, cancelFn := context.WithTimeout(ctx2, time.Duration(time.Second*10))
	defer cancelFn()
	sb, _ := workflowprogram.CreateSchedulerManager(ctx3, configProperties)
	sb.Start(ctx3)
	<-ctx3.Done()
	checkScehedulerBroker(t, sb)

}
func checkScehedulerBroker(t *testing.T, sb *workflowprogram.SchedulerBroker) {
	tasks := sb.SchedulerMap
	if len(tasks) != 2 {
		t.Fatalf("Number of task != 2 : %d", len(tasks))
	}
	for key, value := range tasks {
		t.Logf("Key: %s", key)
		t.Logf("BrokerID: %s -> TaskID: %s", value.Broker.BrokerID, value.SchedulerID)

	}

}
func setupSchedulerBroker(logger *logger.HilcoLogger) (*workflowprogram.SchedulerBroker, context.Context) {

	ctx := context.Background()
	configProperties := make(map[string]string)
	configProperties[workflowprogram.C_NumberOfSchedulers] = "2"
	now := time.Now()
	processingDate := now.Format(time.RFC3339)

	//setup(logger)
	var processingContext = workflowprogram.ProcessingContext{
		ProcessingDate: processingDate,
		Logger:         logger,
	}
	ctx2 := context.WithValue(ctx, workflowprogram.C_ProcessingContext, processingContext)
	ctx3, cancelFn := context.WithTimeout(ctx2, time.Duration(time.Second*30))
	defer cancelFn()
	sb, _ := workflowprogram.CreateSchedulerManager(ctx3, configProperties)
	sb.Start(ctx3)
	return sb, ctx3
}
func TestBakeAPersonCakeWithSchedulerBroker(t *testing.T) {
	var logger = logger.NewMultiWithFile(false)
	setup(logger)
	sb, ctx3 := setupSchedulerBroker(logger)

	person := domain.Person{
		ExternalID: uuid.New().String(),
		InternalID: "Tim",
		LastName:   "Sample",
		FirstName:  "Tim",
	}
	aBusinessProcessID := uuid.New().String()
	aBusinessProcessTemplate := "Baking"
	anEffectiveDate := "anEffectiveDate"
	personBusinessProcess, _ := activateBusinessProcessTemplateForPrincipal(ctx3,
		sb,
		person,
		aBusinessProcessID,
		aBusinessProcessTemplate,
		anEffectiveDate)
	fmt.Println(personBusinessProcess)
	internalId := personBusinessProcess.InternalID

	replayChannel := make(chan string)
	/*
		GetIngredients
		MixWetIngredients
		MixDryIngredients
		Bake
	*/

	logger.Info("===============================================================")
	logger.Info("===================GetIngredients===========+==================")
	event := buildTaskEvent("GetIngredients", internalId)
	processInboundMessage(ctx3, sb, replayChannel, []byte(event))
	logger.Info("###############################################################")

	logger.Info("===============================================================")
	logger.Info("======================MixWetIngredients=== ====================")
	event = buildTaskEvent("MixWetIngredients", internalId)
	processInboundMessage(ctx3, sb, replayChannel, []byte(event))
	logger.Info("###############################################################")

	logger.Info("===============================================================")
	logger.Info("==========================MixDryIngredients====================")
	event = buildTaskEvent("MixDryIngredients", internalId)
	processInboundMessage(ctx3, sb, replayChannel, []byte(event))
	logger.Info("###############################################################")
	logger.Info("===============================================================")
	logger.Info("==========================Bake ================================")
	event = buildTaskEvent("Bake", internalId)
	//processInboundMessage(ctx3, sb, replayChannel, []byte(event))
	logger.Info("###############################################################")

	<-ctx3.Done()
}
func processInboundMessage(ctx context.Context,
	messageProcessor workflowprogram.MessageProcessor,
	replayChannel chan string,
	jsonMessage []byte) (string, error) {
	return messageProcessor.ProcessInboundMessage(ctx, replayChannel, jsonMessage)
}
func activateBusinessProcessTemplateForPrincipal(ctx context.Context,
	messageProcessor workflowprogram.MessageProcessor,
	aPerson domain.Person,
	aBusinessProcessID string,
	aBusinessProcessTemplate interface{},
	anEffectiveDate string) (*domain.PersonBusinessProcess, error) {
	return messageProcessor.ActivateBusinessProcessTemplateForPrincipal(ctx, aPerson, aBusinessProcessID, aBusinessProcessTemplate, anEffectiveDate)

}
func buildTaskEvent(task string, aPersonBusinessProcessId string) string {
	_, jsonMessage := util.BuildBakeMessage("Event", task, aPersonBusinessProcessId)
	event := string(jsonMessage)

	message := datastructures.SchedulerInboundMessage{
		SchedulerCode:              "aSchedulerId",
		SchedulerMessageType:       "Event",
		SchedulerMessageKey:        "0001",
		SchedulerMessageJSONString: event,
	}

	jsonMessage2, err := json.Marshal(&message)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	command2 := string(jsonMessage2)
	return command2
}

/*
func buildStartCommand() string {
	schedulerId, jsonMessage := util.BuildBakeMessage("Command", "Start", "aUUID")
	command := string(jsonMessage)

	message := datastructures.SchedulerInboundMessage{
		SchedulerCode:              schedulerId,
		SchedulerMessageType:       "Command",
		SchedulerMessageKey:        "0001",
		SchedulerMessageJSONString: command,
	}

	jsonMessage2, err := json.Marshal(&message)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	command2 := string(jsonMessage2)
	return command2
}
*/
func TestPersonAdapter(t *testing.T) {
	fmt.Println("Starting test for processing an event")
	bpt := buildBusinessProcessTemplate()
	pbp := buildPersonBusinessProcess(bpt)
	event := buildEventMessage("GetIngredients")
	var m datastructures.Message
	err := json.Unmarshal(event, &m)
	if err != nil {
		fmt.Println("Error")
	}
	//Test 1
	pbp.Execute(m)
	var pass bool = true
	we := pbp.WaitingExpectations
	for _, ee := range we {
		if (ee.ID == "20") || (ee.ID == "30") {

		} else {
			pass = false
		}
	}
	fmt.Printf("Pass=%t", pass)
}
func TestTimeoutBatchJob(t *testing.T) {
	var logger = logger.NewMultiWithFile(false)
	sb, ctx3 := setupSchedulerBroker(logger)
	sb.ProcessTimeouts(ctx3)

}
/*
func TestAntlrPencilCalculator(t *testing.T) {
	calc.Execute()
}
*/
