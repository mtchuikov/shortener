package middlewares

import "net/http"

func OnlyMethod(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if method != req.Method {
			errMsg := "method not allowed"
			http.Error(rw, errMsg, http.StatusMethodNotAllowed)
			return
		}

		next(rw, req)
	}
}
