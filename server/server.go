package server

import (
	"app-invite-service/component"
	"app-invite-service/component/tokenprovider"
	"app-invite-service/middleware"
	"app-invite-service/modules/user/usertransport/ginuser"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	mmysql "github.com/golang-migrate/migrate/v4/database/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
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

func (s *Server) RunMigration(dbConnectionStr string) {
	sqlDB, err := sql.Open("mysql", dbConnectionStr)
	if err != nil {
		log.Fatalln("cannot open migration database:", err)
	}

	driver, _ := mmysql.WithInstance(sqlDB, &mmysql.Config{})
	dbMigration, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalln("cannot open migration database:", err)
	}

	if err := dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Fail to run migration: ", err)
	}
}

// Start start http server
func (s *Server) Start() {
	// Create context that listens for the interrupt signal from the OS.
	// Reference: https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	if s.AppEnv == "dev" {
		gin.SetMode(gin.DebugMode)
		r.Use(gin.Logger())
	}

	appCtx := component.NewAppContext(s.DBConn, s.RedisConn, s.SecretKey, s.TokenConfig)
	r.Use(middleware.Recover(appCtx))

	v1 := r.Group("/api/v1")

	v1.POST("/register", ginuser.Register(appCtx))
	v1.POST("/login", ginuser.Login(appCtx))
	v1.POST("/login/invitation", ginuser.LoginWithInviteToken(appCtx))
	v1.GET("/token/validation", ginuser.ValidateInviteToken(appCtx))

	tasks := v1.Group("/users", middleware.RequiredAuth(appCtx))
	{
		tasks.POST("/invitation", ginuser.GenerateInviteToken(appCtx))
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Printf("Server run on PORT :%d\n", s.Port)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	if s.ServerReady != nil {
		s.ServerReady <- true
	}

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
