package delivery

import (
	"TPForum/internal/pkg/domain"
	"TPForum/internal/pkg/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)


type ThreadDelivery struct {
	TUcase domain.ThreadUsecase
}

func NewThreadDelivery(router *mux.Router, usecase domain.ThreadUsecase) {
	threadDelivery := &ThreadDelivery{
		TUcase: usecase,
	}

	router.HandleFunc("/forum/{slug}/create", threadDelivery.Create).Methods("POST", "OPTIONS")
	router.HandleFunc("/thread/{slug_or_id}/vote", threadDelivery.Vote).Methods("POST", "OPTIONS")
	router.HandleFunc("/forum/{slug}/threads", threadDelivery.GetForumThreads).Methods("GET", "OPTIONS")
	router.HandleFunc("/thread/{slug_or_id}/details", threadDelivery.GetThreadBySlugOrId).Methods("GET", "OPTIONS")
	router.HandleFunc("/thread/{slug_or_id}/details", threadDelivery.UpdateThreadBySlugOrId).Methods("POST", "OPTIONS")
}

func (d *ThreadDelivery) Create(w http.ResponseWriter, r *http.Request)  {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	thread := &domain.Thread{}
	if err := decoder.Decode(thread); err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	thread.Forum = mux.Vars(r)["slug"]

	err := d.TUcase.Create(thread)
	switch err {
	case domain.NotFoundError:
		utils.WriteResponseError(w, http.StatusNotFound,"Can't find user or forum")
		return
	case domain.ConflictError:
		exists, err := d.TUcase.GetThreadBySlug(thread.Slug)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		data, err := json.Marshal(exists)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.WriteResponse(w, http.StatusConflict, data)
		return
	case nil:
		data, err := json.Marshal(*thread)
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.WriteResponse(w, http.StatusCreated, data)
		return
	default:
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (d *ThreadDelivery) GetForumThreads(w http.ResponseWriter, r *http.Request) {
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

	threads, err := d.TUcase.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		switch err {
		case domain.NotFoundError:
			utils.WriteResponseError(w, http.StatusNotFound, "Can't find forum with slug: " + slug)
			return
		default:
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	data, err := json.Marshal(threads)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}

func (d *ThreadDelivery) GetThreadBySlugOrId(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]

	thread, err := d.TUcase.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		utils.WriteResponseError(w, http.StatusNotFound, "Can't find thread: " + slugOrId)
		return
	}

	data, err := json.Marshal(thread)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}

func (d *ThreadDelivery) UpdateThreadBySlugOrId(w http.ResponseWriter, r *http.Request) {
	slugOrId := mux.Vars(r)["slug_or_id"]

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	thread := &domain.Thread{}
	if err := decoder.Decode(thread); err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	err := d.TUcase.UpdateThreadBySlugOrId(slugOrId, thread)

	if err != nil {
		utils.WriteResponseError(w, http.StatusNotFound, "Can't find thread: " + slugOrId)
		return
	}

	data, err := json.Marshal(*thread)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}

func (d *ThreadDelivery) Vote(w http.ResponseWriter, r *http.Request)  {
	slugOrId := mux.Vars(r)["slug_or_id"]

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	vote := &domain.Vote{}
	if err := decoder.Decode(vote); err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
	}

	thread, err := d.TUcase.Vote(slugOrId, *vote)

	if err != nil {
		utils.WriteResponseError(w, http.StatusNotFound, "Can't find thread: " + slugOrId)
		return
	}

	data, err := json.Marshal(thread)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteResponse(w, http.StatusOK, data)
}
