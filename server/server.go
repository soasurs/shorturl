package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"shorturl/server/model"
	"shorturl/server/model/schema"
	"shorturl/transform/transform"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
)

func main() {
	// Dialect to gRPC service
	transformService := fmt.Sprint(os.Getenv("SHORTURL_TRANSFORM_SVC_SERVICE_HOST"), ":", os.Getenv("SHORTURL_TRANSFORM_SVC_SERVICE_PORT"))
	conn, err := grpc.Dial(transformService, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*5))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	transformClient := transform.NewTransformClient(conn)

	// Connect to PostgreSQL
	db, err := schema.Connect(os.Getenv("POSTGRES_DSN"))
	if err != nil {
		panic(err)
	}

	if err := schema.Migrate(context.Background(), db); err != nil {
		panic(err)
	}
	modelClient := model.NewClient(db)
	// Start HTTP server
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/resolve", func(c echo.Context) error {
		key := c.QueryParam("key")
		if key == "" {
			return c.String(http.StatusBadRequest, "null key")
		}

		origin, err := modelClient.ShortenMap().FindOne(c.Request().Context(), key)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"url": origin.Url,
		})
	})

	e.GET("/shorten", func(c echo.Context) error {
		url := c.QueryParam("url")
		if url == "" {
			return c.String(http.StatusBadRequest, "null url")
		}

		transformResponse, err := transformClient.Shorten(c.Request().Context(), &transform.ShortenRequest{Url: url})
		if err != nil {
			return err
		}

		shorten, err := modelClient.ShortenMap().Create(c.Request().Context(), url, transformResponse.Shorten)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"key": shorten.Key,
		})
	})

	e.Start(":3000")
}
