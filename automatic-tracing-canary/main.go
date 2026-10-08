package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"message": "success",
		})
	})
	router.GET("/health", func(c *gin.Context) {
		slog.Info("health request handled")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	chainTargetURL := os.Getenv("CHAIN_TARGET_URL")
	client := &http.Client{Transport: http.DefaultTransport}

	router.GET("/outbound", func(c *gin.Context) {
		request, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, "http://127.0.0.1:"+port+"/health", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response, err := client.Do(request)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		defer response.Body.Close()

		c.JSON(http.StatusOK, gin.H{"upstream_status": response.StatusCode})
	})
	router.GET("/chain", func(c *gin.Context) {
		if chainTargetURL == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "CHAIN_TARGET_URL is not configured",
				"runtime": "go",
			})
			return
		}

		request, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, chainTargetURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response, err := client.Do(request)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		defer response.Body.Close()

		c.JSON(response.StatusCode, gin.H{
			"runtime":         "go",
			"upstream_status": response.StatusCode,
		})
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
}
