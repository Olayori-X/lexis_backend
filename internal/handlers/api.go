package handlers

import (
	auth "github.com/Olayori-X/notes/internal/handlers/auth"
	statement "github.com/Olayori-X/notes/internal/handlers/statement"
	middleware "github.com/Olayori-X/notes/internal/middleware"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func Handler(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)

	r.Route("/auth", func(router chi.Router) {
		router.Post("/login", auth.LoginHandler)
		router.Post("/signup", auth.SignupHandler)
		router.Post("/verify-otp", auth.VerifyOtpHandler)
		router.Post("/forgotpassword", auth.ForgotPasswordHandler)
		router.Post("/changepassword", auth.ChangePasswordHandler)
	})

	// r.Route("/user", func(router chi.Router) {
	// 	router.Post("/getuserprofile", user.GetProfileHandler)
	// 	router.Post("/updateuserprofile", user.UpdateUserProfileHandler)
	// 	router.Get("/users", GetUsersHandler)
	// 	router.Get("/search", user.SearchUsersHandler)
	// })

	r.Route("/statements", func(router chi.Router) {

		// Middle ware for /account authorization
		router.Use(middleware.Authorization)

		router.Post("/addstatement", statement.AddStatementHandler)
		router.Get("/", statement.GetStatementsHandler)
		router.Get("/findstatements", statement.SearchStatementsHandler)
		router.Put("/editstatement/{id}", statement.UpdateStatementHandler)
		router.Delete("/deletestatement/{id}", statement.DeleteStatementHandler)
	})

}
