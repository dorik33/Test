package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dorik33/Test/internal/client"
	"github.com/dorik33/Test/internal/models"
	"github.com/dorik33/Test/internal/store"
	"github.com/gorilla/mux"
)

// Message represents error response
// @Description Error response structure
type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

// @Summary Get songs list
// @Description Get paginated list of songs with filters
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group query string false "Group name filter"
// @Param song query string false "Song name filter"
// @Param limit query int false "Results limit (default 10)"
// @Param offset query int false "Results offset"
// @Success 200 {array} models.Song
// @Failure 500 {object} Message
// @Router /songs [get]
func (api *API) GetSongsHandler(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	songs, err := api.store.SongRepository.GetSongs(r.Context(), group, song, limit, offset)
	if err != nil {
		api.logger.Infof("Get songs failed: %v\n", err)
		api.sendErrorResponse(w, http.StatusInternalServerError, "Error with database. Please Try later")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songs)
}

// @Summary Get song lyrics
// @Description Get song lyrics with pagination
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param limit query int false "Lines limit"
// @Param offset query int false "Lines offset"
// @Success 200 {object} models.SongText
// @Failure 400 {object} Message
// @Failure 404 {object} Message
// @Failure 500 {object} Message
// @Router /songText/{id} [get]
func (api *API) GetTextSongByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.logger.Infof("Invalid ID provided in request: %v", err)
		api.sendErrorResponse(w, http.StatusBadRequest, "Use integer for the id parameter")
		return
	}

	text, err := api.store.SongRepository.GetSongTextByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		if errors.Is(err, store.ErrSongNotFound) {
			api.logger.Infof("Song not found. ID: %d", id)
			api.sendErrorResponse(w, http.StatusNotFound, "Song not found")
			return
		}

		api.logger.Infof("Database error GetSongTextByID. ID: %d, Error: %v", id, err)
		api.sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	s := strings.Split(text, "\n")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 || limit > len(s) {
		limit = len(s)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 || offset > len(s) {
		offset = 0
	}

	msg := models.SongText{
		Text: strings.Join(s[offset:limit], " "),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

// @Summary Delete song
// @Description Delete song by ID
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Success 204
// @Failure 400 {object} Message
// @Failure 404 {object} Message
// @Failure 500 {object} Message
// @Router /song/{id} [delete]
func (api *API) DeleteSonghandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.logger.Infof("Invalid ID provided in request: %v", err)
		api.sendErrorResponse(w, http.StatusBadRequest, "Use integer for the id parameter")
		return
	}
	err = api.store.SongRepository.DeleteSong(r.Context(), id)
	if err != nil {
		api.logger.Infof("Error with delete song: %v\n", err)
		if errors.Is(err, store.ErrSongNotFound) {
			api.sendErrorResponse(w, http.StatusNotFound, "Song not found")
			return
		}
		api.logger.Infof("Database error: %v", err)
		api.sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Update song
// @Description Update existing song
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Song data"
// @Success 204
// @Failure 400 {object} Message
// @Failure 404 {object} Message
// @Failure 500 {object} Message
// @Router /song/{id} [put]
func (api *API) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		api.logger.Infof("Invalid ID provided in request: %v", err)
		api.sendErrorResponse(w, http.StatusBadRequest, "Use integer for the id parameter")
		return
	}

	var song models.Song
	err = json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		api.logger.Infof("Error with decode json: %v", err)
		msg := Message{
			StatusCode: http.StatusBadRequest,
			Message:    "Song is not valid",
			IsError:    true,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}

	err = api.store.SongRepository.UpdateSong(r.Context(), id, song)
	if err != nil {
		api.logger.Infof("Error with update song: %v\n", err)
		if errors.Is(err, store.ErrSongNotFound) {
			api.sendErrorResponse(w, http.StatusNotFound, "Song not found")
			return
		}
		api.logger.Infof("Database error: %v", err)
		api.sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AddSongHandler godoc
// @Summary Add new song
// @Description Add new song to database
// @Tags songs
// @Accept  json
// @Produce  json
// @Param request body models.SongRequest true "Song request"
// @Success 204
// @Failure 400 {object} Message
// @Failure 500 {object} Message
// @Router /song [post]
func (api *API) AddSongHandler(w http.ResponseWriter, r *http.Request) {
	var songReq models.SongRequest
	if err := json.NewDecoder(r.Body).Decode(&songReq); err != nil {
		api.logger.Infof("Error with decode song request %v\n", err)
		api.sendErrorResponse(w, http.StatusBadRequest, "Request is not valid")
		return
	}
	api.logger.Debugf("Calling external API: %s, Response: %v", api.config.ApiBaseURL, songReq)
	songDetail, err := client.FetchSongInfo(api.config.ApiBaseURL, songReq)
	if err != nil {
		api.logger.Infof("Error with api: %v\n", err)
		api.sendErrorResponse(w, http.StatusInternalServerError, "Internal error. Try again later")
		return
	}

	song := models.Song{
		GroupName:   songReq.GroupName,
		SongName:    songReq.SongName,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	err = api.store.SongRepository.AddSong(r.Context(), song)
	if err != nil {
		api.logger.Infof("Error with add song: %v\n", err)
		api.sendErrorResponse(w, http.StatusInternalServerError, "Error with database. Please Try later")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
