package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/mirror520/notification"
	"github.com/mirror520/notification/conf"
	"github.com/mirror520/notification/message"
	"github.com/mirror520/notification/provider"
)

func main() {
	app := &cli.App{
		Name:  "notification",
		Usage: "sends out alerts and messages to users",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Usage:   "work directory",
				EnvVars: []string{"NOTIFICATION_PATH"},
				Value:   ".",
			},
			&cli.IntFlag{
				Name:    "port",
				Usage:   "service port",
				Value:   8080,
				EnvVars: []string{"NOTIFICATION_PORT"},
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cli *cli.Context) error {
	path := cli.Path("path")

	f, err := os.Open(path + "/config.yaml")
	if err != nil {
		f, err = os.Open(path + "/config.example.yaml")
		if err != nil {
			return err
		}
	}

	var cfg *conf.Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return err
	}

	log, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer log.Sync()

	zap.ReplaceGlobals(log)

	providers := make(map[string]message.Provider)
	for name, profile := range cfg.Providers {
		provider, err := provider.NewProvider(profile)
		if err != nil {
			log.Error(err.Error(), zap.String("provider", name))
			continue
		}

		providers[name] = provider
	}

	if len(providers) == 0 {
		return errors.New("no provider")
	}

	r := gin.Default()
	r.Use(cors.Default())

	svc := notification.NewService(providers)

	apiV1 := r.Group("/notification/v1")
	{
		// POST /send
		{
			endpoint := notification.SendEndpoint(svc)
			apiV1.POST("/send", notification.SendHandler(endpoint))
		}

		// GET /providers/:id/credit
		{
			endpoint := notification.CreditEndpoint(svc)
			apiV1.GET("/providers/:id/credit", notification.CreditHandler(endpoint))
		}
	}

	port := cli.Int("port")
	go r.Run(":" + strconv.Itoa(port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit
	log.Info(sign.String())

	return nil
}
