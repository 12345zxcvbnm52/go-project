package inventorydata

import (
	"kenshop/pkg/errors"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
)

func GormInventoryErrHandle(err error) error {
	if err == nil {
		return nil
	}

	if errors.IfWithCoder(err) {
		return err
	}

	switch errors.Cause(err) {
	case gorm.ErrDuplicatedKey:
		err = errors.WithCoder(err, errors.CodeInventoryAlreadyExist, "")
	case gorm.ErrRecordNotFound:
		err = errors.WithCoder(err, errors.CodeInventoryNotFound, "")
	case redsync.ErrFailed:
		err = errors.WithCoder(err, errors.CodeRedlockLockFailed, "")
	case redsync.ErrLockAlreadyExpired:
		err = errors.WithCoder(err, errors.CodeRedlockUnlockFailed, "")
	case redsync.ErrExtendFailed:
		err = errors.WithCoder(err, errors.CodeRedlockExtendFailed, "")
	default:
		err = errors.WithCoder(err, errors.CodeInternalError, "")
	}
	return err
}
