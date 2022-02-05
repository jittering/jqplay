//go:generate rice embed-go
package server

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/owenthereal/jqplay/config"
	"github.com/owenthereal/jqplay/server/middleware"
	"github.com/pkg/browser"
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

	// var err error

	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS, rice.LocateWorkingDirectory},
	}

	publicBox := conf.MustFindBox("../web/public")

	// tmpl := template.New("index.tmpl")
	// tmpl.Delims("#{", "}")
	// tmpl, err = tmpl.Parse(publicBox.MustString("index.tmpl"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	gin.SetMode(ginMode)
	router := gin.New()
	router.Use(
		middleware.Timeout(25*time.Second),
		middleware.LimitContentLength(10),
		gin.Recovery(),
	)
	if s.Config.Verbose {
		router.Use(middleware.Logger())
	}
	// router.SetHTMLTemplate(tmpl)

	// router.StaticFS("/css", conf.MustFindBox("public/css").HTTPBox())
	// router.StaticFS("/images", conf.MustFindBox("../web/public/images").HTTPBox())
	// router.StaticFS("/fonts", conf.MustFindBox("public/bower_components/bootstrap/dist/fonts").HTTPBox())

	// dynamic routes
	router.GET("/jq/input", h.handleJqInput)
	router.POST("/jq/commandline", h.handleJqCommandLine)
	router.GET("/jq/version", h.handleJqVersion)
	router.GET("/jq", h.handleJqGet)
	router.POST("/jq", h.handleJqPost)
	// for jmespath
	router.POST("/jmespath", h.handleJmespathPost)

	// sharing, removed for now
	// router.POST("/s", h.handleJqSharePost)
	// router.GET("/s/:id", h.handleJqShareGet)

	// static files
	staticRoute(router, "/", publicBox, "index.html")
	staticRoute(router, "", publicBox, "material-icons.css")
	staticRoute(router, "", publicBox, "robots.txt")
	staticRoute(router, "", publicBox, "favicon.png")
	router.StaticFS("/images", conf.MustFindBox("../web/public/images").HTTPBox())
	router.StaticFS("/build", conf.MustFindBox("../web/public/build").HTTPBox())

	srv := &http.Server{
		Addr:    s.Config.Host + ":" + s.Config.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !strings.Contains(err.Error(), "Server closed") {
			log.WithError(err).Fatal("error starting sever")
		}
	}()

	url := "http://"
	if s.Config.Host == "0.0.0.0" {
		url += "127.0.0.1"
	} else {
		url += s.Config.Host
	}
	url += ":" + s.Config.Port
	if s.Config.NoOpen {
		fmt.Println("> server running at", url)
	} else {
		go func() {
			time.Sleep(time.Millisecond * 250)
			err := browser.OpenURL(url)
			if err != nil {
				fmt.Println("> server running at", url)
			} else {
				fmt.Println("> opening", url, "in default browser")
			}
		}()
	}

	defer func() {
		if h.lastCommand != "" {
			fmt.Println(h.lastCommand)
		}
	}()

	<-stop
	fmt.Println()
	ctx, cancel := context.WithTimeout(context.Background(), 28*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

func staticRoute(router *gin.Engine, path string, box *rice.Box, filename string) {
	if path == "" {
		path = "/" + filename
	}
	str := box.MustString(filename)
	ctype := mime.TypeByExtension(filepath.Ext(filename))
	router.GET(path, func(c *gin.Context) {
		c.Header("Content-Type", ctype)
		c.String(200, str)
	})
}
