//go:build wireinject

package dependency

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/weCredit/internal/database"
	"github.com/weCredit/internal/http/api"
	"github.com/weCredit/internal/http/controller"
	"github.com/weCredit/internal/pkg/config"
	"github.com/weCredit/internal/pkg/security"
	"github.com/weCredit/internal/pkg/util"
	"github.com/weCredit/internal/repository"
	"github.com/weCredit/internal/service"
)

func NewConfig(opt config.Options) (config.WeCreditConfig, error) {
	wire.Build(
		config.NewConfig, // This should call your NewConfig function which takes Options
	)
	return config.WeCreditConfig{}, nil
}

func NewDatabaseConfig(cfg config.WeCreditConfig) (*pgxpool.Pool, error) {
	wire.Build(
		database.NewDB,
	)
	return &pgxpool.Pool{}, nil
}

func NewWeCredit(cfg config.WeCreditConfig, db *pgxpool.Pool) (*api.WeCreditApi, error) {
	wire.Build(
		util.NewAppUtil,
		repository.NewTransactioner,
		security.NewJwtSecurityManager,
		repository.NewLoginCodeRepository,
		repository.NewUserRepository,

		service.NewUserService,

		controller.NewUserController,

		api.NewWeCreditApi,
	)
	return &api.WeCreditApi{}, nil
}
