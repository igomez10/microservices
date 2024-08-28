package pattern

import (
	"fmt"
	"net/http"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
)

type Pattern struct {
	Pattern string
}

// Adds the request pattern into the context, this is only required because chi.mux
// does not provide a way to get the pattern from the request. This middleware will update the
// string pointer saved in the context as "pattern"
func (p *Pattern) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracerhelper.GetTracer().Start(r.Context(), fmt.Sprintf(r.Method+"_%s", r.URL.Path))
		defer span.End()

		r = r.WithContext(ctx)

		r = r.WithContext(contexthelper.SetRequestPatternInContext(r.Context(), p.Pattern))
		next.ServeHTTP(w, r)
	})
}
