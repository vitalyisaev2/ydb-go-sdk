package conn

import (
	"time"

	"google.golang.org/grpc"

	"github.com/ydb-platform/ydb-go-sdk/v3/internal/meta"
	"github.com/ydb-platform/ydb-go-sdk/v3/trace"
)

type Config interface {
	DialTimeout() time.Duration
	Trace() trace.Driver
	ConnectionTTL() time.Duration
	GrpcDialOptions() []grpc.DialOption
	Meta() meta.Meta
	UseDNSResolver() bool
}
