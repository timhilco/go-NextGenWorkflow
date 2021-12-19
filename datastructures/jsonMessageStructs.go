package datastructures

// PublisherHistoryItem is a publisher history array element
type PublisherHistoryItem struct {
	PublisherID                    string `json:"publisherId"`
	PublisherApplicationName       string `json:"publisherApplicationName"`
	PublisherApplicationInstanceID string `json:"PublisherApplicationInstanceID"`
	MessageID                      string `json:"messageId"`
	MessageTopic                   string `json:"messageTopic"`
	MessageSubTopic                string `json:"messageSubTopic"`
	EventName                      string `json:"eventName"`
	MessageTimestamp               string `json:"messageTimestamp"`
	SequenceNumber                 int    `json:"sequenceNUmber"`
}

// PersonIdentificationSystemOfRecord is a
type PersonIdentificationSystemOfRecord struct {
	SystemOfRecordSystemID              string `json:"systemOfRecordSystemId"`
	SystemOfRecordApplicationName       string `json:"systemOfRecordApplicationName"`
	SystemOfRecordApplicationInstanceID string `json:"systemOfRecordApplicationInstanceId"`
	SystemOfRecordDatabaseSchema        string `json:"systemOfRecordDatabaseSchema"`
	PlatformInternalID                  string `json:"platformInternalId"`
	PlatformExternalID                  string `json:"platformExternalId"`
	PlatformRoleType                    string `json:"platformRoleType"`
	PlatformClientID                    string `json:"latformClientId"`
}

// RelatedResources is a
type RelatedResources struct {
	RelatedResourcesItem      string `json:"relatedResourceItem"`
	RelatedResourceIdentifier string `json:"relatedResourceIdentifier"`
	RelatedResourceState      string `json:"relatedResourceState"`
	RelatedResourceDescrption string `json:"relatedResourceDescrption"`
}

//CommandBody id
type CommandBody struct {
	ClientReferenceNumber string `json:"clientReferenceNumber"`
	Task                  string `json:"task"`
	EffectiveDate         string `json:"effectiveDate"`
}

// EventMessage is the overall message structure
type EventMessage struct {
	Header *Header      `json:"header"` //PgEventHeader *PgEventHeader `json:"pgEventHeader"`
	Body   *CommandBody `json:"body"`
}

// CommandMessage is the overall message structure
type CommandMessage struct {
	Header *Header      `json:"header"` //PgEventHeader *PgEventHeader `json:"pgEventHeader"`
	Body   *CommandBody `json:"body"`
}

// Message holds the overall message structure
type Message struct {
	Message *EventMessage `json:"message"`
}

// Header is the message header structure
type Header struct {
	MessageID          string `json:"messageId"`
	MessageType        string `json:"messageType"`
	MessageNamespace   string `json:"messageNamespace"`
	MessageVersion     string `json:"messageVersion"`
	MessageTopic       string `json:"messageTopic"`
	MessageSubTopic    string `json:"messageSubTopic"`
	EventName          string `json:"eventName"`
	EventBodyNamespace string `json:"eventBodyNamespace"`
	Tag                string `json:"tag"`
	Action             string `json:"action"`
	BusinessDomain     string `json:"businessDomain"`
	MessageTimestamp   string `json:"messageTimestamp"`

	TagObjectID                        string                              `json:"tagObjectId"`
	CorrelationID                      string                              `json:"correlationId"`
	CorrelationIDType                  string                              `json:"correlationIdType"`
	NormalizedClientID                 string                              `json:"normalizedClientId"`
	AgentID                            string                              `json:"agentId"`
	GlobalPersonIdentfier              string                              `json:"globalPersonIdentifier"`
	PublisherID                        string                              `json:"publisherId"`
	PublisherApplicationName           string                              `json:"publisherApplicationName"`
	PublisherApplicationInstanceID     string                              `json:"publisherApplicationInstanceID"`
	PublishingPlatformsHistory         *[]PublisherHistoryItem             `json:"publishingPlatformsHistory"`
	PersonIdentificationSystemOfRecord *PersonIdentificationSystemOfRecord `json:"personIdentificationSystemOfRecord"`
	RelatedResources                   *[]RelatedResources                 `json:"relatedResources"`
	IsSyntheticEvent                   bool                                `json:"isSyntheticEvent"`
	PersonName                         string                              `json:"personName"`
}

// PgEventHeader is a event type header
type PgEventHeader struct {
	Value string `json:"value"`
}

// OppEventHeader is the header for the runtime events
type OppEventHeader struct {
	PlatformProcessStartTimestamp string `json:"platformProcessStartTimestamp"`
	PlatformProcessEndTimestamp   string `json:"platformProcessEndTimestamp"`
	PlatformProcessElapsedTime    int    `json:"platformProcessElapsedTime"`
}

//SchedulerInboundMessage is for sending messages to the SchedulerTask
type SchedulerInboundMessage struct {
	SchedulerCode              string `json:"schedulerCode"`
	SchedulerMessageType       string `json:"schedulerMessageType"`
	SchedulerMessageKey        string `json:"schedulerMessageKey"`
	SchedulerMessageJSONString string `json:"schedulerMessageJSONString"`
}
