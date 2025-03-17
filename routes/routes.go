package routes

import (
	"database/sql"

	"github.com/BaikalMine/SongService/controllers"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// SetupRouter инициализирует маршруты API, включая swagger UI.
func SetupRouter(db *sql.DB, externalAPIUrl string) *gin.Engine {
	router := gin.Default()

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// Маршруты для песен
	router.GET("/songs", func(c *gin.Context) {
		controllers.GetSongs(c, db)
	})
	router.GET("/songs/:id/lyrics", func(c *gin.Context) {
		controllers.GetSongLyrics(c, db)
	})
	router.POST("/songs", func(c *gin.Context) {
		controllers.AddSong(c, db, externalAPIUrl)
	})
	router.PUT("/songs/:id", func(c *gin.Context) {
		controllers.UpdateSong(c, db)
	})
	router.DELETE("/songs/:id", func(c *gin.Context) {
		controllers.DeleteSong(c, db)
	})

	return router
}
