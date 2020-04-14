package wrappa

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"code.cloudfoundry.org/lager"

	// "github.com/concourse/concourse/atc/metric"
	"github.com/tedsuo/rata"
)

type ConcurrentRequestLimitFlag struct {
	Action string
	Limit  int
}

func (crl *ConcurrentRequestLimitFlag) UnmarshalFlag(value string) error {
	variable, expression, err := parseAssignment(value)
	if err != nil {
		return err
	}
	limit, err := strconv.Atoi(expression)
	if err != nil {
		return formatError(value, "limit must be an integer")
	}
	crl.Action = variable
	crl.Limit = limit
	return nil
}

func parseAssignment(value string) (string, string, error) {
	assignment := strings.Split(value, "=")
	if len(assignment) != 2 {
		return "", "", formatError(value, "value must be an assignment")
	}
	return assignment[0], assignment[1], nil
}

func formatError(value string, reason string) error {
	return fmt.Errorf("invalid concurrent request limit '%s': %s", value, reason)
}

type ConcurrencyLimitsWrappa struct {
	logger                  lager.Logger
	concurrentRequestLimits []ConcurrentRequestLimitFlag
}

func NewConcurrencyLimitsWrappa(
	logger lager.Logger,
	concurrentRequestLimits []ConcurrentRequestLimitFlag,
) Wrappa {
	return ConcurrencyLimitsWrappa{
		logger:                  logger,
		concurrentRequestLimits: concurrentRequestLimits,
	}
}

func (wrappa ConcurrencyLimitsWrappa) Wrap(handlers rata.Handlers) rata.Handlers {
	wrapped := rata.Handlers{}

	for _, limit := range wrappa.concurrentRequestLimits {
		for name, handler := range handlers {
			if limit.Action == name {
				wrapped[name] = wrapHandler(
					wrappa.logger,
					limit.Limit,
					handler,
				)
			} else {
				wrapped[name] = handler
			}
		}
	}

	return wrapped
}

func wrapHandler(logger lager.Logger, limit int, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})
}
