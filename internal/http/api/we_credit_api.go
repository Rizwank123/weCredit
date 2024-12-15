package api

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/weCredit/internal/http/controller"
	"github.com/weCredit/internal/pkg/config"
)

type WeCreditApi struct {
	cfg            config.WeCreditConfig
	UserController controller.UserController
}

// NewWeChatApi creates a new WeCredit instance
//
//	@title						WeChat API
//	@version					1.0
//	@description				WeChat application's set of APIs
//	@termsOfService				https://example.com/terms
//	@contact.name				Mohammad Developer
//	@contact.url				https://rizwank123.github.io
//	@contact.email				md.rizwank431@gmail.com
//	@host						localhost:7700
//	@BasePath					/api/v1
//	@schemes					http https
//	@securityDefinitions.apiKey	JWT
//	@in							header
//	@name						Authorization
func NewWeCreditApi(cfg config.WeCreditConfig, uc controller.UserController) *WeCreditApi {
	return &WeCreditApi{
		cfg:            cfg,
		UserController: uc,
	}
}

func (b WeCreditApi) SetupRoutes(e *echo.Echo) {
	apiV1 := e.Group("/api/v1")

	auth := echojwt.JWT([]byte(b.cfg.AuthSecret))

	userApi := apiV1.Group("/users")
	userApi.POST("/login", b.UserController.Login)
	userApi.POST("", b.UserController.RegisterUser)
	userApi.POST("/init/login", b.UserController.InitLogin)
	secureApi := apiV1.Group("/users")
	secureApi.Use(auth)
	secureApi.GET("/:id", b.UserController.FindByID)

}
