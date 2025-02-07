package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/romulets/aipim/domain"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a mapping/parsing server",
	Run:   runServe,
}
var port int

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 8777, "Sets a port to serve on")
}

func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func toPainless(c *gin.Context) {
	var clm domain.CloudtrailLogMapping

	if err := c.BindJSON(&clm); err != nil {
		c.Data(
			http.StatusInternalServerError, "text/plain",
			[]byte(fmt.Sprintf("Error parsing JSON: %s", err)),
		)
		return
	}

	c.Data(http.StatusOK, "text/plain", []byte(clm.ToString()))
}

func fromPainless(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.Data(
			http.StatusInternalServerError, "text/plain",
			[]byte(fmt.Sprintf("Error reading Painless from body: %s", err)),
		)
		return
	}
	clm := domain.CloudtrailLogMapping{
		DefaultRelatedEntities: []string{},
		Sources:                []domain.MappedSource{},
	}
	if err := clm.Scan(string(raw)); err != nil {
		c.Data(
			http.StatusInternalServerError, "text/plain",
			[]byte(fmt.Sprintf("Error parsing Painless: %s", err)),
		)
		return
	}

	c.JSON(http.StatusOK, clm)
}

func runServe(cmd *cobra.Command, args []string) {
	router := gin.Default()
	router.Use(corsMiddleware())

	router.Use(static.Serve("/", static.LocalFile("./static", false)))

	router.GET("/api/health", healthcheck)
	router.POST("/api/mapping/to-painless", toPainless)
	router.POST("/api/mapping/from-painless", fromPainless)
	router.Run(fmt.Sprintf(":%d", port))
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
