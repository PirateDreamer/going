package going

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Houserqu/ginc"
	"github.com/PirateDreamer/going/ginx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 优雅开关机
func GranceRun(router *gin.Engine) {
	srv := &http.Server{
		Addr:    viper.GetString("server.addr"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func InitHttp() {
	ginx.InitContainer()
	var router *gin.Engine
	if !viper.GetBool("server.gin_not_default") {
		router = gin.Default()
	} else {
		router = gin.New()
	}
	router.Use(ginx.TraceMiddleware())
	ginx.R = router.Group("/api")

	ginx.AuthR = ginc.R.Group("/auth")
	// ginx.AuthR.Use(ginx.AuthMiddleware())
}
