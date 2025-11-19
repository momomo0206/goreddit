package web

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/momomo0206/goreddit"
)

func NewHandler(store goreddit.Store, sessions *scs.SessionManager, csrfkey []byte) *Handler {
	h := &Handler{
		Mux:      chi.NewMux(),
		store:    store,
		sessions: sessions,
	}

	threads := ThreadHandler{store: store, sessions: sessions}
	posts := PostHandler{store: store, sessions: sessions}
	comments := CommentHandler{store: store, sessions: sessions}

	h.Use(middleware.Logger)
	h.Use(csrf.Protect(
		csrfkey,
		csrf.Secure(false),
		csrf.TrustedOrigins([]string{"localhost:3000"}),
	))
	h.Use(sessions.LoadAndSave)

	// h.Use(csrf.Protect(
	// 	csrfkey,
	// 	csrf.Secure(false),
	// 	csrf.SameSite(csrf.SameSiteLaxMode),
	// 	csrf.Path("/"),
	// 	csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		reason := csrf.FailureReason(r)
	// 		log.Printf("=== CSRF Verification Failed ===")
	// 		log.Printf("Reason: %v", reason)
	// 		log.Printf("Method: %s, URL: %s", r.Method, r.URL.Path)
	// 		log.Printf("Origin: %s", r.Header.Get("Origin"))
	// 		log.Printf("Referer: %s", r.Header.Get("Referer"))

	// 		// Cookieの確認
	// 		cookie, err := r.Cookie("_gorilla_csrf")
	// 		if err != nil {
	// 			log.Printf("Cookie Error: %v", err)
	// 		} else {
	// 			log.Printf("Cookie Value: %s", cookie.Value)
	// 		}

	// 		// POSTデータの確認
	// 		r.ParseForm()
	// 		log.Printf("Form Token: %s", r.FormValue("gorilla.csrf.Token"))
	// 		log.Printf("================================")

	// 		http.Error(w, fmt.Sprintf("CSRF verification failed: %v", reason), http.StatusForbidden)
	// 	})),
	// ))

	h.Get("/", h.Home())
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", threads.List())
		r.Get("/new", threads.Create())
		r.Post("/", threads.Store())
		r.Get("/{id}", threads.Show())
		r.Post("/{id}/delete", threads.Delete())
		r.Get("/{id}/new", posts.Create())
		r.Post("/{id}", posts.Store())
		r.Get("/{threadID}/{postID}", posts.Show())
		r.Get("/{threadID}/{postID}/vote", posts.Vote())
		r.Post("/{threadID}/{postID}", comments.Store())
	})
	h.Get("/comments/{id}/vote", comments.vote())

	return h
}

type Handler struct {
	*chi.Mux

	store    goreddit.Store
	sessions *scs.SessionManager
}

func (h *Handler) Home() http.HandlerFunc {
	type data struct {
		Posts []goreddit.Post
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		pp, err := h.store.Posts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{Posts: pp})
	}
}
