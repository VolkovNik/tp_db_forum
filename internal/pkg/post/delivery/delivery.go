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

type PostDelivery struct {
	PUcase domain.PostUsecase
}

func NewPostDelivery(router *mux.Router, usecase domain.PostUsecase) {
	postDelivery := &PostDelivery{
		PUcase: usecase,
	}

	router.HandleFunc("/thread/{slug_or_id}/create", postDelivery.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/thread/{slug_or_id}/posts", postDelivery.GetThreadPosts).Methods("GET", "OPTIONS")
	router.HandleFunc("/post/{id}/details", postDelivery.UpdateById).Methods("POST", "OPTIONS")
	router.HandleFunc("/post/{id}/details", postDelivery.Details).Methods("GET", "OPTIONS")
}

func (p *PostDelivery) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	slugOrId := mux.Vars(r)["slug_or_id"]

	decoder := json.NewDecoder(r.Body)
	posts := &domain.Posts{}
	if err := decoder.Decode(posts); err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	err := p.PUcase.Create(slugOrId, posts)

	switch err {
	case domain.NotFoundError:
		utils.WriteResponseError(w, http.StatusNotFound,
			fmt.Sprintf("Can't find thread: %s", slugOrId))
		return
	case domain.ParentError:
		utils.WriteResponseError(w, http.StatusConflict, "Parent post was created in another thread")
	case nil:
		data, err := json.Marshal(*posts)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.WriteResponse(w, http.StatusCreated, data)
	default:
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}
}

func (p *PostDelivery) UpdateById(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	post := &domain.Post{}
	if err = decoder.Decode(post); err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	err = p.PUcase.UpdateById(id, post)

	if err != nil {
		utils.WriteResponseError(w, http.StatusNotFound, "Can't find post with id: " + idString)
		return
	}

	data, err := json.Marshal(*post)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}

func (p *PostDelivery) Details(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	related :=  r.URL.Query().Get("related")

	postFull, err := p.PUcase.Details(id, related)

	if err != nil {
		utils.WriteResponseError(w, http.StatusNotFound, "Can't find post with id: " + idString)
		return
	}

	data, err := json.Marshal(postFull)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}

func (p *PostDelivery) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]

	var desc bool
	var limit int
	var since int
	var err error

	limitString := r.URL.Query().Get("limit")
	sinceString := r.URL.Query().Get("since")
	sort 		:= r.URL.Query().Get("sort")
	descString 	:= r.URL.Query().Get("desc")

	if limitString == "" {
		limit = 100
	} else {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if sinceString == "" {
		since = 0
	} else {
		since, err = strconv.Atoi(sinceString)
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

	posts, err := p.PUcase.GetThreadPosts(slugOrId, limit, since, sort, desc)

	if err != nil {
		if err == domain.NotFoundError {
			utils.WriteResponseError(w, http.StatusNotFound, "Can't find thread with : " + slugOrId)
			return
		}
	}

	data, err := json.Marshal(posts)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}
