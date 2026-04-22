package domain

import "errors"

var ErrMessageNotFound = errors.New("message with given ID not found")
var ErrNotFound = errors.New("not found")
