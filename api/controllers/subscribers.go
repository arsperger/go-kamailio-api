package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.com/voip-services/go-kamailio-api/api/models"
	"gitlab.com/voip-services/go-kamailio-api/internal/respond"

	"github.com/gorilla/mux"
)

// GetSubscribers gets subscribers
func (a *App) getSubscribers(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	offset, _ := strconv.Atoi(params["page"])
	limit, _ := strconv.Atoi(params["rows"])

	if offset < 0 {
		offset = 0
	}

	if limit > 10 || limit < 1 {
		limit = 10
	}

	subs, err := models.GetSubscribers(a.Ctx, a.DB, offset, limit)
	if err != nil {
		respond.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	respond.JSON(w, http.StatusOK, subs)
	return
}

// GetSubscribers gets online subscribers
func (a *App) getSubscribersOnline(w http.ResponseWriter, r *http.Request) {

	subs, err := models.GetSubscribersOnline(a.jsonrpcHTTPAddr, a.httpClient)
	if err != nil {
		respond.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, subs)
	return
}

// GetSubscriber gets subscribers
func (a *App) getSubscriber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	subs, err := models.GetSubscriberByID(a.Ctx, a.DB, id)
	if err != nil {
		respond.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	respond.JSON(w, http.StatusOK, subs)
	return
}

func (a *App) deleteSubscriber(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"message": "subscriber successfully deleted", "status": "success"}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	_, err := models.GetSubscriberByID(a.Ctx, a.DB, id)
	if err != nil {
		respond.ERROR(w, http.StatusNotFound, err)
		return
	}

	err = models.DeleteSubscriber(a.Ctx, a.DB, id)
	if err != nil {
		respond.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, resp)
	return
}

// CreateVenue parses request, validates data and saves the new venue
func (a *App) createSubscriber(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"message": "sip device successfully created", "status": "success"}

	subscriber := &models.Subscribers{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if err = json.Unmarshal(body, &subscriber); err != nil {
		respond.ERROR(w, http.StatusBadRequest, err)
		return
	}

	subscriber.Prepare()

	if sbname := models.GetSubscriberByUserName(a.Ctx, a.DB, subscriber.Username); sbname != nil {
		respond.ERROR(w, http.StatusBadRequest, sbname)
		return
	}

	if err = subscriber.Validate(); err != nil {
		respond.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if err = subscriber.Save(a.Ctx, a.DB); err != nil {
		respond.ERROR(w, http.StatusBadRequest, err)
		return
	}

	respond.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) updateSubscriber(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"message": "sip device successfully updated", "status": "success"}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	_, err := models.GetSubscriberByID(a.Ctx, a.DB, id)
	if err != nil {
		respond.ERROR(w, http.StatusNotFound, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond.ERROR(w, http.StatusBadRequest, err)
		return
	}

	subscriber := models.Subscribers{}
	if err = json.Unmarshal(body, &subscriber); err != nil {
		respond.ERROR(w, http.StatusBadRequest, err)
		return
	}

	subscriber.Prepare()

	err = subscriber.UpdateSubscriber(a.Ctx, a.DB, id)
	if err != nil {
		respond.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusOK, resp)
	return
}
