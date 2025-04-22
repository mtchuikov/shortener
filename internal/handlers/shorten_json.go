package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mtchuikov/shortener/internal/models"
)

func (h *shorten) HandleJSON(rw http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(rw, req.Body, h.maxURLSizeJSON)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "payload too large", http.StatusBadRequest)
		return
	}

	var reqData models.ShortURLRequest
	err = json.Unmarshal(payload, &reqData)
	if err != nil {
		http.Error(rw, "invalid payload", http.StatusBadRequest)
		return
	}

	originalURL, err := models.NewOriginalURL(reqData.URL)
	if err != nil {
		code, msg, err := modelErrorToCodeAndMsg(err)
		if err != nil {
			logUnexpectedError(h.log, err, resolveOp)
		}

		http.Error(rw, msg, code)
		return
	}

	statusCode := http.StatusCreated

	shortenURL, err := h.shortener.Serve(req.Context(), originalURL)
	if err != nil {
		code, msg, err := serviceErrorToCodeAndMsg(err)
		if err != nil {
			logUnexpectedError(h.log, err, resolveOp)
		}

		if code != http.StatusConflict {
			http.Error(rw, msg, code)
			return
		}

		statusCode = http.StatusConflict
	}

	respData := models.ShortURLResponse{Result: shortenURL.String()}
	payload, err = json.Marshal(respData)
	if err != nil {
		http.Error(rw, "failed to marshal", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(payload)
}
