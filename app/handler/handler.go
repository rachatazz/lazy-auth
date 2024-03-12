package handler

import (
	"net/http"
	"time"

	"lazy-auth/app/errs"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func HandleOk(c *gin.Context, data any, meta any) {
	method := c.Request.Method
	code := http.StatusOK
	if method == http.MethodPost {
		code = http.StatusCreated
	}

	c.JSON(
		code,
		gin.H{"status": "ok", "data": data, "meta": meta, "timestamp": time.Now()},
	)
}

func HandleError(c *gin.Context, err interface{}) {
	switch e := err.(type) {
	case errs.AppError:
		c.AbortWithStatusJSON(e.Code, gin.H{"status": "error", "message": e.Message, "timestamp": time.Now()})

	case validator.ValidationErrors:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": e.Error(), "timestamp": time.Now()})

	case error:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "unexpected error", "timestamp": time.Now()})
	}
}

type ValidateType int

const (
	ValidateQuery ValidateType = iota
	ValidateBody
)

func ValidationPipe(c *gin.Context, obj any, vType ValidateType) error {
	// https://gin-gonic.com/docs/examples/binding-and-validation/
	// https://github.com/go-playground/validator

	var err error
	switch vType {
	case ValidateQuery:
		err = c.ShouldBindQuery(obj)
	case ValidateBody:
		err = c.ShouldBind(obj)
	}

	return err
}
