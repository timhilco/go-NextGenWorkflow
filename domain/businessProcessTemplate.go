package domain

import (
	"fmt"

	"github.com/timhilco/go-NextGenWorkflow/datastructures"
)

// BusinessProcessTemplate is a business process template
type BusinessProcessTemplate struct {
	ID                string             `bson:"id,omitempty"`
	EventExpectations []EventExpectation `bson:"eventExpectations,omitempty"`
}

//StepInterface  is an inteface for Step
type StepInterface interface {
}

//BreakStep is a type of Step
type BreakStep struct {
}

func (bs *BreakStep) ExecuteAnEventInAnEventExpectationForABusinessProcessTemplate(
	anEvent EventDefinition,
	anEventExpectation EventExpectation,
	aBusinessProcess BusinessProcesss) {

}

//CompositeStep is a type of Step
type CompositeStep struct {
}

//ConditionalStep is a type of Step
type ConditionalStep struct {
}
type SchedulerEvent struct {
	EventID       string
	EffectiveDate string
}

//Step is a step
type Step struct {
	Block string `bson:"block,omitempty"`
	Args  string `bson:"args,omitempty"`
	Code  string `bson:"code,omitempty"`
	F     string `bson:"f,omitempty"` // Build Function
	/*
	   // Create that object
	   const that = {};
	   // Define private properties
	   that.block = spec.block;
	   const args = spec.block.args;
	   const code = spec.block.code;
	   that.f = new Function(args, code);
	   // Define private methods
	*/
}

func (s *Step) ExecuteAnEventInThisForAnActiviation(jsonEvent datastructures.Message,
	anEventExpectation *EventExpectation,
	aBusinessProcess *PersonBusinessProcess) {
	aBusinessProcess.DontExpect(anEventExpectation)
	logger.Info("----------------------------------------------------------------------------")
	logger.Info("------- Step exectueAnEventInThisForAnActivation ---------------------------")
	logger.Info("------- that.f(anEvent, anEventExpectation, anWorkflow); -------------------")
	logger.Info("----------------------------------------------------------------------------")
	ee := anEventExpectation.nextExpectationsForAnWorkflow(aBusinessProcess)
	aBusinessProcess.ExpectAll(ee)
	if aBusinessProcess.IsDone() {
		aBusinessProcess.State = C_Status_Complete
	}
	/*
	   const executeAnEventInThisForAnActiviation = function (anEvent, anEventExpectation, anWorkflow) {
	       anWorkflow.dontExpectAnEventExpectation(anEventExpectation);
	       that.f(anEvent, anEventExpectation, anWorkflow);
	       anWorkflow.expectAll(anEventExpectation.nextExpectationsForAnWorkflow(anWorkflow));

	   };
	   that.executeAnEventInThisForAnActiviation = executeAnEventInThisForAnActiviation;
	   //
	*/
}
func (s *Step) String() string {
	t := "Step -> \n"
	return t
	/*
		t = t + s.id + "\n"
		t = t + ee.step.String() + "\n"
		t = t + ee.event.String()
		return t
	*/
	/*
	   const toString = function (iOffset) {
	       const offset = iOffset || 0;
	       let sOffset = "---";
	       for (let index = 0; index < offset; index++) {
	           sOffset = sOffset + "----";
	       }
	       return "Step -> \n" + sOffset
	       + "args: " + that.block.args + "\n" + sOffset
	       + "code: " + that.block.code + "<-- Step end";
	   };
	   that.toString = toString;

	*/
}

/*
func (s *Step) toDataObject() {

	   const toDataObject = function () {
	       const obj = {};
	       obj.block = that.block;
	       return obj;
	   };
	   that.toDataObject = toDataObject;

}
*/
// EventDefinition is
type EventDefinition struct {

	// that.id = spec.id;
	//that.effectiveDate = spec.effectiveDate || new Date();
	//that.businessProcessId = spec.businessProcessId;
	// Define private methods
	ID   string `bson:"id,omitempty"`
	Name string `bson:"name,omitempty"`
}

