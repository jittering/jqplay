package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/owenthereal/jqplay/config"
	"github.com/owenthereal/jqplay/jq"
	"github.com/sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
)

const (
	jqExecTimeout = 15 * time.Second
)

type JQHandlerContext struct {
	*config.Config
	JQ string
}

func (c *JQHandlerContext) Asset(path string) string {
	return fmt.Sprintf("%s/%s", c.AssetHost, path)
}

func (c *JQHandlerContext) ShouldInitJQ() bool {
	return c.JQ != ""
}

type JQHandler struct {
	Config *config.Config
}

func (h *JQHandler) handleIndex(c *gin.Context) {
	c.HTML(200, "index.tmpl", &JQHandlerContext{Config: h.Config})
}

func (h *JQHandler) handleJqPost(c *gin.Context) {
	var j *jq.JQ
	if err := c.BindJSON(&j); err != nil {
		err = fmt.Errorf("error parsing JSON: %s", err)
		h.logger(c).WithError(err).Info("error parsing JSON")
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), jqExecTimeout)
	defer cancel()

	// Evaling into ResponseWriter sets the status code to 200
	// appending error message in the end if there's any
	if err := j.Eval(ctx, c.Writer); err != nil {
		fmt.Fprint(c.Writer, err.Error())
		h.logger(c).WithError(err).Info("jq error")
	}
}

func (h *JQHandler) handleJqGet(c *gin.Context) {
	jq := &jq.JQ{
		J: c.Query("j"),
		Q: c.Query("q"),
	}

	var jqData string
	if err := jq.Validate(); err == nil {
		d, err := json.Marshal(jq)
		if err == nil {
			jqData = string(d)
		}
	}

	c.HTML(http.StatusOK, "index.tmpl", &JQHandlerContext{Config: h.Config, JQ: jqData})
}

func (h *JQHandler) handleJqSharePost(c *gin.Context) {
	c.String(http.StatusNotImplemented, "snippets not enabled")
}

func (h *JQHandler) handleJqShareGet(c *gin.Context) {
	id := c.Param("id")

	h.logger(c).WithField("id", id).Info("snippets not enabled")
	c.Redirect(http.StatusFound, "/")
}

func (h *JQHandler) logger(c *gin.Context) *logrus.Entry {
	l, _ := c.Get("logger")
	return l.(*logrus.Entry)
}

func shouldLogJQError(err error) bool {
	return err == jq.ExecTimeoutError || err == jq.ExecCancelledError
}
