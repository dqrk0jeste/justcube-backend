package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker("dadsdfdsfsdfdsfsdfsdfsdffrkjmbdx")
	require.NoError(t, err)

	username := "darko"
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	var payload *Payload

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker("dadsdfdsfsdfdsfsdfsdfsdffrkjmbdx")
	require.NoError(t, err)

	token, err := maker.CreateToken("darko", -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	var payload *Payload
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}