func (e *EventDefinition) match(jsonMessage datastructures.Message) bool {
	s1 := e.Name
	s2 := jsonMessage.Message.Header.EventName
	return s1 == s2
	/*
	   const match = function (anEventDefinition) {
	       return that.id === anEventDefinition.id;
	   };
	   that.match = match;
	*/
}

func (e *EventDefinition) String() string {
	s := "Event Definition -> "
	s = s + e.ID + ":" + e.Name + "\n"
	return s

	/*
	   const toString = function (iOffset) {
	       const offset = iOffset || 0;
	       let sOffset = "---";
	       for (let index = 0; index < offset; index++) {
	           sOffset = sOffset + "----";
	       }
	       let text = "";
	       text = "EventDefinition -->\n" + sOffset
	           + "id: "             + that.id + "\n" + sOffset
	           +  "effectiveDate: " + that.effectiveDate + "\n" + sOffset
	           +  "businessProcessId: " + that.businessProcessId +" <-- EventDefinition end"
	       return text;
	   };
	*/
} /*
func (e *EventDefinition) toDataObject() {


	   const toDataObject = function () {
	       const obj = {};
	       obj.id = that.id;
	       return obj;

}
*/
//EventExpectation is a event to expect
type EventExpectation struct {
	ID            string             `bson:"id"`
	Step          Step               `bson:"step"`
	Event         EventDefinition    `bson:"event"`
	TimeoutAction TimeoutAction      `bson:"timeoutAction"`
	Prerequisites []EventExpectation `bson:"prerequisites"`
	/*

	   = (businessProcessTemplate, spec) => {
	           // Create that object
	           const that = {};
	           // Define private properties
	           const event = businessProcessTemplate.event(spec.event);
	           const step = businessProcessTemplate.step(spec.step);
	           const timeoutAction = businessProcessTemplate.timeoutAction(businessProcessTemplate, spec.timeoutFunction);
	           that.id = spec.id;
	           let prerequisites = {};
	           if (spec.source === "MongoDB") {
	               logger.log("debug", "Creating expectation object from DB");
	               const prerequisitesArray = new List();
	               const specArray = spec.prerequisites;
	               specArray.forEach(function (each) {
	                   each.source = "MongoDB";
	                   const anExpection = businessProcessTemplate.eventExpectation(each);
	                   prerequisitesArray.push(anExpection);
	               });
	               prerequisites = prerequisitesArray;
	           } else {
	               prerequisites = spec.prerequisites || new List();
	           }
	       }

	*/
}

func (ee *EventExpectation) addAllPrerequisitesTo(aCollection []*EventExpectation) []*EventExpectation {
	//fmt.Printf("Entering addAllPrerequisitesTo -> %v \n ", aCollection)
	var includes bool = false
	for _, item := range aCollection {

		if ee.ID == item.ID {
			includes = true
		}
	}
	if includes {
		//fmt.Printf("Exiting addAllPrerequisitesTo Includes If Stops -> %v \n", ee)
		//eeFakeEE := EventExpectation{}
		//var fakeCollection []*EventExpectation = []*EventExpectation{&eeFakeEE}
		return aCollection
	}
	aCollection = append(aCollection, ee)
	for _, anEventExpectation := range ee.Prerequisites {

		//fmt.Printf("addAllPrerequisitesTo: in for loop -> %v -> %v \n ", ee, aCollection)
		aCollection = anEventExpectation.addAllPrerequisitesTo(aCollection)
	}
	var s string
	for _, ee := range aCollection {
		s = s + ee.ID + " "
	}
	s = "aCollection = [ " + s + " ]"
	ls := fmt.Sprintf("Business Process Template -> Exiting addAllPrerequisitesTo -> %s", s)
	logger.Info(ls)
	return aCollection
}

/*
const addAllPrerequisitesTo = function (aCollection) {
	const includes = aCollection.has(that, function (a, b) {
		return a.id === b.id;
		});
		if (includes) {
			return that;
		}
		aCollection.push(that);
		prerequisites.forEach(function (anEventExpectation) {
			anEventExpectation.addAllPrerequisitesTo(aCollection);
			});
			};
			that.addAllPrerequisitesTo = addAllPrerequisitesTo;
*/

