package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"organization_structure/internal/model/in"
	"organization_structure/internal/model/out"
	"organization_structure/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func HUpdateDepartment(s *usecase.UseCase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var content in.UpdateDepartment
		var depId = chi.URLParam(r, "id")

		err := json.NewDecoder(r.Body).Decode(&content)

		if err != nil {
			sendError(w, 400, err.Error(), nil)
			return
		}

		dep, err := s.UpdateDepartment(r.Context(), depId, content)

		if err != nil {
			sErr, ok := errors.AsType[*out.CustomError](err)

			if ok {
				sendError(w, sErr.Code, sErr.Error(), err)
				return
			}

			sendError(w, 500, err.Error(), nil)
			return
		}

		sendSuccess(w, 200, "Department updated", dep)
	}
}
