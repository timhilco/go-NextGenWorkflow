package databases

import "github.com/timhilco/go-NextGenWorkflow/util/logger"

type DatabaseContext struct {
	URL    string
	Logger *logger.HilcoLogger
}
