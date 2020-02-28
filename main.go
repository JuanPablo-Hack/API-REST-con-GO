package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Creamos la estrucutura de nuestro evento
type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

//Declaramos un evento vacio con nuestros eventos
type allEvents []event

//Alimentamos de datos nuestro arreglo de eventos
var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduccion de una API-REST con GO",
		Description: "Ven y aprende como crear una API rest desde 0 con golang",
	},
}

//Declaramos las rutas
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido a mi API-REST con Go")
}

//Creamos la funcion de un evento nuevo (metodo POST)
func createEvent(w http.ResponseWriter, r *http.Request) {
	//Declaramos la variable de un nuevo evento
	var newEvent event
	//Decimos que lea el contenido
	reqBody, err := ioutil.ReadAll(r.Body)
	//Validacion de que los datos se estan entregando
	if err != nil {
		fmt.Fprintf(w, "Ingrese los datos con el título y la descripción del evento solo para actualizar")
	}
	//Lo convertimos en formato json
	json.Unmarshal(reqBody, &newEvent)
	//Agregamos el valor en el arreglo de eventos
	events = append(events, newEvent)
	//Escribimos el estado de la peticion
	w.WriteHeader(http.StatusCreated)
	//Mostramos la informacion que se ingreso
	json.NewEncoder(w).Encode(newEvent)
}

//Obtenemos los datos de un evento en especifico metodo GET
func getOneEvent(w http.ResponseWriter, r *http.Request) {
	//Declaramos el id del evento
	eventID := mux.Vars(r)["id"]
	//Recorremos el arreglo para ver si existe y lo mostramos
	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

//Obtenemos todos los eventos metodo GET
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

//Funcion para actualizar un evento
func updateEvent(w http.ResponseWriter, r *http.Request) {
	//Declaramos el id que se va a actualizar
	eventID := mux.Vars(r)["id"]
	//Declramos una variable de evento a actualizar
	var updatedEvent event
	//Si el cuerpo del request esta bien que lo lea
	reqBody, err := ioutil.ReadAll(r.Body)
	// Si existe algun error entonces
	if err != nil {
		fmt.Fprintf(w, "Ingrese los datos con el título y la descripción del evento solo para actualizar")
	}
	//Mandamos en formta JSON el cuerpo del request y el evento a actualizar
	json.Unmarshal(reqBody, &updatedEvent)
	//Recorremos todo el array de eventos para poder ingresar la informacion
	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	// Obtenga la identificación de la url
	eventID := mux.Vars(r)["id"]

	// Obtenga los detalles de un evento existente
	// Use el identificador en blanco para evitar crear un valor que no se usará
	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func main() {
	//Iniciamos el servidor
	router := mux.NewRouter().StrictSlash(true)
	//Mandamos llamar a las rutas
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))

}
