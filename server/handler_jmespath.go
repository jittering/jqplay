package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/owenthereal/jqplay/jmes"
	"github.com/sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
)

type JmesPathInput struct {
	Input  string
	Search string
}

// Run JMESPATH search
func (h *JQHandler) handleJmespathPost(c *gin.Context) {
	var j *JmesPathInput
	if err := c.BindJSON(&j); err != nil {
		err = fmt.Errorf("error parsing JSON: %s", err)
		logrus.WithError(err).Info("error parsing JSON")
		c.String(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), jqExecTimeout)
	defer cancel()

	// Evaling into ResponseWriter sets the status code to 200
	// appending error message in the end if there's any
	if err := jmes.Eval(ctx, j.Input, j.Search, c.Writer); err != nil {
		fmt.Fprint(c.Writer, err.Error())
		logrus.WithError(err).Info("jq error")
	}
	h.lastCommand = j.Search
}
