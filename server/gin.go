package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	mongdblib "github.com/opensourceways/mongodb-lib/mongodblib"
	"github.com/opensourceways/server-common-lib/interrupts"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/opensourceways/software-package-server/config"
	"github.com/opensourceways/software-package-server/docs"
	softwarepkgapp "github.com/opensourceways/software-package-server/softwarepkg/app"
	"github.com/opensourceways/software-package-server/softwarepkg/controller"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/clavalidatorimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/messageimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/pkgmanagerimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/repositoryimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/sigvalidatorimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/softwarepkgadapter"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/translationimpl"
	"github.com/opensourceways/software-package-server/softwarepkg/infrastructure/useradapterimpl"
)

func StartWebServer(port int, timeout time.Duration, cfg *config.Config) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logRequest())

	setRouter(r, cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.ListenAndServe(srv, timeout)
}

// setRouter init router
func setRouter(engine *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "Software Package"
	docs.SwaggerInfo.Description = "set header: 'PRIVATE-TOKEN=xxx'"

	v1 := engine.Group(docs.SwaggerInfo.BasePath)
	setApiV1(v1, cfg)

	engine.UseRawPath = true
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func setApiV1(v1 *gin.RouterGroup, cfg *config.Config) {
	initSoftwarePkgService(v1, cfg)

	controller.AddRouteForSigController(v1, sigvalidatorimpl.SigValidator())

	controller.AddRouteForCLAController(
		v1, clavalidatorimpl.Instance(),
	)
}

func initSoftwarePkgService(v1 *gin.RouterGroup, cfg *config.Config) {
	repo := softwarepkgadapter.NewsoftwarePkgAdapter(
		mongdblib.DAO(cfg.Mongo.Collections.SoftwarePkg),
	)

	controller.AddRouteForSoftwarePkgController(
		v1,
		softwarepkgapp.NewSoftwarePkgService(
			repo,
			pkgmanagerimpl.Instance(),
			messageimpl.Producer(),
		),
		useradapterimpl.UserAdapter(),
	)

	controller.AddRouteForSoftwarePkgCommentController(
		v1,
		softwarepkgapp.NewSoftwarePkgCommentAppService(
			repo,
			translationimpl.Translation(),
			repositoryimpl.NewSoftwarePkgComment(&cfg.Postgresql.Config),
		),
	)
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		logrus.Infof(
			"| %d | %d | %s | %s |",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}
