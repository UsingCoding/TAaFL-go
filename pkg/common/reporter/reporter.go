package reporter

import "github.com/sirupsen/logrus"

type Reporter interface {
	Info(...interface{})
	Error(...interface{})
	WithError(err error, args ...interface{})
}

func New(impl *logrus.Logger) Reporter {
	return &reporter{
		impl: impl,
	}
}

type reporter struct {
	impl *logrus.Logger
}

func (r *reporter) Info(args ...interface{}) {
	r.impl.Info(args...)
}

func (r *reporter) Error(args ...interface{}) {
	r.impl.Error(args)
}

func (r *reporter) WithError(err error, args ...interface{}) {
	r.impl.WithError(err).Error(args)
	r.impl.Fatal()
}
