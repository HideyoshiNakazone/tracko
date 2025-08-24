package internal_errors

import "errors"

var ErrRestrictedConfigField = errors.New("attempted to modify restricted config field")
