package pipeline

import (
	"errors"
)

var ErrCanceled = errors.New("context canceled")

// Job is a long-running function which returns nil to indicate successful completion, or an error to indicate abnormal
// completion. Within the context of a Pipeline, a Job which returns an error will cause the entire pipeline to be
// canceled.
type Job func() error
