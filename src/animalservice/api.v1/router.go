package api

import (
	"github.com/PennTex/PetWhisperer/src/animalservice"
	"github.com/PennTex/PetWhisperer/src/animalservice/repositories"
	"github.com/gorilla/mux"
)

func NewRouter(router *mux.Router) *mux.Router {
	animalRepository := repositories.CloudDatastoreRepository{}
	animalService := animalservice.NewAnimalService(animalRepository)
	animalAPI := NewAnimalAPI(animalService)

	router.HandleFunc("/animals", animalAPI.PostAnimal).Methods("POST")
	router.HandleFunc("/animals", animalAPI.GetAnimals).Methods("GET")
	router.HandleFunc("/animals/{animalID:[0-9a-z-]{36}}", animalAPI.GetAnimal).Methods("GET")
	router.HandleFunc("/animals/{animalID:[0-9a-z-]{36}}", animalAPI.DeleteAnimal).Methods("DELETE")
	router.HandleFunc("/users/{userID}/animals", animalAPI.GetUsersAnimals).Methods("GET")

	return router
}
