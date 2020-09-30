package go_rnd_common

import (
	"errors"
	"github.com/croatiangrn/scill_errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
)

var (
	ErrGenericErr      = errors.New("generic_err")
	ErrDBEmpty         = errors.New("db_empty")
	ErrLanguageIDEmpty = errors.New("default_language_id_empty")
)

// Deprecated: Use RND.ThrowStatusOK instead
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

// Deprecated: Use RND.HttpErrorWithSlug instead
func ThrowStatusBadRequest(msg string, c *gin.Context) {
	err := HttpError{
		Error:      msg,
		StatusCode: http.StatusBadRequest,
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, err)
}

// Deprecated: Use RND.HttpErrorWithSlug instead
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

// Deprecated: Use RND.ThrowStatusUnauthorized instead
func ThrowStatusUnauthorized(errMsg string, c *gin.Context) {
	err := HttpError{
		Error:      errMsg,
		StatusCode: http.StatusUnauthorized,
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, err)
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

type RND struct {
	DB *gorm.DB

	// Fallback languageID
	DefaultLanguageID int

	// Fallback language locale
	DefaultLanguageShortcode string
}

func NewRND(r RND) (*RND, error) {
	if r.DB == nil {
		return nil, scill_errors.EmptyDBPointer
	}

	if r.DefaultLanguageID == 0 {
		return nil, scill_errors.EmptyLanguageID
	}

	if len(r.DefaultLanguageShortcode) == 0 {
		r.DefaultLanguageShortcode = "en"
	} else {
		r.DefaultLanguageShortcode = strings.ToLower(r.DefaultLanguageShortcode)
	}

	return &r, nil
}

func (r *RND) getGenericErr(languageID int) string {
	errorName := ""
	query := `SELECT error_message FROM error_messages WHERE error_key = ? AND language_id = ?`

	if languageID == 0 {
		languageID = r.DefaultLanguageID
	}

	r.DB.Debug().Raw(query, scill_errors.GenericErr.Error(), languageID).Row().Scan(&errorName)
	return errorName
}

func (r *RND) getErrorName(err error, languageID int) (string, error) {
	errorName := ""
	query := `SELECT error_message FROM error_messages WHERE error_key = ? AND language_id = ?`

	if languageID == 0 {
		languageID = r.DefaultLanguageID
	}

	if err := r.DB.Debug().Raw(query, err.Error(), languageID).Row().Scan(&errorName); err != nil {
		return "", scill_errors.GenericErr
	}

	return errorName, nil
}

func (r *RND) HttpErrorWithSlug(err error, languageID int, ctx *gin.Context) {
	errName, gotError := r.getErrorName(err, languageID)
	statusCode := http.StatusBadRequest

	if gotError != nil && errors.Is(gotError, scill_errors.GenericErr) {
		errName = r.getGenericErr(languageID)
		err = scill_errors.GenericErr
		statusCode = http.StatusInternalServerError
	} else if errors.Is(gotError, scill_errors.RecordNotFound){
		statusCode = http.StatusNotFound
	}

	e := HttpErrorWithErrorSlug{
		Error:      errName,
		ErrorSlug:  err.Error(),
		StatusCode: statusCode,
	}

	ctx.AbortWithStatusJSON(statusCode, e)
}

func (r *RND) ThrowStatusUnauthorized(err error, languageID int, c *gin.Context) {
	errName, gotError := r.getErrorName(err, languageID)
	if gotError != nil && errors.Is(gotError, scill_errors.GenericErr) {
		errName = r.getGenericErr(languageID)
		err = scill_errors.GenericErr
	}

	e := HttpErrorWithErrorSlug{
		Error:      errName,
		ErrorSlug:  err.Error(),
		StatusCode: http.StatusUnauthorized,
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, e)
}

func (r *RND) ThrowStatusCreated(i interface{}, c *gin.Context) {
	if i != nil {
		c.JSON(http.StatusCreated, i)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "OK",
	})
}

func (r *RND) ThrowStatusOK(i interface{}, c *gin.Context) {
	if i != nil {
		c.JSON(http.StatusOK, i)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "OK",
	})
}
