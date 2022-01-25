package domain

import "github.com/timhilco/go-NextGenWorkflow/datastructures"

type StepFactory struct{}

func (f StepFactory) BuildStep(specs map[string]string, action func()) IStep {
	step := &Step{Block: specs["block"]}
	return step
}
func (f StepFactory) BuildCompositeStep(steps []IStep) IStep {

	step := &CompositeStep{Steps: steps}
	return step
}

//StepInterface  is an inteface for Step
type IStep interface {
	ExecuteAnEventInThisForAnActiviation(jsonEvent datastructures.Message,
		anEventExpectation *EventExpectation,
		aBusinessProcess *PersonBusinessProcess)
	String() string
}

//BreakStep is a type of Step
type BreakStep struct {
}

func (bs *BreakStep) ExecuteAnEventInThisForAnActiviation(jsonEvent datastructures.Message,
	anEventExpectation *EventExpectation,
	aBusinessProcess *PersonBusinessProcess) {
	logger.Info("In (bs *BreakStep) ExecuteAnEventInThisForAnActiviation")
}
func (bs *BreakStep) String() string {
	text := "aBreakStep"
	return text
}

//CompositeStep is a type of Step
type CompositeStep struct {
	Steps []IStep
}

func (cs *CompositeStep) ExecuteAnEventInThisForAnActiviation(jsonEvent datastructures.Message,
	anEventExpectation *EventExpectation,
	aBusinessProcess *PersonBusinessProcess) {
	for _, step := range cs.Steps {
		step.ExecuteAnEventInThisForAnActiviation(jsonEvent,
			anEventExpectation,
			aBusinessProcess)

	}

}
func (cs *CompositeStep) String() string {
	text := "aCompositeStep"
	return text
}

//ConditionalStep is a type of Step
type ConditionalStep struct {
	If   bool
	Then EventExpectation
	Else EventExpectation
}

func (cs *ConditionalStep) ExecuteAnEventInThisForAnActiviation(jsonEvent datastructures.Message,
	anEventExpectation *EventExpectation,
	aBusinessProcess *PersonBusinessProcess) {
	logger.Info("In (cs *ConditionalStep) ExecuteAnEventInThisForAnActiviation")
}
func (cs *ConditionalStep) String() string {
	text := "aConditionalStep"
	return text
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
	logger.Info(s.Block)
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