func (ee *EventExpectation) allPrerequisites() []*EventExpectation {
	//fmt.Printf("Entering addAllPrerequisitesTo -> for: %v\n", ee)
	var result []*EventExpectation
	result = ee.addAllPrerequisitesTo(result)
	//fmt.Printf("Exiting addAllPrerequisitesTo -> for: %v\n", result)
	var result2 []*EventExpectation
	var s string
	// remove ee from result
	for _, e := range result {
		if ee.ID != e.ID {
			result2 = append(result2, e)
			s = s + e.ID + " "
		}
	}
	s = "[ " + s + " ]"
	ls := fmt.Sprintf("Business Process Template -> Exiting addAllPrerequisitesTo -> for: %s", s)
	logger.Info(ls)
	return result2

	/*
					function allPrerequisites() {
		       const result = new List();
		       addAllPrerequisitesTo(result);
		       result.delete(that);
		       return result;
		   }
	*/

}
func (ee *EventExpectation) String() string {
	s := "Event Expectation -> "
	s = s + ee.ID + "\n"
	s = s + ee.Step.String()
	s = s + ee.Event.String()
	p := ""
	for _, prereq := range ee.Prerequisites {
		p = p + prereq.ID + " "
	}
	s = s + "PreRegs: [" + p + "]\n"

	s = s + "-------- End Event Expectation -------------\n"
	return s

	/*
	   const toString = function (iOffset) {
	       const offset = iOffset || 0;
	       let sOffset = "---";
	       for (let index = 0; index < offset; index++) {
	           sOffset = sOffset + "----";
	       }
	       let text = "Prerequisites [ \n-----";
	       prerequisites.forEach(function (eventExpectation) {
	           text = text + eventExpectation.toString(offset+4) + "\n" + sOffset;
	       });
	       text = text + "] ";
	       return "-> EventExpectation(" + that.id + ")\n" + sOffset
	           + event.toString(offset) + " \n" + sOffset
	           + step.toString(offset) + " \n"  + sOffset
	           + timeoutAction.toString(offset) + " : \n" + sOffset
	           + text
	           + "<------- end eventExpection --------------------\n";
	   };
	   that.toString = toString;
	   //
	*/
}

/*
func (ee *EventExpectation) hasNoPrerequisites() {

	   const hasNoPrerequisites = function () {
	       if (prerequisites.length === 0) {
	           return true;
	       } else {
	           return false;
	       }
	   };
	   that.hasNoPrerequisites = hasNoPrerequisites;

}
*/
// hasPrerequisites
func (ee *EventExpectation) hasPrerequisites(eventExpection *EventExpectation) bool {

	var b bool = false
	for _, ee := range ee.Prerequisites {
		if ee.ID == eventExpection.ID {
			b = true
		}
	}
	//fmt.Printf("In ee.hasPrerequisites ->  %s hasPrequisites = %t \n", ee.ID, b)

	return b
	/*
	   const hasPrerequisites = function (eventExpectation) {
	       return prerequisites.has(eventExpectation, function (a, b) {
	           return a.id === b.id;
	       });
	   };
	   that.hasPrerequisites = hasPrerequisites;
	*/
}
func (ee *EventExpectation) hasNoPrerequisitesIn(eventExpectations []EventExpectation) bool {
	//fmt.Printf("Entering ee.hasNoPrerequisitesIn - ee= %s IN %v\n", ee.ID, eventExpectations)
	var b bool = true
	allPreReqs := ee.allPrerequisites()

	for _, preReq := range allPreReqs {

		var b2 bool = false
		for _, ee := range eventExpectations {
			if preReq.ID == ee.ID {
				b2 = true
			}
		}
		if b2 {
			b = false
		}
	}

	ls := fmt.Sprintf("Business Process Template -> Exiting ee.hasNoPrerequisitesIn - ee=%s result=%t ", ee.ID, b)
	logger.Info(ls)
	return b
	/*
	   const hasNoPrerequisitesIn = function (collection) {
	       const allPreReqs = allPrerequisites();
	       let bool = true;
	       allPreReqs.forEach(function (anEventExpectation) {
	           const b2 = collection.has(anEventExpectation, function (a, b) {
	               return a.id === b.id;
	           });
	           if (b2) {
	               bool = false;
	           }
	       });
	       return bool;

	   };
	   that.hasNoPrerequisitesIn = hasNoPrerequisitesIn;
	*/
}

