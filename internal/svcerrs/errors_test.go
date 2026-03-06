package svcerrs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yaBliznyk/newsportal/internal/svcerrs"
)

func Test_InvalidFieldError_Is(t *testing.T) {
	t.Parallel()

	err := svcerrs.NewInvalidFieldError("test_field", "test reason")
	require.ErrorIs(t, err, svcerrs.ErrInvalidData)
}

func Test_InvalidFieldError_Error(t *testing.T) {
	t.Parallel()

	err := svcerrs.NewInvalidFieldError("test_field", "test reason")
	require.EqualError(t, err, "invalid field: \"test_field\", reason: test reason")
}

func Test_BusinessLogicError_Is(t *testing.T) {
	t.Parallel()

	err := svcerrs.NewBusinessLogicError("test_alias")
	require.ErrorIs(t, err, svcerrs.ErrBusinessLogic)
}

func Test_BusinessLogicError_Error(t *testing.T) {
	t.Parallel()

	err := svcerrs.NewBusinessLogicError("test_alias")
	require.EqualError(t, err, "wrong business logic: test_alias")
}
