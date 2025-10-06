package router

import (
	"net/http"

	"github.com/livingpool/config"
	"github.com/livingpool/constants"
	"github.com/livingpool/handler"
	"github.com/livingpool/service"
	"github.com/livingpool/views"
)

func Init(conf *config.Config) *http.ServeMux {
	router := http.NewServeMux()
	playerPool := service.NewPlayerPool(constants.PLAYER_POOL_CAP)
	gameHandler := handler.NewGameHandler(views.NewTemplates(), playerPool, &service.RealTimeProvider{})
	leaderboardHandler := handler.NewLeaderboardHandler(views.NewTemplates(), service.NewLeaderboard(conf.DB))

	// http.FS can be used to create a http.Filesystem
	var staticFS = http.FS(views.StaticFiles)
	fs := http.FileServer(staticFS)

	// Serve static files
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve all other requests
	router.HandleFunc("/", gameHandler.Home)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	router.HandleFunc("GET /return", gameHandler.ReturnHome)
	router.HandleFunc("POST /new", gameHandler.NewGame)
	router.HandleFunc("GET /check", gameHandler.CheckGuess)
	router.HandleFunc("POST /save-record", leaderboardHandler.SaveRecord)
	router.HandleFunc("GET /show-board", leaderboardHandler.ShowLeaderboard)

	return router
}
