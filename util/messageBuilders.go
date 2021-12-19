package util

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/timhilco/go-NextGenWorkflow/datastructures"
)

//BuildBakeMessage is for building event and command messages
func BuildBakeMessage(messageType string, task string, clientReferenceNumber string) (string, []byte) {

	tag := "Cake"
	action := task
	eventName := tag + ":" + action
	aUUID := uuid.New().String()
	timeNow := time.Now()
	timeStamp := timeNow.Format(time.RFC3339)
	type EventBody struct {
		BusinessProcessID string `json:"businessProcessId"`
		Task              string `json:"task"`
	}

	var pisr = datastructures.PersonIdentificationSystemOfRecord{
		SystemOfRecordSystemID:              "SR_SI",
		SystemOfRecordApplicationName:       "SR_AN",
		SystemOfRecordApplicationInstanceID: "SR_AII",
		SystemOfRecordDatabaseSchema:        "SR-DS",
		PlatformInternalID:                  "SR-PII",
		PlatformExternalID:                  "SR-PEI",
		PlatformRoleType:                    "SR-PRT",
		PlatformClientID:                    "SR-PCI",
	}

	var rri = datastructures.RelatedResources{
		RelatedResourcesItem:      "RR-I",
		RelatedResourceIdentifier: "RR-ID",
		RelatedResourceState:      "RR-S",
		RelatedResourceDescrption: "RR-D",
	}
	var rr []datastructures.RelatedResources
	rr = append(rr, rri)

	var phi = datastructures.PublisherHistoryItem{
		PublisherID:                    "PHI-ID",
		PublisherApplicationName:       "PHI-AN",
		PublisherApplicationInstanceID: "PHI-PAII",
		MessageID:                      "PHI-MI",
		MessageTopic:                   "PHI-MT",
		MessageSubTopic:                "PHI-MST",
		EventName:                      "PHI-EV",
		MessageTimestamp:               "PHI-MT",
		SequenceNumber:                 1,
	}
	var ph []datastructures.PublisherHistoryItem
	ph = append(ph, phi)

	var h = datastructures.Header{
		MessageID:        aUUID,
		MessageType:      messageType,
		MessageNamespace: "com.alight.messages/events/person/goalEvent",
		MessageVersion:   "1.0",
		MessageTopic:     "Baking",
		BusinessDomain:   "Person",
		EventName:        eventName,
		MessageTimestamp: timeStamp,
		Action:           action,
		Tag:              tag,

		PersonName:               "Tim",
		TagObjectID:              "aTagObjectId",
		CorrelationID:            aUUID,
		EventBodyNamespace:       "com.alight.messages/events/operations/platformProcessingEvent",
		CorrelationIDType:        "Session",
		NormalizedClientID:       "P0095",
		AgentID:                  "N/A",
		GlobalPersonIdentfier:    "N/A",
		PublisherID:              "alight",
		PublisherApplicationName: "meQ",

		PublisherApplicationInstanceID:     "1",
		PublishingPlatformsHistory:         &ph,
		PersonIdentificationSystemOfRecord: &pisr,
		RelatedResources:                   &rr,

		IsSyntheticEvent: false,
	}
	if clientReferenceNumber == "" {
		clientReferenceNumber = aUUID
	}
	var b = datastructures.CommandBody{
		ClientReferenceNumber: clientReferenceNumber,
		Task:                  task,
		EffectiveDate:         "anEffectiveDate",
	}
	var m = datastructures.EventMessage{
		Header: &h,
		Body:   &b,
	}
	var e = datastructures.Message{
		Message: &m,
	}
	/*
		fmt.Println(h)
		fmt.Println(b)
		fmt.Println(m)
		fmt.Println(e)
	*/
	jsonEvent, err := json.Marshal(&e)
	if err != nil {
		fmt.Println(err)
		return "error", nil
	}

	return aUUID, jsonEvent

}
