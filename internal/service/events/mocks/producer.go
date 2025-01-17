// Code generated by http://github.com/gojuno/minimock (v3.3.13). DO NOT EDIT.

package mocks

//go:generate minimock -i HomeWork_1/internal/service/events.Producer -o producer.go -n Producer -p mocks

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// Producer implements events.Producer
type Producer struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcSendSyncMessage          func(msg []byte) (partition int32, offset int64, err error)
	inspectFuncSendSyncMessage   func(msg []byte)
	afterSendSyncMessageCounter  uint64
	beforeSendSyncMessageCounter uint64
	SendSyncMessageMock          mProducerSendSyncMessage
}

// NewProducer returns a mock for events.Producer
func NewProducer(t minimock.Tester) *Producer {
	m := &Producer{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SendSyncMessageMock = mProducerSendSyncMessage{mock: m}
	m.SendSyncMessageMock.callArgs = []*ProducerSendSyncMessageParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mProducerSendSyncMessage struct {
	optional           bool
	mock               *Producer
	defaultExpectation *ProducerSendSyncMessageExpectation
	expectations       []*ProducerSendSyncMessageExpectation

	callArgs []*ProducerSendSyncMessageParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// ProducerSendSyncMessageExpectation specifies expectation struct of the Producer.SendSyncMessage
type ProducerSendSyncMessageExpectation struct {
	mock      *Producer
	params    *ProducerSendSyncMessageParams
	paramPtrs *ProducerSendSyncMessageParamPtrs
	results   *ProducerSendSyncMessageResults
	Counter   uint64
}

// ProducerSendSyncMessageParams contains parameters of the Producer.SendSyncMessage
type ProducerSendSyncMessageParams struct {
	msg []byte
}

// ProducerSendSyncMessageParamPtrs contains pointers to parameters of the Producer.SendSyncMessage
type ProducerSendSyncMessageParamPtrs struct {
	msg *[]byte
}

// ProducerSendSyncMessageResults contains results of the Producer.SendSyncMessage
type ProducerSendSyncMessageResults struct {
	partition int32
	offset    int64
	err       error
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmSendSyncMessage *mProducerSendSyncMessage) Optional() *mProducerSendSyncMessage {
	mmSendSyncMessage.optional = true
	return mmSendSyncMessage
}

// Expect sets up expected params for Producer.SendSyncMessage
func (mmSendSyncMessage *mProducerSendSyncMessage) Expect(msg []byte) *mProducerSendSyncMessage {
	if mmSendSyncMessage.mock.funcSendSyncMessage != nil {
		mmSendSyncMessage.mock.t.Fatalf("Producer.SendSyncMessage mock is already set by Set")
	}

	if mmSendSyncMessage.defaultExpectation == nil {
		mmSendSyncMessage.defaultExpectation = &ProducerSendSyncMessageExpectation{}
	}

	if mmSendSyncMessage.defaultExpectation.paramPtrs != nil {
		mmSendSyncMessage.mock.t.Fatalf("Producer.SendSyncMessage mock is already set by ExpectParams functions")
	}

	mmSendSyncMessage.defaultExpectation.params = &ProducerSendSyncMessageParams{msg}
	for _, e := range mmSendSyncMessage.expectations {
		if minimock.Equal(e.params, mmSendSyncMessage.defaultExpectation.params) {
			mmSendSyncMessage.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmSendSyncMessage.defaultExpectation.params)
		}
	}

	return mmSendSyncMessage
}

// ExpectMsgParam1 sets up expected param msg for Producer.SendSyncMessage
func (mmSendSyncMessage *mProducerSendSyncMessage) ExpectMsgParam1(msg []byte) *mProducerSendSyncMessage {
	if mmSendSyncMessage.mock.funcSendSyncMessage != nil {
		mmSendSyncMessage.mock.t.Fatalf("Producer.SendSyncMessage mock is already set by Set")
	}

	if mmSendSyncMessage.defaultExpectation == nil {
		mmSendSyncMessage.defaultExpectation = &ProducerSendSyncMessageExpectation{}
	}

	if mmSendSyncMessage.defaultExpectation.params != nil {
		mmSendSyncMessage.mock.t.Fatalf("Producer.SendSyncMessage mock is already set by Expect")
	}

	if mmSendSyncMessage.defaultExpectation.paramPtrs == nil {
		mmSendSyncMessage.defaultExpectation.paramPtrs = &ProducerSendSyncMessageParamPtrs{}
	}
	mmSendSyncMessage.defaultExpectation.paramPtrs.msg = &msg

	return mmSendSyncMessage
}

// Inspect accepts an inspector function that has same arguments as the Producer.SendSyncMessage
func (mmSendSyncMessage *mProducerSendSyncMessage) Inspect(f func(msg []byte)) *mProducerSendSyncMessage {
	if mmSendSyncMessage.mock.inspectFuncSendSyncMessage != nil {
		mmSendSyncMessage.mock.t.Fatalf("Inspect function is already set for Producer.SendSyncMessage")
	}

	mmSendSyncMessage.mock.inspectFuncSendSyncMessage = f

	return mmSendSyncMessage
}

// Return sets up results that will be returned by Producer.SendSyncMessage
func (mmSendSyncMessage *mProducerSendSyncMessage) Return(partition int32, offset int64, err error) *Producer {
	if mmSendSyncMessage.mock.funcSendSyncMessage != nil {
		mmSendSyncMessage.mock.t.Fatalf("Producer.SendSyncMessage mock is already set by Set")
	}

	if mmSendSyncMessage.defaultExpectation == nil {
		mmSendSyncMessage.defaultExpectation = &ProducerSendSyncMessageExpectation{mock: mmSendSyncMessage.mock}
	}
	mmSendSyncMessage.defaultExpectation.results = &ProducerSendSyncMessageResults{partition, offset, err}
	return mmSendSyncMessage.mock
}

// Set uses given function f to mock the Producer.SendSyncMessage method
func (mmSendSyncMessage *mProducerSendSyncMessage) Set(f func(msg []byte) (partition int32, offset int64, err error)) *Producer {
	if mmSendSyncMessage.defaultExpectation != nil {
		mmSendSyncMessage.mock.t.Fatalf("Default expectation is already set for the Producer.SendSyncMessage method")
	}

	if len(mmSendSyncMessage.expectations) > 0 {
		mmSendSyncMessage.mock.t.Fatalf("Some expectations are already set for the Producer.SendSyncMessage method")
	}

	mmSendSyncMessage.mock.funcSendSyncMessage = f
	return mmSendSyncMessage.mock
}

// When sets expectation for the Producer.SendSyncMessage which will trigger the result defined by the following
// Then helper
func (mmSendSyncMessage *mProducerSendSyncMessage) When(msg []byte) *ProducerSendSyncMessageExpectation {
	if mmSendSyncMessage.mock.funcSendSyncMessage != nil {
		mmSendSyncMessage.mock.t.Fatalf("Producer.SendSyncMessage mock is already set by Set")
	}

	expectation := &ProducerSendSyncMessageExpectation{
		mock:   mmSendSyncMessage.mock,
		params: &ProducerSendSyncMessageParams{msg},
	}
	mmSendSyncMessage.expectations = append(mmSendSyncMessage.expectations, expectation)
	return expectation
}

// Then sets up Producer.SendSyncMessage return parameters for the expectation previously defined by the When method
func (e *ProducerSendSyncMessageExpectation) Then(partition int32, offset int64, err error) *Producer {
	e.results = &ProducerSendSyncMessageResults{partition, offset, err}
	return e.mock
}

// Times sets number of times Producer.SendSyncMessage should be invoked
func (mmSendSyncMessage *mProducerSendSyncMessage) Times(n uint64) *mProducerSendSyncMessage {
	if n == 0 {
		mmSendSyncMessage.mock.t.Fatalf("Times of Producer.SendSyncMessage mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmSendSyncMessage.expectedInvocations, n)
	return mmSendSyncMessage
}

func (mmSendSyncMessage *mProducerSendSyncMessage) invocationsDone() bool {
	if len(mmSendSyncMessage.expectations) == 0 && mmSendSyncMessage.defaultExpectation == nil && mmSendSyncMessage.mock.funcSendSyncMessage == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmSendSyncMessage.mock.afterSendSyncMessageCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmSendSyncMessage.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// SendSyncMessage implements events.Producer
func (mmSendSyncMessage *Producer) SendSyncMessage(msg []byte) (partition int32, offset int64, err error) {
	mm_atomic.AddUint64(&mmSendSyncMessage.beforeSendSyncMessageCounter, 1)
	defer mm_atomic.AddUint64(&mmSendSyncMessage.afterSendSyncMessageCounter, 1)

	if mmSendSyncMessage.inspectFuncSendSyncMessage != nil {
		mmSendSyncMessage.inspectFuncSendSyncMessage(msg)
	}

	mm_params := ProducerSendSyncMessageParams{msg}

	// Record call args
	mmSendSyncMessage.SendSyncMessageMock.mutex.Lock()
	mmSendSyncMessage.SendSyncMessageMock.callArgs = append(mmSendSyncMessage.SendSyncMessageMock.callArgs, &mm_params)
	mmSendSyncMessage.SendSyncMessageMock.mutex.Unlock()

	for _, e := range mmSendSyncMessage.SendSyncMessageMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.partition, e.results.offset, e.results.err
		}
	}

	if mmSendSyncMessage.SendSyncMessageMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmSendSyncMessage.SendSyncMessageMock.defaultExpectation.Counter, 1)
		mm_want := mmSendSyncMessage.SendSyncMessageMock.defaultExpectation.params
		mm_want_ptrs := mmSendSyncMessage.SendSyncMessageMock.defaultExpectation.paramPtrs

		mm_got := ProducerSendSyncMessageParams{msg}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.msg != nil && !minimock.Equal(*mm_want_ptrs.msg, mm_got.msg) {
				mmSendSyncMessage.t.Errorf("Producer.SendSyncMessage got unexpected parameter msg, want: %#v, got: %#v%s\n", *mm_want_ptrs.msg, mm_got.msg, minimock.Diff(*mm_want_ptrs.msg, mm_got.msg))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmSendSyncMessage.t.Errorf("Producer.SendSyncMessage got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmSendSyncMessage.SendSyncMessageMock.defaultExpectation.results
		if mm_results == nil {
			mmSendSyncMessage.t.Fatal("No results are set for the Producer.SendSyncMessage")
		}
		return (*mm_results).partition, (*mm_results).offset, (*mm_results).err
	}
	if mmSendSyncMessage.funcSendSyncMessage != nil {
		return mmSendSyncMessage.funcSendSyncMessage(msg)
	}
	mmSendSyncMessage.t.Fatalf("Unexpected call to Producer.SendSyncMessage. %v", msg)
	return
}

// SendSyncMessageAfterCounter returns a count of finished Producer.SendSyncMessage invocations
func (mmSendSyncMessage *Producer) SendSyncMessageAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSendSyncMessage.afterSendSyncMessageCounter)
}

// SendSyncMessageBeforeCounter returns a count of Producer.SendSyncMessage invocations
func (mmSendSyncMessage *Producer) SendSyncMessageBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSendSyncMessage.beforeSendSyncMessageCounter)
}

