package server

import (
	"app-invite-service/component/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	mmysql "github.com/golang-migrate/migrate/v4/database/mysql"
	"gorm.io/gorm"

	"app-invite-service/common"
	"app-invite-service/component"
	"app-invite-service/component/tokenprovider"
	"app-invite-service/config"
	"app-invite-service/middleware"
	"app-invite-service/module/user/usertransport/ginuser"
)

// Server represents server
type Server struct {
	ServerReady chan bool
	Port        int
	AppEnv      string
	SecretKey   string
	DBConn      *gorm.DB
	RedisConn   *redis.Client
	TokenConfig *tokenprovider.TokenConfig
}

func RunMigration(mysqlURL string) {
	sqlDB, err := sql.Open("mysql", mysqlURL)
	defer func() {
		_ = sqlDB.Close()
	}()
	if err != nil {
		log.Fatalf("cannot open migration database: %v", err)
	}

	driver, _ := mmysql.WithInstance(sqlDB, &mmysql.Config{})
	dbMigration, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf("cannot open migration database: %v", err)
	}

	if err := dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("fail to run migration: %v", err)
	}
}

// Start start http server
func Start(serverReady chan bool, cfg *config.Config) {
	// Create context that listens for the interrupt signal from the OS.
	// Reference: https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// global log level
	l := logger.New(cfg.Logger.Level)

	dbConn, err := gorm.Open(mysql.Open(cfg.MySQL.URL), &gorm.Config{})
	if err != nil {
		l.Fatal("app - Run - config.GetDBConn: %s", err)
	}

	tokenConfig, err := tokenprovider.NewTokenConfig(cfg.App.AtExpiry, cfg.App.RtExpiry)
	if err != nil {
		l.Fatal("app - Run - tokenprovider.NewTokenConfig: %s", err)
	}

	appCtx := component.NewAppContext(
		dbConn,
		redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			Password: cfg.Redis.Password,
			DB:       0,
		}),
		cfg.App.SecretKey,
		tokenConfig,
	)

	routes := InitRoutes(cfg, appCtx)

	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", cfg.App.Port),
		Handler: routes,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		l.Info("Server run on PORT :%d", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	if serverReady != nil {
		serverReady <- true
	}

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	l.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		l.Fatal("Server forced to shutdown: %v", err)
	}

	l.Info("Server exiting")
}

func InitRoutes(cfg *config.Config, appCtx component.AppContext) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	if cfg.App.ENV == common.AppEnvDev {
		gin.SetMode(gin.DebugMode)
		r.Use(gin.Logger())
	}

	r.Use(gin.Recovery())
	r.Use(middleware.Recover(appCtx))
	r.Use(middleware.CORSMiddleware(cfg))
	r.Use(middleware.GinContextToContextMiddleware())

	v1 := r.Group("/api/v1")

	v1.POST("/register", ginuser.Register(appCtx))
	v1.POST("/login", ginuser.Login(appCtx))
	v1.POST("/login/invitation", ginuser.LoginWithInviteToken(appCtx))

	v1.GET("/token/validation", ginuser.ValidateInvitationToken(appCtx))
	v1.GET(
		"/token/invitation",
		middleware.RequiredAuth(appCtx),
		middleware.RequiredAdmin(appCtx),
		ginuser.ListInvitationToken(appCtx),
	)
	v1.PATCH(
		"/token/invitation/:id",
		middleware.RequiredAuth(appCtx),
		middleware.RequiredAdmin(appCtx),
		ginuser.UpdateInvitationToken(appCtx),
	)

	v1.GET(
		"users/invitation",
		middleware.RequiredAuth(appCtx),
		middleware.RequiredAdmin(appCtx),
		ginuser.GenerateInviteToken(appCtx),
	)

	return r
}
