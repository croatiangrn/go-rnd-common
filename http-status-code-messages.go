package go_rnd_common

import (
	"errors"
	"fmt"
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

func (r *RND) GetErrorName(err error, languageID int) (string, error) {
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

func (r *RND) GetErrorfName(err error, languageID int, values ...interface{}) (string, error) {
	errorName := ""
	query := `SELECT error_message FROM error_messages WHERE error_key = ? AND language_id = ?`

	if languageID == 0 {
		languageID = r.DefaultLanguageID
	}

	if err := r.DB.Debug().Raw(query, err.Error(), languageID).Row().Scan(&errorName); err != nil {
		return "", scill_errors.GenericErr
	}

	return fmt.Sprintf(errorName, values...), nil
}

func (r *RND) HttpErrorWithSlug(err error, languageID int, ctx *gin.Context) {
	errName, gotError := r.GetErrorName(err, languageID)
	statusCode := http.StatusBadRequest

	if gotError != nil && errors.Is(gotError, scill_errors.GenericErr) {
		errName = r.getGenericErr(languageID)
		err = scill_errors.GenericErr
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, scill_errors.RecordNotFound) {
		statusCode = http.StatusNotFound
	} else if errors.Is(err, scill_errors.GenericErr) {
		statusCode = http.StatusInternalServerError
	}

	e := HttpErrorWithErrorSlug{
		Error:      errName,
		ErrorSlug:  err.Error(),
		StatusCode: statusCode,
	}
	ctx.AbortWithStatusJSON(statusCode, e)
}

func (r *RND) HttpErrorfWithSlug(err error, languageID int, ctx *gin.Context, values ...interface{}) {
	errName, gotError := r.GetErrorfName(err, languageID, values...)
	statusCode := http.StatusBadRequest

	if gotError != nil && errors.Is(gotError, scill_errors.GenericErr) {
		errName = r.getGenericErr(languageID)
		err = scill_errors.GenericErr
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, scill_errors.RecordNotFound) {
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
	errName, gotError := r.GetErrorName(err, languageID)
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

type SCILLServiceResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
	ErrorSlug  string `json:"error_slug,omitempty"`
}

func (s *SCILLServiceResponse) formatOK(message string) {
	s.StatusCode = http.StatusOK
	s.Message = message
}

func (s *SCILLServiceResponse) formatError(error, errorSlug string, statusCode int) {
	s.Error = error
	s.ErrorSlug = errorSlug
	s.StatusCode = statusCode
}

func (s *SCILLServiceResponse) ThrowStatusOK(message string, c *gin.Context) {
	s.formatOK(message)
	c.JSON(s.StatusCode, s)
}

func (s *SCILLServiceResponse) HttpErrorWithSlug(r *RND, err error, languageID int, ctx *gin.Context) {
	errName, gotError := r.GetErrorName(err, languageID)
	statusCode := http.StatusBadRequest

	if gotError != nil && errors.Is(gotError, scill_errors.GenericErr) {
		errName = r.getGenericErr(languageID)
		err = scill_errors.GenericErr
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, scill_errors.RecordNotFound) {
		statusCode = http.StatusNotFound
	} else if errors.Is(err, scill_errors.GenericErr) {
		statusCode = http.StatusInternalServerError
	}

	s.formatError(errName, err.Error(), statusCode)
	ctx.AbortWithStatusJSON(statusCode, s)
}

func (s *SCILLServiceResponse) HttpErrorfWithSlug(r *RND, err error, languageID int, ctx *gin.Context, values ...interface{}) {
	errName, gotError := r.GetErrorfName(err, languageID, values...)
	statusCode := http.StatusBadRequest

	if gotError != nil && errors.Is(gotError, scill_errors.GenericErr) {
		errName = r.getGenericErr(languageID)
		err = scill_errors.GenericErr
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, scill_errors.RecordNotFound) {
		statusCode = http.StatusNotFound
	} else if errors.Is(err, scill_errors.Unauthorized) {
		statusCode = http.StatusUnauthorized
	}

	s.formatError(errName, err.Error(), statusCode)
	ctx.AbortWithStatusJSON(statusCode, s)
}