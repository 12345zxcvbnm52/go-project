package test

import "pkg/errors"

func WithMessage(err error, msg string) error {
	return errors.WithMessage(err, msg)
}

func WithStack(err error) error {
	return errors.WithStack(err)
}

func WithCode(err error) error {
	return errors.WithCode(0, "wa")
}

func Wrap() error {
	return errors.Wrap(errors.New("wqqtqw"), "qfw")
}
