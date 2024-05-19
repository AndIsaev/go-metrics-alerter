package storage

import "errors"

var ErrIncorrectMetricValue = errors.New("incorrect value for metric")
var ErrNotInitializedStorage = errors.New("the storage is not initialized")
var ErrKeyErrorStorage = errors.New("data not found")
