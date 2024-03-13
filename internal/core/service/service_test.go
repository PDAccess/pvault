package service

import (
	"context"
	"crypto/rand"
	"os"
	"testing"

	"github.com/pdaccess/pvault/internal/core/domain"
	"github.com/pdaccess/pvault/internal/core/ports"
	"github.com/pdaccess/pvault/internal/repo/mock"
	"github.com/rs/zerolog/log"
)

var (
	ctx  context.Context
	impl ports.Service
)

func TestMain(m *testing.M) {
	ctx = log.With().
		Str("component", "module").
		Logger().WithContext(context.Background())

	storeageKey := make([]byte, 32)
	if _, err := rand.Read(storeageKey); err != nil {
		log.Ctx(ctx).Err(err).Msg("random read")
		os.Exit(1)
	}

	var err error

	impl, err = New(storeageKey, mock.New(), mock.NewAllValidValidator())
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("service init")
		os.Exit(1)
	}

	ctx = context.WithValue(ctx, domain.UserTokenIn, "empty")

	m.Run()
}
