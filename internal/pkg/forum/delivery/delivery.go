package delivery

import (
	"TPForum/internal/pkg/domain"
	"TPForum/internal/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ForumDelivery struct {
	FUcase domain.ForumUsecase
}

func NewForumDelivery(router *mux.Router, usecase domain.ForumUsecase) {
	forumDelivery := &ForumDelivery{
		FUcase: usecase,
	}

	router.HandleFunc("/forum/create", forumDelivery.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/forum/{slug}/details", forumDelivery.GetForumDetails).Methods("GET", "OPTIONS")
	router.HandleFunc("/forum/{slug}/users", forumDelivery.GetForumUsers).Methods("GET", "OPTIONS")
}

func (f *ForumDelivery) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	forum := &domain.Forum{}
	if err := decoder.Decode(forum); err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	_, err := f.FUcase.Create(forum)
	switch err {
	case domain.NotFoundError:
		utils.WriteResponseError(w, http.StatusNotFound,
			fmt.Sprintf("Can't find user with nickname: %s", forum.User))
		return
	case domain.ConflictError:
		exists, err := f.FUcase.GetForumDetails(forum.Slug)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		data, err := json.Marshal(*exists)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.WriteResponse(w, http.StatusConflict, data)
	case nil:
		data, err := json.Marshal(*forum)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.WriteResponse(w, http.StatusCreated, data)
	default:
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

}

func (f *ForumDelivery) GetForumDetails(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	forum := &domain.Forum{}
	var err error

	forum, err = f.FUcase.GetForumDetails(slug)

	if err != nil {
		utils.WriteResponseError(w, http.StatusNotFound, fmt.Sprintf("Can't find forum with slug: %s", slug))
		return
	}

	data, err := json.Marshal(*forum)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}

func (f *ForumDelivery) GetForumUsers(w http.ResponseWriter, r *http.Request)  {
	slug := mux.Vars(r)["slug"]

	var desc bool
	var limit int
	var err error

	limitString := r.URL.Query().Get("limit")
	since := r.URL.Query().Get("since")
	descString := r.URL.Query().Get("desc")

	if limitString == "" {
		limit = 100
	} else {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if descString == "" || descString == "false" {
		desc = false
	} else {
		desc = true
	}

	users, err := f.FUcase.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		switch err {
		case domain.NotFoundError:
			utils.WriteResponseError(w, http.StatusNotFound,"Can't find forum with slug: " + slug)
			return
		default:
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	data, err := json.Marshal(users)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}
