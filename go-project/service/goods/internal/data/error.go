package goodsdata

import (
	"kenshop/pkg/errors"

	"gorm.io/gorm"
)

func GormGoodsErrHandle(err error) error {
	if err == nil {
		return err
	}

	if errors.IfWithCoder(err) {
		return err
	}

	switch err {
	case gorm.ErrDuplicatedKey:
		err = errors.WithCoder(err, errors.CodeGoodsAlreadyExist, "")
	case gorm.ErrRecordNotFound:
		err = errors.WithCoder(err, errors.CodeGoodsNotFound, "")
	default:
		err = errors.WithCoder(err, errors.CodeInternalError, "")
	}
	return err
}

func GormBannerErrHandle(err error) error {
	if err == nil {
		return err
	}

	if errors.IfWithCoder(err) {
		return err
	}

	switch err {
	case gorm.ErrDuplicatedKey:
		err = errors.WithCoder(err, errors.CodeBannerAlreadyExist, "")
	case gorm.ErrRecordNotFound:
		err = errors.WithCoder(err, errors.CodeBannerNotFound, "")
	default:
		err = errors.WithCoder(err, errors.CodeInternalError, "")
	}
	return err
}

func GormCategoryErrHandle(err error) error {
	if err == nil {
		return err
	}

	if errors.IfWithCoder(err) {
		return err
	}

	switch err {
	case gorm.ErrDuplicatedKey:
		err = errors.WithCoder(err, errors.CodeCategoryAlreadyExist, "")
	case gorm.ErrRecordNotFound:
		err = errors.WithCoder(err, errors.CodeCategoryNotFound, "")
	default:
		err = errors.WithCoder(err, errors.CodeInternalError, "")
	}
	return err
}

func GormBrandErrHandle(err error) error {
	if err == nil {
		return err
	}

	if errors.IfWithCoder(err) {
		return err
	}

	switch err {
	case gorm.ErrDuplicatedKey:
		err = errors.WithCoder(err, errors.CodeBrandAlreadyExist, "")
	case gorm.ErrRecordNotFound:
		err = errors.WithCoder(err, errors.CodeBrandNotFound, "")
	default:
		err = errors.WithCoder(err, errors.CodeInternalError, "")
	}
	return err
}

func GormCategoryBrandErrHandle(err error) error {
	if err == nil {
		return err
	}

	if errors.IfWithCoder(err) {
		return err
	}

	switch err {
	case gorm.ErrDuplicatedKey:
		err = errors.WithCoder(err, errors.CodeCategoryBrandAlreadyExist, "")
	case gorm.ErrRecordNotFound:
		err = errors.WithCoder(err, errors.CodeCategoryBrandNotFound, "")
	default:
		err = errors.WithCoder(err, errors.CodeInternalError, "")
	}
	return err
}
