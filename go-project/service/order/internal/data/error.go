package orderdata

import (
	"kenshop/pkg/errors"

	"gorm.io/gorm"
)

func GormOrderErrHandle(err error) error {
	if err == nil {
		return nil
	}

	if errors.IfWithCoder(err) {
		return err
	}

	switch errors.Cause(err) {
	case gorm.ErrDuplicatedKey:
		err = errors.WithCoder(err, errors.CodeOrderAlreadyExist, "")
	case gorm.ErrRecordNotFound:
		err = errors.WithCoder(err, errors.CodeOrderNotFound, "")
	default:
		err = errors.WithCoder(err, errors.CodeInternalError, "")
	}
	return err
}
