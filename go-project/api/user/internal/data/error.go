package userdata

import (
	errors "kenshop/pkg/error"

	"gorm.io/gorm"
)

func GrpcErrHandle(err error) error {
	if err == nil {
		return err
	}
	switch err {
	case gorm.ErrDuplicatedKey:
		return errors.WithCoder(err, errors.CodeUserAlreadyExist, "")
	case gorm.ErrRecordNotFound:
		return errors.WithCoder(err, errors.CodeUserNotFound, "")
	default:
		return errors.WithCoder(err, errors.CodeInternalError, "")
	}
}
