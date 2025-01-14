package registry

import "goken/pkg/errors"

type RegistorOptions struct {
	Address string `mapstructure:"address" json:"address"`
	Schema  string `mapstructure:"schema" json:"schema"`
}

func NewRegistorOptions() *RegistorOptions {
	return &RegistorOptions{}
}

func (o *RegistorOptions) Validate() []error {
	errs := []error{}
	if o.Address == "" || o.Schema == "" {
		errs = append(errs, errors.New("错误的RegistorOptions"))
	}
	return errs
}
