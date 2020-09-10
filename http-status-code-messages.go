package go_rnd_common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ThrowStatusOk(i interface{}, c *gin.Context) {
	if i != nil {
		c.JSON(http.StatusOK, i)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "OK",
	})
}

type HttpError struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status"`
}

func ThrowStatusBadRequest(msg string, c *gin.Context) {
	err := HttpError{
		Error:      msg,
		StatusCode: http.StatusBadRequest,
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, err)
}

func ThrowStatusInternalServerError(msg string, c *gin.Context) {
	err := HttpError{
		Error:      msg,
		StatusCode: http.StatusInternalServerError,
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, err)
}

func ThrowUniqueViolationErr(err string, c *gin.Context) {
	formattedErr := getColumnNameForDbErr(err)

	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": "ERR_DUPLICATE_ENTRY_" + formattedErr,
	})
}

func ThrowStatusUnauthorized(errMsg string, c *gin.Context) {
	err := HttpError{
		Error:      errMsg,
		StatusCode: http.StatusUnauthorized,
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, err)
}

func RecordNotFound(err error) bool {
	return err.Error() == "record not found"
}

var beautifiedErrorKeys = map[string]string{
	"user_email":    "EMAIL",
	"user_steam_id": "STEAM_ID",
}

func getColumnNameForDbErr(title string) string {
	errToReturn := strings.Split(title, "(")[1]
	errToReturn = errToReturn[:len(errToReturn)-2]

	if val, ok := beautifiedErrorKeys[errToReturn]; ok {
		return val
	}

	return errToReturn
}

type HttpErrorWithErrorSlug struct {
	Error      string `json:"error"`
	ErrorSlug  string `json:"error_slug"`
	StatusCode int    `json:"status_code"`
}

func ThrowAnErrorWithErrorSlug(errName string, errSlug string, statusCode int, ctx *gin.Context) {
	e := HttpErrorWithErrorSlug{
		Error:      errName,
		ErrorSlug:  errSlug,
		StatusCode: statusCode,
	}

	ctx.AbortWithStatusJSON(statusCode, e)
}
