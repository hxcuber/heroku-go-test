package repository

import (
	"github.com/pkg/errors"
)

var (
	errNestedTx = errors.New("db txn nested in db txn")
)
