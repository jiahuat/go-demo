package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"

	whttp "westone.com/wscp-restful/pkg/cluster/delivery/http"
	metrics "westone.com/wscp-restful/pkg/metrics"
	option "westone.com/wscp-restful/pkg/option"
	swag "westone.com/wscp-restful/pkg/swagger"
	"westone.com/wscp-restful/pkg/tracing"
)

var LogLevel *string
var LogFile *string
var AlsoLogToStdErr *string

func execute() error {
	rootCmd := &cobra.Command{
		Use:   "wscp-restful",
		Short: "WSCP RESTFUL Server",
		Run: func(cmd *cobra.Command, args []string) {

			// init config
			var c option.Config
			if err := viper.Unmarshal(&c); err != nil {
				log.Fatalln("viper init config with err", err)
			}
			log.Printf("addr:%s, cluster name:%s \n", c.Http.Addr, c.Cluster.Name)
			startServer(&c)
		},
	}

	// flag
	flags := rootCmd.PersistentFlags()
	LogLevel = flags.String("log_level", "4", "Set log level 1 2 3 4")
	LogFile = flags.String("log_file", "./rest.log", "Set log file")
	AlsoLogToStdErr = flags.String("alsologtostderr", "false", "true or false")
	// init klog
	klog.InitFlags(nil)
	flags.Set("log_file", ".")
	flags.Set("alsologtostderr", "true")
	flag.Set("v", *LogLevel)
	log.Println("this is from cobra, loglevel ", *LogLevel)
	klog.V(1).Infoln("hello this is from klog Level: 1")
	klog.V(2).Infoln("hello this is from klog Level: 2")
	klog.V(3).Infoln("hello this is from klog Level: 3")
	klog.V(4).Infoln("hello this is from klog Level: 4")
	// viper
	cobra.OnInitialize(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("read config with err", err)
		}

		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Println("Config file changed:", e.Name)
		})
		viper.WatchConfig()
	})

	return rootCmd.Execute()
}

func startServer(config *option.Config) {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(metrics.GinMetricsMiddleware())
	router.Use(tracing.Jaeger())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.POST("/mt", func(c *gin.Context) {
		m := make(map[string]interface{})
		if err := c.ShouldBindJSON(&m); err != nil {
			log.Println("should build json with err", err)
		}
		log.Println("the response is : ", m)
		span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "doPing1")
		defer span.Finish()
		span.SetTag("response", "this is the response!!!!")
		c.JSON(200, m)
	})

	clusterRouter := router.Group("/cluster")
	// todo: handle err
	whttp.RegisHandlers(clusterRouter, config)
	// 条件加载
	swag.RegisSwagger(router)

	srv := &http.Server{
		Addr:    config.Http.Addr,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 2)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shuddownChan := make(chan struct{}, 1)
	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		close(shuddownChan)
	}()
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	case <-shuddownChan:
		log.Println("Server shutdown!")
	}
	log.Println("Server exiting")
}