//Expects is
func (ee *EventExpectation) Expects(jsonEvent datastructures.Message) bool {
	return ee.Event.match(jsonEvent)
	/*
	   const expects = function (anEvent) {
	       return event.match(anEvent);
	   };
	   that.expects = expects;
	*/
}

//ExecuteFor is
func (ee *EventExpectation) ExecuteFor(jsonEvent datastructures.Message, aBusinessProcess *PersonBusinessProcess) {

	ee.Step.ExecuteAnEventInThisForAnActiviation(jsonEvent, ee, aBusinessProcess)
}

/*
   const executeFor = function (anEvent, AnWorkflow) {
       step.executeAnEventInThisForAnActiviation(anEvent, that, AnWorkflow);
   };
   that.executeFor = executeFor;
   //
*/

func (ee *EventExpectation) processTimeoutFor(personBusinessProcess *PersonBusinessProcess, event SchedulerEvent) bool {

	//logger.log("debug", "in processTimeoutFor: " + that.toString());
	if ee.TimeoutAction.TriggerFunction != "" {
		//logger.log("debug", "in processTimeoutFor: timeoutAction");
		triggerStep := ee.TimeoutAction.HasTimeoutTriggered(event, ee, personBusinessProcess)
		if triggerStep {
			//logger.log("debug", "in processTimeoutFor- triggeringStep");
			message := datastructures.Message{}
			ee.TimeoutAction.ActionStep.ExecuteAnEventInThisForAnActiviation(message, ee, personBusinessProcess)
			//logger.log("debug", "in processTimeoutFor- After Step execute");
			return true
		} else {
			//logger.log("debug", "in processTimeoutFor: No timeout triggered");
			return false
		}
	} else {
		//logger.log("debug", "in processTimeoutFor: that.timeoutAction is undefined");
		return false
	}

	/*
	   const processTimeoutFor = function (AnWorkflow, anEvent) {
	       logger.log("debug", "in processTimeoutFor: " + that.toString());
	       if (timeoutAction) {
	           logger.log("debug", "in processTimeoutFor: timeoutAction");
	           const triggerStep = timeoutAction.hasTimeoutTriggered(anEvent, that, AnWorkflow);
	           if (triggerStep) {
	               logger.log("debug", "in processTimeoutFor- triggeringStep");
	               timeoutAction.actionStep.executeAnEventInThisForAnActiviation(anEvent, that, AnWorkflow);
	               logger.log("debug", "in processTimeoutFor- After Step execute");
	               return true;
	           } else {
	               logger.log("debug", "in processTimeoutFor: No timeout triggered");
	               return false;
	           }
	       } else {
	           logger.log("debug", "in processTimeoutFor: that.timeoutAction is undefined");
	           return false;
	       }

	   };
	   that.processTimeoutFor = processTimeoutFor;
	*/
}
func (ee *EventExpectation) nextExpectationsForAnWorkflow(aBusinessProcess *PersonBusinessProcess) []EventExpectation {

	collection := aBusinessProcess.ExpectationsFollowing(ee)

	return collection
	/*
	   const nextExpectationsForAnWorkflow = function (anWorkflow) {

	       const collection = anWorkflow.expectationsFollowing(that);
	       return collection;
	   };
	   that.nextExpectationsForAnWorkflow = nextExpectationsForAnWorkflow;
	*/
}

