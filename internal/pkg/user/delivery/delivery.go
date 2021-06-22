package delivery

import (
	"TPForum/internal/pkg/domain"
	"TPForum/internal/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type UserDelivery struct {
	UUcase domain.UserUsecase
}

func NewUserDelivery(router *mux.Router, usecase domain.UserUsecase) {
	userDelivery := &UserDelivery{
		UUcase: usecase,
	}

	router.HandleFunc("/user/{nickname}/create", userDelivery.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/{nickname}/profile", userDelivery.GetProfileInfo).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/{nickname}/profile", userDelivery.UpdateProfileInfo).Methods("POST", "OPTIONS")

}

func (u *UserDelivery) Create(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	decoder := json.NewDecoder(r.Body)
	user := &domain.User{}
	if err := decoder.Decode(user); err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	user.Nickname = nickname

	_, err := u.UUcase.Create(*user)

	if err != nil {
		users := domain.Users{}
		users, err := u.UUcase.SelectByEmailOrNickname(user.Nickname, user.Email)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		}
		data, err := json.Marshal(users)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.WriteResponse(w, http.StatusConflict, data)
		return
	}

	data, err := json.Marshal(*user)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusCreated, data)

}

func (u UserDelivery) GetProfileInfo(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]
	user := domain.User{}

	user, err := u.UUcase.GetProfileInfo(nickname)
	if err != nil {
		utils.WriteResponseError(w, http.StatusNotFound,
			fmt.Sprintf("Can't find user with nickname: %s", nickname))
		return
	}
	data, err := json.Marshal(user)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}

func (u UserDelivery) UpdateProfileInfo(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	user := &domain.User{}
	if err := decoder.Decode(user); err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	user.Nickname = nickname

	err := u.UUcase.UpdateProfileInfo(user)
	switch err {
	case domain.NotFoundError:
		utils.WriteResponseError(w, http.StatusNotFound, fmt.Sprintf("Can't find user by nickname: %s", nickname))
		return
	case domain.ConflictError:
		utils.WriteResponseError(w, http.StatusConflict, "This email already used by another user")
	case nil:
		data, err := json.Marshal(*user)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.WriteResponse(w, http.StatusOK, data)
	default:
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}
}
