package utils

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrBadRequest(err error, data any) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
		Data:           data,
	}
}

func ErrUnauthorizedRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized request.",
		ErrorText:      err.Error(),
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
	Data       any    `json:"data,omitempty"`  // donnee complmentaire
}

func ToPgText(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}

func ToPgInt4(value int) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(value),
		Valid: true,
	}
}

func ToPgDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{
			Valid: false,
		}
	}
	return pgtype.Date{
		Time:  *t,
		Valid: true,
	}
}

// toPointer returns a pointer to the given value (works for any type).
func ToPointer[T any](value T) *T {
	return &value
}

type ContextKey struct {
	Name string
}

func (k *ContextKey) String() string {
	return "context value: " + k.Name
}

func RenderError(w http.ResponseWriter, r *http.Request, method string, err error) {
	render.Render(w, r, ErrRender(err))
	GetLogger(r.Context()).Warn(method, zap.String("", err.Error()))
}

func RenderBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	var errorsFormulaire = ValidationErrorsToMap(err)
	if errorsFormulaire != nil {
		render.Render(w, r, ErrBadRequest(errors.New("erreur de validation"), errorsFormulaire))
		return
	}

	RenderError(w, r, "UpdateFormation", err)
}

func ValidationErrorsToMap(err error) map[string]interface{} {
	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		errMap := map[string]interface{}{}
		for _, e := range validateErrs {
			AddErrorToMap(errMap, e.Namespace(), e.Tag())
		}
		return errMap
	}
	return nil
}

func RenderBadRequestMap(w http.ResponseWriter, r *http.Request, err map[string]string) {
	render.Render(w, r, ErrBadRequest(errors.New("erreur de validation"), err))

}

func AddErrorToMap(errMap map[string]any, path string, message string) {
	matches := extractTokens(path)
	field := matches[1]

	if len(matches) == 2 {
		errMap[field] = message
	} else if len(matches) == 4 {

		index, _ := strconv.Atoi(matches[2])
		subField := matches[3]
		// Création de la hiérarchie (slice de maps)
		if _, ok := errMap[field]; !ok {
			errMap[field] = make([]map[string]string, 0)
		}

		fieldMap := errMap[field].([]map[string]string)

		// Assurer que la slice a la capacité suffisante pour l'index
		for i := len(fieldMap); i <= index; i++ {
			fieldMap = append(fieldMap, nil) // Ajouter des éléments nil jusqu'à l'index souhaité
		}
		errMap[field] = fieldMap // Mettre à jour la map principale

		// Initialiser la map à l'index si elle n'existe pas
		if fieldMap[index] == nil {
			fieldMap[index] = make(map[string]string)
		}

		fieldMap[index][subField] = message
	} else {
		return
	}
}

func extractTokens(path string) []string {
	var tokens []string
	parts := strings.Split(path, ".")
	for _, part := range parts {
		if strings.Contains(part, "[") && strings.Contains(part, "]") {
			// Cas pour les tableaux (ex: Promotions[0])
			regex := regexp.MustCompile(`^([^\[]+)\[(\d+)\]$`)
			matches := regex.FindStringSubmatch(part)
			if len(matches) == 3 {
				tokens = append(tokens, matches[1], fmt.Sprintf(matches[2]))
			} else {
				tokens = append(tokens, part) // Si le format n'est pas correct, on garde la partie entière
			}
		} else {
			// Cas pour les champs simples (ex: FormationRequest, Name, Debut)
			tokens = append(tokens, part)
		}
	}
	return tokens
}
