package handlers

import (
	"errors"
	"net/http"
	"organization_structure/internal/model/out"
	"organization_structure/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func HDeleteDepartment(s *usecase.UseCase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var id = chi.URLParam(r, "id")
		var mode = r.URL.Query().Get("mode")
		var reassign_to_department_id = r.URL.Query().Get("reassign_to_department_id")

		err := s.DeleteDepartment(r.Context(), id, mode, new(reassign_to_department_id))

		if err != nil {
			sErr, ok := errors.AsType[*out.CustomError](err)

			if ok {
				sendError(w, sErr.Code, sErr.Error(), err)
				return
			}

			sendError(w, 500, err.Error(), nil)
			return
		}

		sendSuccess(w, 204, "Department is deleted", nil)
	}
}
