package domain

import (
	"fmt"

	"github.com/timhilco/go-NextGenWorkflow/datastructures"
	loggerInterface "github.com/timhilco/go-NextGenWorkflow/util/logger"
)

// BusinessProcesss is an type of hilcoEvent
type BusinessProcesss struct {
	EventID string
}

var logger = loggerInterface.NewMultiWithFile(false)

const (
	C_Status_Activated = "Activated"
	C_Status_InProcess = "InProcess"
	C_Status_Complete  = "Complete"
)

type ObjectID string

// PersonBusinessProcess is an type of hilcoEvent
type PersonBusinessProcess struct {
	InternalID                     string                  `bson:"internalID,omitempty"`
	BusinessProcessReferenceNumber string                  `bson:"businessReferenceNumber,omitempty"`
	PersonGlobalIdentifier         string                  `bson:"personGlobalIdentifier,omitempty"`
	EffectiveDate                  string                  `bson:"effectveDate,omitempty"` //Date
	State                          string                  `bson:"state,omitempty"`
	BusinessProcessTemplate        BusinessProcessTemplate `bson:"-"`
	BusinessProcessTemplateID      string                  `bson:"businessProcessTemplate,omitempty"`
	WaitingExpectations            []EventExpectation      `bson:"-"`
	WaitingExpectationsID          []string                `bson:"waitingExpectationsID,omitempty"`
}

func (p *PersonBusinessProcess) Break() {

}
func (p *PersonBusinessProcess) EventOutOfSequence(anEvent EventDefinition) {

}
func (p *PersonBusinessProcess) SetBusinessProcessTemplate(aTemplate BusinessProcessTemplate) {

}

//Execute an Event against a Business Process
func (p *PersonBusinessProcess) Execute(jsonMessage datastructures.Message) {

	id := jsonMessage.Message.Header.MessageID
	eventName := jsonMessage.Message.Header.EventName
	s := fmt.Sprintf("Person Business Process -> Execute: Id= %s , EventName= %s ", id, eventName)
	logger.Info(s)
	expectation := p.nextExpectationAtEventIfAbsent(jsonMessage)
	expectation.ExecuteFor(jsonMessage, p)

}

/*
		const execute = function (anEvent) {
		logger.log('debug', 'Executing Event in Workflow :' + anEvent.toString());
		const expectation = nextExpectationAtEventIfAbsent(anEvent);
		expectation.executeFor(anEvent, that);
	};
	that.execute = execute;
*/

//NextExpectationAtEventIfAbsent is
func (p *PersonBusinessProcess) nextExpectationAtEventIfAbsent(jsonMessage datastructures.Message) EventExpectation {
	id := jsonMessage.Message.Header.MessageID
	eventName := jsonMessage.Message.Header.EventName
	s := fmt.Sprintf("Person Business Process -> nextExpectationAtEventIfAbsent: Id= %s , EventName= %s ", id, eventName)
	logger.Info(s)
	waitingExpections := p.WaitingExpectations
	var expectations []EventExpectation
	for _, expectation := range waitingExpections {
		if expectation.Expects(jsonMessage) {
			expectations = append(expectations, expectation)
		}
	}
	if len(expectations) == 0 {
		logger.Info("Person Business Process nextExpectationAtEventIfAbsent -> No expectations found")

	} else {
		// logger.log('debug',"Expectations:"+ expectations.toString());
		expectation := expectations[0]
		// logger.log('debug',"Found expectation:"+ expectation);
		s := fmt.Sprintf("Person Business Process nextExpectationAtEventIfAbsent -> Expectations found: EE:%s", expectation.ID)
		logger.Info(s)
		return expectation
	}
	return EventExpectation{}
}

/*
	function nextExpectationAtEventIfAbsent(anEvent /* , bloc) {
		const expectations = waitingExpectations.filter(
			function (expectation) {
				return expectation.expects(anEvent);
			});
		// logger.log('debug',"Expectations length = "+ expectations.length );
		if (expectations.length === 0) {
			logger.log('debug', "No expectations found");
		} else {
			// logger.log('debug',"Expectations:"+ expectations.toString());
			const expectation = expectations.one();
			// logger.log('debug',"Found expectation:"+ expectation);
			return expectation;
		}

	}
*/

// ToString is
func (p *PersonBusinessProcess) String() string {
	id := p.BusinessProcessReferenceNumber
	person := p.PersonGlobalIdentifier
	state := p.State
	text := id + " \n"
	text = text + person + " \n"
	text = text + "state= " + state + " \n"
	waitingExpectations := p.WaitingExpectations
	for _, we := range waitingExpectations {

		text = text + we.String() + " \n"
	}
	text = text + "--------------"
	s := fmt.Sprintf("Person Business Process -> \n %s", text)
	return s

}

