package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type PaginateCtx struct{}

type Pagination struct {
	Start        int64                  `json:"start"`
	Size         int64                  `json:"size"`
	Filters      MRT_ColumnFiltersState `json:"filters"`
	GlobalFilter string                 `json:"globalFilter"`
	Sorting      MRT_SortingState       `json:"sorting"`
}

type MRT_ColumnFiltersState []MRT_ColumnFilter

type MRT_ColumnFilter struct {
	Id    string `json:"id"`
	Value string `json:"value"` // Le type de Value peut varier
}

type MRT_SortingState []MRT_Sorting

type MRT_Sorting struct {
	Id   string `json:"id"`
	Desc bool   `json:"desc"`
}

type Meta struct {
	TotalRowCount int
}
type ApiResponse[T any] struct {
	Data []T
	Meta Meta
}

func (rd *ApiResponse[T]) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Récupérer le paramètre 'pagination' de la requête GET
		paginationJSON := r.URL.Query().Get("pagination")

		var pagination Pagination

		// Désérialiser la chaîne JSON en structure Go
		err := json.Unmarshal([]byte(paginationJSON), &pagination)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erreur lors de la désérialisation de 'pagination': %v", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), PaginateCtx{}, &pagination)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
