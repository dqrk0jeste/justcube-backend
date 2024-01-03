package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		ID:           uuid.New(),
		Username:     "darko",
		PasswordHash: "darko",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	fmt.Printf("%v", user)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NotZero(t, user.ID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
}
