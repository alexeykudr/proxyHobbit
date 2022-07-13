package handler

import (
	"awesomeProject/pkg/repository"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	repository *repository.Repository
}

func NewHandler(repository *repository.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["ok"] = "true"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (h *Handler) rebootByUrl(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data, _ := h.repository.ProxyPorts.GetIdBySlug(ps.ByName("url"))

	resp := make(map[int]string)
	if data != 0 {
		// INSERT LOGIC OF RUNING REBOOT OF ROUTER
		resp[data+10] = "Proxy reloading now!"

		// TODO handel h error
		h.execProxyCommand()
		w.WriteHeader(http.StatusOK)
	} else {
		resp[-1] = "error, cant find"
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(resp)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
	// select id from proxyPorts where generated url like slug
}

func (h *Handler) generateSlug(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	r.ParseForm()
	s := r.Form
	request_port := s["id"]

	a, _ := strconv.Atoi(request_port[0])
	system_port_id := a
	// fmt.Println(request_port)

	portId, actualUrl, err := h.repository.ProxyPorts.GenerateSlug(system_port_id)

	if err != nil {
		// log.Printf("Error happened in CreatingSlug Err: %s", err)
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Error"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[int]string)
	if portId >= 1 {
		resp[portId] = actualUrl
		// h.execProxyCommand()
	} else {
		resp[0] = "error, cant find"
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

func (h *Handler) updateInterval(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	s := r.Form
	portId := s["id"][0]
	interval := s["interval"][0]

	a, _ := strconv.Atoi(portId)
	system_port_id := a

	// fmt.Printf("%T\n", portId)
	// fmt.Println(portId, interval)

	returned_interval, err := h.repository.UpdateReconnectInterval(system_port_id, interval)
	// fmt.Println(id)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Bad Request"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]int)
	resp["message"] = returned_interval
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

type test_struct struct {
	Test string
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// decoder := json.NewDecoder(r.Body)
	// var t test_struct
	// err := decoder.Decode(&t)
	// if err != nil {
	//     panic(err)
	// }
	// log.Println(t.Test)

	r.ParseForm()
	s := r.Form
	username := s["login"][0]
	password := s["password"][0]
	err := h.repository.CreateSimpleUser(username, password)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := make(map[string]string)
		resp["message"] = "Bad Request"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	resp := make(map[string]string)
	resp["message"] = "Ok"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

	// TODO UNIQ USERS
}

func (h *Handler) InitRoutes() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", h.HealthCheck)
	router.GET("/reboot/:url", h.rebootByUrl)
	router.GET("/generateSlug/", h.generateSlug)
	router.POST("/updateInterval", h.updateInterval)
	router.POST("/createUser", h.createUser)
	// router.POST isert into db (createProxyPort)
	return router
}