/*
	const toString = function () {
		let text = '';
		text = text + "-----------------Workflow ----------------------\n";
		text = text + "id: " + that.id + "\n";
		text = text + that.effectiveDate.toString() + "\n";
		text = text + businessProcessTemplate.toString() + "\n";
		text = text + "-----------------Waiting Expectations -------------\n";
		waitingExpectations.forEach(function (eventExpectation) {
			text = text + ">>" + eventExpectation.toString() + ';\n';
		});
		text = text + "---------------------------------------------------";
		return text;
	};
	that.toString = toString;
*/

//DontExpectAnEventExpectation is
func (p *PersonBusinessProcess) DontExpect(anEventExpectation *EventExpectation) {

	s := fmt.Sprintf("Person Business Process -> DontExpectAnEventExpectation: EE:%s", anEventExpectation.ID)
	logger.Info(s)
	var newWaitingExpectations []EventExpectation
	for _, expectation := range p.WaitingExpectations {
		if expectation.ID != anEventExpectation.ID {
			newWaitingExpectations = append(newWaitingExpectations, expectation)

		}
	}
	p.WaitingExpectations = newWaitingExpectations
}

/*
	const dontExpectAnEventExpectation = function (anEventExpectation) {
		waitingExpectations.delete(anEventExpectation);

	};
	that.dontExpectAnEventExpectation = dontExpectAnEventExpectation;
*/

//ExpectAll is
func (p *PersonBusinessProcess) ExpectAll(eventExpectations []EventExpectation) {
	var e string
	for _, ee := range eventExpectations {
		e = e + ee.ID + " "
	}
	e = "[ " + e + " ]"
	s := fmt.Sprintf("Person Business Process -> ExpectAll: EE: %s", e)
	logger.Info(s)
	for _, ee := range eventExpectations {
		p.Expect(ee)
	}

}

/*
	const expectAll = function (aCollection) {
		aCollection.forEach(function (each) {
			that.expect(each);
		});

	};
	that.expectAll = expectAll;
*/

//Expect is
func (p *PersonBusinessProcess) Expect(eventExpectation EventExpectation) {
	s := fmt.Sprintf("Person Business Process -> Expect: EE:%s", eventExpectation.ID)
	logger.Info(s)
	we := p.WaitingExpectations
	we = append(we, eventExpectation)
	p.WaitingExpectations = we
	eeIds := make([]string, 0)
	for _, ee := range we {
		eeIds = append(eeIds, ee.ID)
	}
	p.WaitingExpectationsID = eeIds
	p.State = C_Status_InProcess
	s = fmt.Sprintf("Person Business Process -> Expect - Updated: %s", p.InternalID)
	logger.Info(s)
}

/*
	const expect = function (anEventExpectation) {
		waitingExpectations.add(anEventExpectation);

	};
	that.expect = expect;

	anWorkflow.expectationsFollowing( that);
*/

//ExpectationsFollowing is
func (p *PersonBusinessProcess) ExpectationsFollowing(eventExpectation *EventExpectation) []EventExpectation {
	s := fmt.Sprintf("Person Business Process -> ExpectationsFollowing: %sEE: ", eventExpectation.ID)
	logger.Info(s)
	return p.BusinessProcessTemplate.ExpectationsFollowingAnExpectationWaitingExpectations(eventExpectation, p.WaitingExpectations)
}

/*
	const expectationsFollowing = function (anEventExpectation) {
		return businessProcessTemplate.expectationsFollowingAnExpectationWaitingExpectations(anEventExpectation, waitingExpectations);
	};
	that.expectationsFollowing = expectationsFollowing;
	//
*/

//IsDone is
func (p *PersonBusinessProcess) IsDone() bool {
	logger.Info("Person Business Process -> IsDone ")
	return len(p.WaitingExpectations) == 0
}

/*
	const isDone = function () {
		return waitingExpectations.length === 0;
	};
	that.isDone = isDone;
*/

//ProcessTimeouts is
func (p *PersonBusinessProcess) ProcessTimeouts(schedulerEvent SchedulerEvent) {
	logger.Info("Person Business Process -> ProcessTimeouts ")
	foundTimedoutExpectations := false

	for {
		for _, anEventExpectation := range p.WaitingExpectations {
			//logger.log('debug', "Processing waiting expectation for " + anEventExpectation);
			foundTimeout := anEventExpectation.processTimeoutFor(p, schedulerEvent)
			if foundTimeout {
				foundTimedoutExpectations = true
			} else {
				foundTimedoutExpectations = false
			}
		}
		if !foundTimedoutExpectations {
			break
		}
	}

	/*
			const processTimeouts = function (anEvent) {

			logger.log('debug', "Processing timeouts for " + that.id);
			logger.log('debug', that.toString());
			let foundTimedoutExpectations = false;

			function checker() {
				waitingExpectations.forEach(function (anEventExpectation) {
					logger.log('debug', "Processing waiting expectation for " + anEventExpectation);
					let foundTimeout = anEventExpectation.processTimeoutFor(that, anEvent);
					if (foundTimeout) {
						foundTimedoutExpectations = true;
					} else {
						foundTimedoutExpectations = false;
					}
				});
			}
			do {
				checker();
			}
			while (foundTimedoutExpectations);

		};
		that.processTimeouts = processTimeouts;
	*/
}
