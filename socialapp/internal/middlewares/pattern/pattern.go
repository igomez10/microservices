package pattern

import (
	"net/http"
	"socialapp/internal/contexthelper"
)

type Pattern struct {
	Pattern string
}

// Adds the request pattern into the context, this is only required because chi.mux
// does not provide a way to get the pattern from the request. This middleware will update the
// string pointer saved in the context as "pattern"
func (p *Pattern) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(contexthelper.SetRequestPatternInContext(r.Context(), p.Pattern))
		next.ServeHTTP(w, r)
	})
}
