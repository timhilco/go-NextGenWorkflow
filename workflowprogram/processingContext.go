package workflowprogram

import (
	"sync"

	"github.com/timhilco/go-NextGenWorkflow/util/logger"
)

// ProcessingContext is for providing runtime contect
const (
	C_ProcessingContext = "ProcessingContext"
)

type ProcessingContext struct {
	//Broker         string
	ProcessingDate string
	Logger         *logger.HilcoLogger
	//SyncWaitGroup  *sync.WaitGroup
	//TermChan       chan bool
	//PersonBusinessProcessMap map[string]domain.PersonBusinessProcess
	//BusinessProcessTemplate domain.BusinessProcessTemplate
}

//GroupRebalanceConfig is
type GroupRebalanceConfig struct {
	ProcessorType string
	ProducerTopic string
}

// KafkaPublisherConsumerProcessingContext is the processing context for kafka
type KafkaPublisherConsumerProcessingContext struct {
	Broker               string
	ConsumerGroupName    string
	ConsumerTopics       []string
	PublisherTopics      []string
	Logger               logger.HilcoLogger
	Callback             func(string)
	SyncWaitGroup        *sync.WaitGroup
	TermChan             chan bool
	SchedulerTask        SchedulerTask
	ProcessorType        string
	GroupRebalanceConfig GroupRebalanceConfig
}

// ProcessCommander start stop commands
type ProcessCommander interface {
	Start(context ProcessingContext)
	Stop(context ProcessingContext)
}
