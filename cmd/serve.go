package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
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

}

func fromPainless(c *gin.Context) {

}

func runServe(cmd *cobra.Command, args []string) {
	router := gin.Default()
	router.GET("/", healthcheck)
	router.POST("/api/mapping/to-painless", toPainless)
	router.POST("/api/mapping/from-painless", fromPainless)
	router.Run(fmt.Sprintf(":%d", port))
}
