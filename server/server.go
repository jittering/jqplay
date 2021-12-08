package server

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/owenthereal/jqplay/config"
	"github.com/owenthereal/jqplay/server/middleware"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
)

func New(c *config.Config) *Server {
	return &Server{c}
}

type Server struct {
	Config *config.Config
}

func (s *Server) Start() error {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	h := &JQHandler{Config: s.Config}

	var err error

	tmpl := template.New("index.tmpl")
	tmpl.Delims("#{", "}")
	tmpl, err = tmpl.ParseFiles("public/index.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(
		middleware.Timeout(25*time.Second),
		middleware.LimitContentLength(10),
		middleware.Secure(s.Config.IsProd()),
		middleware.RequestID(),
		middleware.Logger(),
		gin.Recovery(),
	)
	router.SetHTMLTemplate(tmpl)

	router.Static("/js", "public/js")
	router.Static("/css", "public/css")
	router.Static("/images", "public/images")
	router.Static("/fonts", "public/bower_components/bootstrap/dist/fonts")
	router.StaticFile("/worker-xquery.js", "public/bower_components/ace-builds/src-min-noconflict/worker-xquery.js")
	router.StaticFile("/robots.txt", "public/robots.txt")

	router.GET("/", h.handleIndex)
	router.GET("/jq", h.handleJqGet)
	router.POST("/jq", h.handleJqPost)
	router.POST("/s", h.handleJqSharePost)
	router.GET("/s/:id", h.handleJqShareGet)

	srv := &http.Server{
		Addr:    ":" + s.Config.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.WithError(err).Fatal("error starting server")
		}
	}()

	<-stop
	log.Println("\nShutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 28*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
