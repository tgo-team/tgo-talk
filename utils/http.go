package utils

import "net/http"

type ResultError struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

func NewResultError(errCode int, errMsg string) *ResultError {

	resultError := &ResultError{}
	resultError.ErrCode = errCode
	resultError.ErrMsg = errMsg

	return resultError
}

func ResponseError400(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ResponseError(w, http.StatusBadRequest, msg)
}


func ResponseError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := ResultError{statusCode, msg}
	w.WriteHeader(statusCode)
	WriteJson(w, err)
}

func ResponseSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := NewResultError(0, "OK")
	WriteJson(w, err)

}

func ResponseSuccessJson(w http.ResponseWriter,obj interface{})  {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	WriteJson(w, obj)
}

func BindJson(r *http.Request,obj interface{}) error  {
	return ReadJson(r.Body,obj)
}