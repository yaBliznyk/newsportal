package svcerrs_test

import (
	"errors"
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

func Test_NewDetailsError(t *testing.T) {
	t.Parallel()

	testErr := errors.New("test error")
	err := svcerrs.NewDetailsError(testErr, []svcerrs.Detail{
		{
			Code:   "test_code",
			Data:   map[string]any{"test_key": "test_value"},
			Reason: "test reason",
		},
	})
	require.Error(t, err)
	var detailsError *svcerrs.DetailsError
	require.ErrorAs(t, err, &detailsError)
	require.Equal(t, testErr, detailsError.Err)
	require.Equal(t, []svcerrs.Detail{
		{
			Code:   "test_code",
			Data:   map[string]any{"test_key": "test_value"},
			Reason: "test reason",
		},
	}, detailsError.Details)
}

func Test_NewDetailsError_Is(t *testing.T) {
	t.Parallel()

	testErr := errors.New("test error")
	err := svcerrs.NewDetailsError(testErr, []svcerrs.Detail{
		{
			Code:   "test_code",
			Data:   map[string]any{"test_key": "test_value"},
			Reason: "test reason",
		},
	})
	require.ErrorIs(t, err, testErr)
}

func Test_NewDetailsError_Error(t *testing.T) {
	t.Parallel()

	testErr := errors.New("test error")
	err := svcerrs.NewDetailsError(testErr, []svcerrs.Detail{
		{
			Code:   "test_code",
			Data:   map[string]any{"test_key": "test_value"},
			Reason: "test reason",
		},
	})
	require.EqualError(t, err, "test error")
}

func Test_NewDetailsError_As(t *testing.T) {
	t.Parallel()

	testErr := svcerrs.NewBusinessLogicError("test_alias")
	err := svcerrs.NewDetailsError(testErr, []svcerrs.Detail{
		{
			Code:   "test_code",
			Data:   map[string]any{"test_key": "test_value"},
			Reason: "test reason",
		},
	})
	var detailsError *svcerrs.DetailsError
	require.ErrorAs(t, err, &detailsError)
	var businessLogicError *svcerrs.BusinessLogicError
	require.ErrorAs(t, err, &businessLogicError)
}
