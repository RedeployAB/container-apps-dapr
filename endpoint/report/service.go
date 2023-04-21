package report

import "errors"

// reporter is the interface that wraps around method Run.
type reporter interface {
	Run(report Report) error
}

// service is service containing settings and a reporter.
type service struct {
	r reporter
}

// NewService returns a new *service.
func NewService(r reporter) (*service, error) {
	if r == nil {
		return nil, errors.New("error creating service: reporter is nil")
	}
	return &service{
		r: r,
	}, nil
}

// Create a report.
func (s service) Create(report Report) error {
	if s.r == nil {
		return errors.New("error creating report: reporter is nil")
	}
	return s.r.Run(report)
}
