package topiclistenerinternal

import (
	"context"

	"github.com/ydb-platform/ydb-go-sdk/v3/internal/empty"
	"github.com/ydb-platform/ydb-go-sdk/v3/internal/topic/topicreadercommon"
	"github.com/ydb-platform/ydb-go-sdk/v3/topic/topicreader"
)

//go:generate mockgen -source event_handler.go -destination event_handler_mock_test.go -package topiclistenerinternal -write_package_comment=false --typed

type EventHandler interface {
	// OnStartPartitionSessionRequest called when server send start partition session request method.
	// You can use it to store read progress on your own side.
	// You must call event.Confirm(...) for start to receive messages from the partition.
	// You can set topiclistener.StartPartitionSessionConfirm for change default settings.
	//
	// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
	OnStartPartitionSessionRequest(ctx context.Context, event *PublicStartPartitionSessionEvent) error

	// OnReadMessages called with batch of messages. Max count of messages limited by internal buffer size
	//
	// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
	OnReadMessages(ctx context.Context, event *PublicReadMessages) error

	// OnStopPartitionSessionRequest called when the server send stop partition message.
	// It means that no more OnReadMessages calls for the partition session.
	// You must call event.Confirm() for allow the server to stop the partition session (if event.Graceful=true).
	// Confirm is optional for event.Graceful=false
	// The method can be called twice: with event.Graceful=true, then event.Graceful=false.
	// It is guaranteed about the method will be called least once.
	//
	// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
	OnStopPartitionSessionRequest(ctx context.Context, event *PublicStopPartitionSessionEvent) error
}

// PublicReadMessages
//
// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
type PublicReadMessages struct {
	PartitionSession topicreadercommon.PublicPartitionSession
	Batch            *topicreader.Batch
}

// PublicStartPartitionSessionEvent
//
// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
type PublicStartPartitionSessionEvent struct {
	PartitionSession topicreadercommon.PublicPartitionSession
	CommittedOffset  int64
	PartitionOffsets PublicOffsetsRange
	confirm          confirmStorage[PublicStartPartitionSessionConfirm]
}

func NewPublicStartPartitionSessionEvent(
	session topicreadercommon.PublicPartitionSession,
	committedOffset int64,
	partitionOffsets PublicOffsetsRange,
) *PublicStartPartitionSessionEvent {
	return &PublicStartPartitionSessionEvent{
		PartitionSession: session,
		CommittedOffset:  committedOffset,
		PartitionOffsets: partitionOffsets,
	}
}

// Confirm
//
// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
func (e *PublicStartPartitionSessionEvent) Confirm() {
	e.ConfirmWithParams(PublicStartPartitionSessionConfirm{})
}

func (e *PublicStartPartitionSessionEvent) ConfirmWithParams(p PublicStartPartitionSessionConfirm) {
	e.confirm.Set(p)
}

// PublicStartPartitionSessionConfirm
//
// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
type PublicStartPartitionSessionConfirm struct {
	readOffset   *int64
	CommitOffset *int64 ``
}

// WithReadOffet
//
// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
func (c PublicStartPartitionSessionConfirm) WithReadOffet(val int64) PublicStartPartitionSessionConfirm {
	c.readOffset = &val

	return c
}

// WithCommitOffset
//
// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
func (c PublicStartPartitionSessionConfirm) WithCommitOffset(val int64) PublicStartPartitionSessionConfirm {
	c.CommitOffset = &val

	return c
}

// PublicOffsetsRange
//
// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
type PublicOffsetsRange struct {
	Start int64
	End   int64
}

// PublicStopPartitionSessionEvent
//
// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
type PublicStopPartitionSessionEvent struct {
	PartitionSession topicreadercommon.PublicPartitionSession

	// Graceful mean about server is waiting for client finish work with the partition and confirm stop the work
	// if the field is false it mean about server stop lease the partition to the client and can assignee the partition
	// to other read session (on this or other connection).
	Graceful        bool
	CommittedOffset int64

	confirm confirmStorage[empty.Struct]
}

func NewPublicStopPartitionSessionEvent(
	partitionSession topicreadercommon.PublicPartitionSession,
	graceful bool,
	committedOffset int64,
) *PublicStopPartitionSessionEvent {
	return &PublicStopPartitionSessionEvent{
		PartitionSession: partitionSession,
		Graceful:         graceful,
		CommittedOffset:  committedOffset,
	}
}

// Confirm
//
// Experimental: https://github.com/ydb-platform/ydb-go-sdk/blob/master/VERSIONING.md#experimental
func (e *PublicStopPartitionSessionEvent) Confirm() {
	e.confirm.Set(empty.Struct{})
}