// Calls returns a list of arguments used in each call to Producer.SendSyncMessage.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmSendSyncMessage *mProducerSendSyncMessage) Calls() []*ProducerSendSyncMessageParams {
	mmSendSyncMessage.mutex.RLock()

	argCopy := make([]*ProducerSendSyncMessageParams, len(mmSendSyncMessage.callArgs))
	copy(argCopy, mmSendSyncMessage.callArgs)

	mmSendSyncMessage.mutex.RUnlock()

	return argCopy
}

// MinimockSendSyncMessageDone returns true if the count of the SendSyncMessage invocations corresponds
// the number of defined expectations
func (m *Producer) MinimockSendSyncMessageDone() bool {
	if m.SendSyncMessageMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.SendSyncMessageMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.SendSyncMessageMock.invocationsDone()
}

// MinimockSendSyncMessageInspect logs each unmet expectation
func (m *Producer) MinimockSendSyncMessageInspect() {
	for _, e := range m.SendSyncMessageMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to Producer.SendSyncMessage with params: %#v", *e.params)
		}
	}

	afterSendSyncMessageCounter := mm_atomic.LoadUint64(&m.afterSendSyncMessageCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.SendSyncMessageMock.defaultExpectation != nil && afterSendSyncMessageCounter < 1 {
		if m.SendSyncMessageMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to Producer.SendSyncMessage")
		} else {
			m.t.Errorf("Expected call to Producer.SendSyncMessage with params: %#v", *m.SendSyncMessageMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSendSyncMessage != nil && afterSendSyncMessageCounter < 1 {
		m.t.Error("Expected call to Producer.SendSyncMessage")
	}

	if !m.SendSyncMessageMock.invocationsDone() && afterSendSyncMessageCounter > 0 {
		m.t.Errorf("Expected %d calls to Producer.SendSyncMessage but found %d calls",
			mm_atomic.LoadUint64(&m.SendSyncMessageMock.expectedInvocations), afterSendSyncMessageCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *Producer) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockSendSyncMessageInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *Producer) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *Producer) minimockDone() bool {
	done := true
	return done &&
		m.MinimockSendSyncMessageDone()
}