/*
func (ee *EventExpectation) toDataObject() {

	   const toDataObject = function () {
	       const obj = {};
	       obj.id = that.id;
	       obj.event = event.toDataObject();
	       obj.step = step.toDataObject();
	       obj.timeoutFunction = timeoutAction.toDataObject();
	       const expectationsArray = [];
	       prerequisites.forEach(function (item) {
	           expectationsArray.push(item.toDataObject());
	       });
	       obj.prerequisites = expectationsArray;
	       return obj;
	   };
	   that.toDataObject = toDataObject;

}
*/
//TimeoutAction is an actio take for a timeout
type TimeoutAction struct {
	TriggerFunction string `bson:"triggerFunction"`
	ActionStep      Step   `bson:"actionStep,omitempty"`
	/*

	   const that = {};
	   // Define private properties
	   if (spec) {
	       that.triggerFunction =
	           new Function(spec.trigger.args, spec.trigger.code);
	       const stepSpec = {
	           block: {
	               args: spec.action.args,
	               code: spec.action.code
	           }
	       };
	       that.actionStep = businessProcessTemplate.step(stepSpec);
	   } else {
	       logger.log("debug", "No timeout policy defined");
	   }

	*/
}

func (t *TimeoutAction) executeAnEventInThisForAnActiviation() {

	/*
	   const executeAnEventInThisForAnActiviation = function (anEvent, anEventExpectation, anWorkflow) {
	       anWorkflow.dontExpectAnEventExpectation(anEventExpectation);
	       that.f(anEvent, anEventExpectation, anWorkflow);
	       anWorkflow.expectAll(anEventExpectation.nextExpectationsForAnWorkflow(anWorkflow));
	       logger.log("debug", "--> Workflow after processing Step: SchedulerTask:executeAnEventInThisForAnActiviation ");
	       logger.log("debug", anWorkflow.toString());
	       logger.log("debug", "------------------------------------------");

	   };
	   that.executeAnEventInThisForAnActiviation = executeAnEventInThisForAnActiviation;
	   //
	*/
}
func (t *TimeoutAction) String() string {
	s := "Timeout -> goes here"
	return s
	/*
	   const toString = function (iOffset) {
	       const offset = iOffset || 1;
	       let sOffset = "---";
	       for (let index = 0; index < offset; index++) {
	           sOffset = sOffset + "----";
	       }
	       if (that.triggerFunction) {
	           const sFunction = that.triggerFunction.toString(offset);
	           const sStep = that.actionStep.toString(offset);
	           return "timeoutAction-> " + "\n"
	           + sOffset +  sFunction + "\n"
	           + sOffset + sStep;
	       } else {
	           return "timeoutAction-> ** No Policy **";
	       }
	   };
	*/
}

/*
func (t *TimeoutAction) toDataObject() {

	   const toDataObject = function () {
	       if (spec) {
	           const obj = {};
	           obj.trigger = spec.trigger;
	           obj.action = spec.action;
	           return obj;
	       }
	   };
	   that.toDataObject = toDataObject;

}
*/
func (t *TimeoutAction) HasTimeoutTriggered(event SchedulerEvent, ee *EventExpectation, personBusinessProcess *PersonBusinessProcess) bool {
	return true
	/*
	   const hasTimeoutTriggered = function (anEvent, anEventExpectation, anWorkflow) {
	       logger.log("debug", "in hasTimeoutTriggered: " + ">>" + that.toString());
	       if (that.triggerFunction) {
	           const bool = that.triggerFunction(anEvent, anEventExpectation, anWorkflow);
	           logger.log("debug", "in hasTimeoutTriggered: " + bool);
	           return bool;
	       }
	       return false;
	   };
	   that.hasTimeoutTriggered = hasTimeoutTriggered;
	*/

}

/*
   // Define private properties
   func  buildTemplate (buildSpec) {
       bptThis.businessProcessTemplateId = buildSpec.businessProcessTemplateId;
       if (buildSpec.source === "MongoDB") {
           logger.log("debug", "Creating businessProcessTemplate from DB Spec");
           printObject(buildSpec);
           const expecationArray = new List();
           const specArray = buildSpec.expectations;

           specArray.forEach(function (each) {
               each.source = "MongoDB";
               const anExpection = bptThis.eventExpectation(each);

               expecationArray.push(anExpection);
           });

           bptThis.expectations = expecationArray;
       } else {
           logger.log("debug", "Creating businessProcessTemplate from Memory Spec");
           bptThis.expectations = buildSpec.expectations;
       }
       return bptThis;
   }
*/

