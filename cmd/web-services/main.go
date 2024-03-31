package webservices

import (
	"context"
	"log"
	"net/http"
	"time"

	"banking/api/handlers"
	"banking/api/responses"
	"banking/api/routes"
	"banking/configs"
	"banking/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestProm = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_histogram",
		Help:    "Histogram of the http request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	}, []string{"path", "method", "status"})
)

func FiberPrometheusMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()

	status := c.Response().StatusCode()
	httpRequestProm.WithLabelValues(c.Path(), c.Method(), http.StatusText(status)).Observe(float64(time.Since(start).Milliseconds()))

	return err
}

func Run() {
	app := fiber.New(
		fiber.Config{
			StrictRouting:     true,
			EnablePrintRoutes: true,
			CaseSensitive:     true,
		},
	)

	app.Use(FiberPrometheusMiddleware)
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	dbPool, err := db.NewPgConn(config)
	if err != nil {
		log.Fatalf("failed open connection to db: %v", err)
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("FAILED PING TO DB: %v", err)
	}

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	deps := handlers.Dependencies{
		Cfg:    config,
		DbPool: dbPool,
	}

	// load Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// register route in another package
	routes.RouteRegister(app, deps)

	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return responses.ReturnTheResponse(c, true, int(404), "Not Found", nil)
	})

	// Here we go!
	log.Fatalln(app.Listen(":" + config.APPPort))
}
