package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
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

func (s *Server) Start(ginMode string) error {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	h := &JQHandler{Config: s.Config}

	var err error

	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS, rice.LocateWorkingDirectory},
	}

	publicBox := conf.MustFindBox("public/root")

	tmpl := template.New("index.tmpl")
	tmpl.Delims("#{", "}")
	tmpl, err = tmpl.Parse(publicBox.MustString("index.tmpl"))
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(ginMode)
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

	jsBox := conf.MustFindBox("public/js")
	router.StaticFS("/js", jsBox.HTTPBox())
	router.StaticFS("/css", conf.MustFindBox("public/css").HTTPBox())
	router.StaticFS("/images", conf.MustFindBox("public/images").HTTPBox())
	router.StaticFS("/fonts", conf.MustFindBox("public/bower_components/bootstrap/dist/fonts").HTTPBox())

	workerFile := jsBox.MustString("worker-xquery.js")
	router.GET("/worker-xquery.js", func(c *gin.Context) {
		c.String(200, workerFile)
	})

	robotsFile := publicBox.MustString("robots.txt")
	router.GET("/robots.txt", func(c *gin.Context) {
		c.String(200, robotsFile)
	})

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
		if err := srv.ListenAndServe(); err != nil && !strings.Contains(err.Error(), "Server closed") {
			log.WithError(err).Fatal("error starting sever")
		}
	}()

	<-stop
	fmt.Println()
	ctx, cancel := context.WithTimeout(context.Background(), 28*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