func (bpt *BusinessProcessTemplate) String() string {
	s := "Business Process Template -> goes here"
	return s
	/*
	   let text = "--> BusinessProcessTemplate: " + bptThis.businessProcessTemplateId + "\n---";
	   bptThis.expectations.forEach(function (eventExpectation) {
	       text = text + eventExpectation.toString(2) ;
	   });
	   text = text + "<---- end business process template ----- " ;
	   return text;
	*/
}

/*
func (bpt *BusinessProcessTemplate) toDataObject() {

	       const obj = {};
	       obj.businessProcessTemplateId = bptThis.businessProcessTemplateId;
	       const expectationsArray = [];
	       bptThis.expectations.forEach(function (item) {
	           expectationsArray.push(item.toDataObject());
	       });
	       obj.expectations = expectationsArray;
	       return obj;
	   };
	   bptThis.toDataObject = toDataObject;

}
*/
//InitialExpectations i
func (bpt *BusinessProcessTemplate) InitialExpectations() []EventExpectation {
	/* This is a quick fix TOTO */
	intial := bpt.EventExpectations[0]
	var initialExpectations []EventExpectation
	initialExpectations = append(initialExpectations, intial)
	return initialExpectations

	/*
	       const initial = bptThis.expectations.filter(function (each) {
	           return each.hasNoPrerequisites();
	       });
	       // logger.log("debug","initial expectations length = "+ initial.length);
	       return initial;
	   };
	*/
}

//ExpectationsFollowingAnExpectationWaitingExpectations is
func (bpt *BusinessProcessTemplate) ExpectationsFollowingAnExpectationWaitingExpectations(
	anEventExpectation *EventExpectation, waitingExpectations []EventExpectation) []EventExpectation {
	var s string
	var r []EventExpectation
	for _, expectation := range bpt.EventExpectations {
		//logger.Info("BusinessProcessTemplate:ExpectationsFollowingAnExpectationWaitingExpectations - Processing  " + expectation.ID)
		b1 := expectation.hasPrerequisites(anEventExpectation)
		b2 := expectation.hasNoPrerequisitesIn(waitingExpectations)
		if b1 && b2 {
			r = append(r, expectation)
			s = s + expectation.ID + " "
		}
	}
	s = "[ " + s + " ]"
	ls := fmt.Sprintf("Business Process Template -> ExpectationsFollowingAnExpectationWaitingExpectations - Entering WaitingExpectation %v : Result: %s", waitingExpectations, s)
	logger.Info(ls)
	return r
	/*
	   }
	    = function (anExpectation, waitingExpectations) {
	       // logger.log("info",this.toString());
	       // expectations.forEach (function (each) {
	       // logger.log("info",each.toString());
	       // const b1 = each.hasPrerequisites (anExpectation);
	       // const b2 = each.hasNoPrerequisitesIn (waitingExpectations);
	       // logger.log("info",b1+ ":"+ b2);
	       // }
	       // );
	       const r = bptThis.expectations.filter(function (expectation) {
	           return (expectation.hasPrerequisites(anExpectation) &&
	               expectation.hasNoPrerequisitesIn(waitingExpectations));
	       });
	       return r;
	   };
	   bptThis.expectationsFollowingAnExpectationWaitingExpectations = expectationsFollowingAnExpectationWaitingExpectations;
	   // End of private methods
	*/
}

//func CreateBakeTemplate ( context) {

