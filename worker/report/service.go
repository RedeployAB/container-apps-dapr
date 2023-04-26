package report

import (
	"errors"
)

// Storer is the interface that wraps around method Store.
type Storer interface {
	Store(r Report) error
}

// Service is the interface that wraps around method Create.
type Service interface {
	Create(r Report) error
}

// service is the implementation of the Service interface.
type service struct {
	s Storer
}

// NewService creates a Service.
func NewService(s Storer) (*service, error) {
	if s == nil {
		return nil, errors.New("error")
	}

	return &service{
		s: s,
	}, nil
}

// Create a report and stores it at the target for the reporter.
func (s service) Create(r Report) error {
	if s.s == nil {
		return errors.New("storer is nil")
	}
	// Do reporting work...
	// Store the report.
	return s.s.Store(r)
}
