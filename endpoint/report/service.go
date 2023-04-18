package report

// reporter is the interface that wraps around method Run.
type reporter interface {
	Run(report Report) error
}

// service is service containing settings and a reporter.
type service struct {
	r reporter
}

// NewService returns a new *service.
func NewService(r reporter) *service {
	return &service{
		r: r,
	}
}

// Create a report.
func (s service) Create(report Report) error {
	return s.r.Run(report)
}