//CreateBakeTemplate is a tempalte for baking a cake
func CreateBakeTemplate() BusinessProcessTemplate {

	businessProcessTemplateID := "Bake;"
	//const  businessProcessTemplate = businessProcessTemplateClass.businessProcessTemplate (null, context);

	/*
		args := "anEvent,anEventExpectation,anbusinessProcessTemplate"
		actionCode := "{console.log(\"Step Action-> GetIngredients timeout action\");}"
		triggerCode := "{   let endDate = anEvent.effectiveDate; " +
			"console.log (endDate); " +
			"let startDate = anWorkflow.effectiveDate; " +
			"console.log (startDate); " +
			"let days = endDate.getDate() - startDate.getDate(); " +
			"console.log (days); " +
			"return days > 1; }"

		triggerArgs := args
		actionArgs := args
	*/
	// add function to build timeout action
	/*
	   timeoutFunction.trigger = trigger;
	   //timeoutFunction.action = action;
	   let spec = {};
	   spec.timeoutFunction = timeoutFunction;
	*/
	// spec is for event expectation
	//bakeTimeoutAction := TimeoutAction{}
	getIngredientsEvent := EventDefinition{
		ID:   "GetIngredients",
		Name: "Cake:GetIngredients",
	}

	block := "printFunction(\"**** Step Action - Block Code -> Getting Ingredients *****\")"
	getIngrediantsStep := Step{
		Block: block,
	}
	getIngredientsEventExpectation := EventExpectation{
		ID:    "10",
		Step:  getIngrediantsStep,
		Event: getIngredientsEvent,
	}

	// wet := FCSEventExpectation event: FCSMixWetEvent new step: (FCSStep
	// block: [:event | Transcript cr; show: "Mixing wet"]) prerequisite: get.

	mixWetIngredientsEvent := EventDefinition{
		ID:   "MixWetIngredients",
		Name: "Cake:MixWetIngredients",
	}
	block2 := "printFunction(\"**** Step Action - Block Code -> Mixing Wet *****\")"
	mixWetIngredientsStep := Step{
		Block: block2,
	}
	prerequisites := []EventExpectation{getIngredientsEventExpectation}
	mixWetIngredientsEventExpectation := EventExpectation{
		ID:            "20",
		Step:          mixWetIngredientsStep,
		Event:         mixWetIngredientsEvent,
		Prerequisites: prerequisites,
	}
	// dry := FCSEventExpectation event: FCSMixDryEvent new step: (FCSStep
	// block: [:event | Transcript cr; show: "Mixing dry"]) prerequisite: get.

	mixDryIngredientsEvent := EventDefinition{
		ID:   "MixDryIngredients",
		Name: "Cake:MixDryIngredients",
	}
	block3 := "printFunction(\"**** Step Action - Block Code -->  Mixing Dry *****\")"
	mixDryIngredientStep := Step{
		Block: block3,
	}
	prerequisites2 := []EventExpectation{getIngredientsEventExpectation}

	mixDryIngredientsEventExpection := EventExpectation{
		ID:            "30",
		Step:          mixDryIngredientStep,
		Event:         mixDryIngredientsEvent,
		Prerequisites: prerequisites2,
	}

	// bake := FCSEventExpectation event: FCSBakeEvent new step: (FCSStep block:
	//[:event | Transcript cr; show: "Baking"]) prerequisites: (Array with: wet
	// * with: dry).

	bakeEvent := EventDefinition{
		ID:   "Bake",
		Name: "Cake:Bake",
	}

	block4 := "printFunction(\"**** Step Action - Block Code -> Baking *****\")"
	bakeStep := Step{
		Block: block4,
	}

	prerequisites3 := []EventExpectation{mixWetIngredientsEventExpectation, mixDryIngredientsEventExpection}

	bakeEventExpectation := EventExpectation{
		ID:            "40",
		Step:          bakeStep,
		Event:         bakeEvent,
		Prerequisites: prerequisites3,
	}

	expectations := []EventExpectation{
		getIngredientsEventExpectation,
		mixWetIngredientsEventExpectation,
		mixDryIngredientsEventExpection,
		bakeEventExpectation,
	}

	specTemplate := BusinessProcessTemplate{
		ID:                businessProcessTemplateID,
		EventExpectations: expectations,
	}
	return specTemplate

}
