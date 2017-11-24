package task

import (
	"fmt"
	"time"
)

type TaskStatus string

const (
	Pending             TaskStatus = "Pending"
	Running             TaskStatus = "Running"
	ResponseImmediately TaskStatus = "ResponseImmediately" // error occurred in plugin
	Finishing           TaskStatus = "Finishing"           // plugin required finish pipeline normally
	Finished            TaskStatus = "Finished"
)

// The finite state machine of task status is:
//                   Pending
//                      +
//                      |
//                      v
//                   Running
//                      +
//                      |
//          +-----------------------+
//          v                       v
// ResponseImmediately          Finishing
//          +                       +
//          +-----------------------+
//                      |
//                      v
//                   Finished

type TaskResultCode uint

type TaskFinished func(task Task, originalStatus TaskStatus)
type TaskRecovery func(task Task, errorPluginName string) bool

type Task interface {
	// Finish sets status to `Finishing`
	Finish()
	// returns flag representing finished
	Finished() bool
	// ResultCode returns current result code
	ResultCode() TaskResultCode
	// Status returns current task status
	Status() TaskStatus
	// SetError sets error message, result code, and status to `ResponseImmediately`
	SetError(err error, resultCode TaskResultCode)
	// Error returns error message
	Error() error
	// StartAt returns task start time
	StartAt() time.Time
	// FinishAt return task finish time
	FinishAt() time.Time
	// AddFinishedCallback adds callback function executing after task status set to Finished
	// Callbacks are only used by plugin instead of model.
	AddFinishedCallback(name string, callback TaskFinished)
	// DeleteFinishedCallback deletes registered Finished callback function, don't call this in task finish callback
	DeleteFinishedCallback(name string)
	// AddRecoveryFunc adds callback function executing after task status set to `ResponseImmediately`,
	// after executing them the status of task will be recovered to `Running`
	AddRecoveryFunc(name string, taskRecovery TaskRecovery)
	// DeleteRecoveryFunc deletes registered recovery function
	DeleteRecoveryFunc(name string)
	// WithValue saves task-life-cycle value
	WithValue(key string, value interface{})
	// Value gets task-life-cycle value
	Value(key string) interface{}
	// Cancel returns a cancellation channel which could be closed to broadcast cancellation of task,
	// if a plugin needs relatively long time to wait I/O or anything else,
	// it should listen this channel to exit current plugin instance.
	Cancel() <-chan struct{}
	// CancelCause returns error message of cancellation
	CancelCause() error
	// Deadline returns deadline of task if the boolean flag set true, or it's not a task with deadline cancellation
	Deadline() (time.Time, bool)
}

////

var (
	Canceled         = fmt.Errorf("context canceled")
	DeadlineExceeded = fmt.Errorf("context deadline exceeded")
)

////

func ToString(value interface{}, maxLength uint64) string {
	switch value.(type) {
	case []uint8:
		return fmt.Sprintf(fmt.Sprintf("%%.%ds", maxLength), value)
	case string:
		return fmt.Sprintf(fmt.Sprintf("%%.%ds", maxLength), value)
	default:
		return fmt.Sprintf(fmt.Sprintf("%%.%dv", maxLength), value)
	}
}
