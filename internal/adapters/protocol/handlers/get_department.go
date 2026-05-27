package handlers

import (
	"errors"
	"net/http"
	"organization_structure/internal/model/out"
	"organization_structure/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HGetDepartment(s *usecase.UseCase) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var depId = chi.URLParam(r, "id")
		var depth = r.URL.Query().Get("depth")
		var includeEmployees = r.URL.Query().Get("include_employees")

		_, err := strconv.Atoi(depId)

		if err != nil {
			sendError(w, 400, err.Error(), nil)
			return
		}

		dep, err := s.GetDepartment(r.Context(), depId, depth, includeEmployees)

		if err != nil {
			sErr, ok := errors.AsType[*out.CustomError](err)

			if ok {
				sendError(w, sErr.Code, sErr.Error(), err)
				return
			}

			sendError(w, 500, err.Error(), nil)
			return
		}

		sendSuccess(w, 200, "Department info is obtained", dep)
	}
}
