package handler

import (
	"awesomeProject/pkg/repository"
	"encoding/json"
	"fmt"
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

func (h *Handler) generateLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	s := r.Form
	id := s["id"][0]

	portId, actualUrl, err := h.repository.ProxyPorts.GenerateSlug(id)

	if err != nil {
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
	r_router_id := s["id"][0]
	r_interval := s["interval"][0]

	system_port_id, _ := strconv.Atoi(r_router_id)
	router_id, err := h.repository.UpdateReconnectInterval(system_port_id, r_interval)

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
	SetInterval(string(r_router_id), r_interval)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]int)
	resp["Ok"] = router_id
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (h *Handler) setConfig(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// create user -> post to localserver to 3ProxyConfig handler

	r.ParseForm()
	s := r.Form
	id := s["id"][0]
	login := s["login"][0]
	password := s["password"][0]
	err := h.repository.CreateSimpleUser(id, login, password)
	SetConfig(id, login, password)
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

}
func (h *Handler) removeinterval(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	s := r.Form
	r_router_id := s["id"][0]

	fmt.Println(r_router_id)
	DeleteInterval(r_router_id)
	// Output:
	// http://api.themoviedb.org/3/tv/popular?another_thing=foo+%26+bar&api_key=key_from_environment_or_flag
}

func (h *Handler) ProxyServerHealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	HealthCheck()
}
func (h *Handler) Reboot(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// add password check to reboot router by id
	r.ParseForm()
	s := r.Form
	id := s["id"][0]
	body, err := Reboot(id)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	resp := make(map[string]string)
	resp["message"] = body
	// jsonResp, err := json.Marshal(resp)
	// if err != nil {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// }
	w.Write([]byte(body))
}

func (h *Handler) rebootByLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	router_id, _ := h.repository.ProxyPorts.GetIdBySlug(ps.ByName("url"))

	resp := make(map[int]string)
	if router_id != 0 {
		// INSERT LOGIC OF RUNING REBOOT OF ROUTER
		resp[router_id] = "Proxy reloading now!"

		// TODO handel h error
		// select router id where generatedUrl = data
		id := strconv.Itoa(router_id)
		_, err := Reboot(id)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
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

// TODO local server var in config viper config git
// TODO send correctly get to local server ip
// TODO config
// TODO get ip adres bellow reload and after (python logic) and log this (python logic)
// change by link - 12.07.22 15:50 - айпи адрес который получили после смены

// TODO compare ip addr , want to get uniq addr
// TODO generate random log pass

func (h *Handler) InitRoutes() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", h.HealthCheck)
	router.GET("/localServerHealth", h.ProxyServerHealthCheck)
	router.GET("/generateLink", h.rebootByLink)
	router.GET("/rebootById", h.Reboot)
	router.GET("/rebootByLink/:link", h.generateLink)
	// post request
	// 
	// 
	router.POST("/setConfig", h.setConfig)
	router.POST("/updateInterval", h.updateInterval)
	router.POST("/removeinterval", h.removeinterval)
	return router
}
