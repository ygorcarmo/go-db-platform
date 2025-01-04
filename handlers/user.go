package handlers

import (
	"log"
	"net/http"

	"db-platform/models"
	"db-platform/storage"
	"db-platform/views/components"
	"db-platform/views/setting"
	appUser "db-platform/views/user"

	"github.com/go-chi/chi/v5"
)

func GetResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	err := appUser.ResetPassword().Render(r.Context(), w)
	if err != nil {
		log.Fatal("Error when trying to render reset password page")
	}
}

func ResetApplicationUserPasswordHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	user := r.Context().Value(models.UserCtx).(*models.AppUser)

	p := r.FormValue("password")
	// TODO: add validation to new password
	np := r.FormValue("new-password")
	np2 := r.FormValue("re-password")

	if np != np2 {
		components.Response(models.CreateResponse("The new passwords don't match", false)).Render(r.Context(), w)
		return
	}

	match, err := hashParams.ComparePasswordAndHash(p, user.Password)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	if !match {
		components.Response(models.CreateResponse("Invalid Current Password.", false)).Render(r.Context(), w)
		return
	}

	hashedP, err := hashParams.HashPasword(np)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	err = s.UpdateApplicationUserPassword(user.Id, hashedP)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	components.Response(models.CreateResponse("Password updated successfully", true)).Render(r.Context(), w)

}

func GetAllUserSettingsPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	users, err := s.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	setting.GetUsersPage(users).Render(r.Context(), w)
}

func GetCreateUserPage(w http.ResponseWriter, r *http.Request) {
	err := appUser.CreateUserPage().Render(r.Context(), w)
	if err != nil {
		log.Fatal("error when trying to render create user page")
	}
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	cUser := r.Context().Value(models.UserCtx).(*models.AppUser)

	u := r.FormValue("username")
	p := r.FormValue("password")
	p2 := r.FormValue("re-password")
	sup := r.FormValue("supervisor")
	sec := r.FormValue("sector")
	isA := r.FormValue("admin")

	if p != p2 {
		components.Response(models.CreateResponse("Passwords do not match", false)).Render(r.Context(), w)
		return
	}

	isAdmin := false
	if isA != "" {
		isAdmin = true
	}

	encodedHash, err := hashParams.HashPasword(p)
	if err != nil {
		components.Response(models.CreateResponse("Something went wrong when hashing the password.", false)).Render(r.Context(), w)
		return
	}

	user := models.AppUser{Username: u, Password: encodedHash, Supervisor: sup, Sector: sec, IsAdmin: isAdmin}

	id, err := s.CreateApplicationUser(user)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	go s.CreateAdminLog(models.AdminLog{UserId: cUser.Id, Action: models.CreateAdminAction, ResourceId: id, ResourceType: models.User, ResourceName: u})

	w.Header().Add("HX-Redirect", "/settings/users")
	w.Write([]byte("user created"))

}

func DeleteUserById(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	cUser := r.Context().Value(models.UserCtx).(*models.AppUser)

	id := chi.URLParam(r, "id")
	err := s.DeleteUserById(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	go func() {
		user, err := s.GetUserById(id)
		if err != nil {
			return
		}
		s.CreateAdminLog(models.AdminLog{UserId: cUser.Id, Action: models.DeleteAdminAction, ResourceId: id, ResourceType: models.User, ResourceName: user.Username})
	}()
	w.WriteHeader(http.StatusOK)
}

func GetEditUserSettingsPage(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	id := chi.URLParam(r, "id")

	user, err := s.GetUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}

	err = appUser.EditUserPage(user).Render(r.Context(), w)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}

}

func UpdateUserSettingsHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	cUser := r.Context().Value(models.UserCtx).(*models.AppUser)

	id := chi.URLParam(r, "id")
	u := r.FormValue("username")
	sup := r.FormValue("supervisor")
	d := r.FormValue("sector")
	a := r.FormValue("admin")

	isAdmin := false
	if a != "" {
		isAdmin = true
	}

	user := models.AppUser{
		Username:   u,
		Supervisor: sup,
		Sector:     d,
		IsAdmin:    isAdmin,
		Id:         id,
	}

	err := s.UpdateApplicationUser(user)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	go s.CreateAdminLog(models.AdminLog{UserId: cUser.Id, Action: models.UpdateSettingsAdminAction, ResourceId: id, ResourceType: models.User, ResourceName: u})

	w.Header().Add("HX-Redirect", "/settings/users")
}

func GetUpdateAppUserCredentials(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	id := chi.URLParam(r, "id")
	name := r.URL.Query().Get("username")

	err := appUser.UpdateAppUserCredentials(id, name).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong when rendering update user app credentials"))
		return
	}
}

func UpdateAppUserCredentialsHandler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	cUser := r.Context().Value(models.UserCtx).(*models.AppUser)

	// TODO: add server side validation
	id := chi.URLParam(r, "id")
	u := r.FormValue("username")
	p := r.FormValue("new-password")
	p2 := r.FormValue("re-password")

	if p != p2 {
		components.Response(models.CreateResponse("Passwords do not match", false)).Render(r.Context(), w)
		return
	}

	h, err := hashParams.HashPasword(p)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	err = s.UpdateApplicationUserCredentials(u, h, id)
	if err != nil {
		components.Response(models.CreateResponse(err.Error(), false)).Render(r.Context(), w)
		return
	}

	w.Header().Add("HX-Redirect", "/settings/users")

	go s.CreateAdminLog(models.AdminLog{UserId: cUser.Id, Action: models.UpdateCredentialsAdminAction, ResourceId: id, ResourceType: models.User, ResourceName: u})

}

func UnlockUser(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid user id"))
		return
	}

	err := s.ResetUserLoginAttempts(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Redirect", "/settings/users")
}
