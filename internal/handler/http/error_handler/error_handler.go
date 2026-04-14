package error_handler

import (
	"errors"
	"net/http"

	"github.com/Vikktttoriya/flight-tracker/internal/handler/dto"
	"github.com/Vikktttoriya/flight-tracker/internal/service/service_errors"
)

func HandleServiceError(w http.ResponseWriter, err error) {
	var svcErr *service_errors.Error
	ok := errors.As(err, &svcErr)
	if !ok {
		dto.RespondError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}

	switch svcErr.Code {
	case service_errors.CodeInvalidArgument:
		dto.RespondError(w, http.StatusBadRequest, "bad_request", svcErr.Message)
	case service_errors.CodeNotFound:
		dto.RespondError(w, http.StatusNotFound, "not_found", svcErr.Message)
	case service_errors.CodeAlreadyExists:
		dto.RespondError(w, http.StatusConflict, "conflict", svcErr.Message)
	case service_errors.CodeInvalidCredentials:
		dto.RespondError(w, http.StatusUnauthorized, "invalid_credentials", svcErr.Message)
	case service_errors.CodeForbidden:
		dto.RespondError(w, http.StatusForbidden, "forbidden", svcErr.Message)
	case service_errors.CodeSelfModification:
		dto.RespondError(w, http.StatusForbidden, "forbidden", svcErr.Message)
	case service_errors.CodeInvalidTransition:
		dto.RespondError(w, http.StatusConflict, "invalid_transition", svcErr.Message)
	default:
		dto.RespondError(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}
