package animalservice

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

type animalPostReq struct {
	Typ      string   `json:"type"`
	Name     string   `json:"name"`
	Birthday int64    `json:"birthday"`
	Owners   []string `json:"owners"`
	ImageURL string   `json:"image_url"`
}

type response struct {
	Error interface{} `json:"error"`
	Data  interface{} `json:"data"`
}

var animalRepo cloudDatastoreRepository

func sendResponse(w http.ResponseWriter, r *http.Request, status int, errorPayload interface{}, dataPayload interface{}) {
	message, err := json.Marshal(response{
		Error: errorPayload,
		Data:  dataPayload,
	})

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(message)
}

func getUsersAnimals(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	ownerID := mux.Vars(r)["userID"]

	animals, err := animalRepo.getByOwnerID(ctx, ownerID)

	if err != nil {
		sendResponse(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	log.Infof(ctx, "Animals retrieved: %s", animals)

	sendResponse(w, r, http.StatusOK, nil, animals)
}

func getAnimals(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	animalIDs := r.URL.Query()["animalID"]

	if animalIDs != nil {
		var animals []animal

		for _, animalID := range animalIDs {
			animal, err := animalRepo.getByID(ctx, animalID)

			if err != nil {
				sendResponse(w, r, http.StatusInternalServerError, err.Error(), nil)
				return
			}

			animals = append(animals, *animal)
		}

		sendResponse(w, r, http.StatusOK, nil, animals)
	} else {
		animals, err := animalRepo.get(ctx)

		if err != nil {
			sendResponse(w, r, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		sendResponse(w, r, http.StatusOK, nil, animals)
	}
}

func getAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	animalID := mux.Vars(r)["animalID"]

	animal, err := animalRepo.getByID(ctx, animalID)

	if err != nil {
		sendResponse(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	sendResponse(w, r, http.StatusOK, nil, animal)
}

func postAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	animalReq := animalPostReq{}

	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &animalReq)

	animalID, err := animalRepo.create(ctx, &animal{
		Typ:       animalReq.Typ,
		Name:      animalReq.Name,
		Birthday:  animalReq.Birthday,
		Owners:    animalReq.Owners,
		ImageURL:  animalReq.ImageURL,
		CreatedAt: time.Now().Unix(),
	})

	if err != nil {
		sendResponse(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	log.Infof(ctx, "New animal ID: %s", animalID)
	sendResponse(w, r, http.StatusOK, nil, animalID)
}

func deleteAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	animalID := mux.Vars(r)["animalID"]

	err := animalRepo.destroy(ctx, animalID)
	if err != nil {
		sendResponse(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	sendResponse(w, r, http.StatusNoContent, nil, nil)
}
