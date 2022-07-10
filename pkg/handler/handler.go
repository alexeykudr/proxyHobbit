package handler

import (
	"awesomeProject/pkg/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

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

func ping(ctx context.Context, ip string) {
	cmd := exec.CommandContext(ctx, "ping", ip, "-t")
	cmd.Stdout = os.Stdin

	go cmd.Run()
	<-ctx.Done()
	fmt.Println("Pid of procces", cmd.Process.Pid)
	cmd.Process.Kill()
}

func (h *Handler) execProxyCommand() {
	ctx, cancel := context.WithTimeout(context.Background(), 5500*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "sleep", "5").Run(); err != nil {
		log.Fatalf("Error in execProxyCommand %s", err.Error())
		fmt.Println("err!")
	}
}

func (h *Handler) reconnectHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data, _ := h.repository.ProxyPorts.GetIdBySlug(ps.ByName("url"))

	resp := make(map[int]string)
	if data != 0 {
		// INSERT LOGIC OF RUNING REBOOT OF ROUTER
		resp[data] = "now reloading!"

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
	// fmt.Println(request_port)

	portId, actualUrl, err := h.repository.ProxyPorts.CreateSlugUrl(request_port[0])

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
	portId := s["portId"][0]
	interval := s["interval"][0]

	// fmt.Printf("%T\n", portId)
	// fmt.Println(portId, interval)

	_, err := h.repository.UpdateReconnectInterval(portId, interval)
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

	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (h *Handler) InitRoutes() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", h.HealthCheck)
	router.GET("/reboot/:url", h.reconnectHandler)
	router.GET("/generateSlug/", h.generateSlug)
	router.POST("/updateInterval", h.updateInterval)

	// log.Fatal(http.ListenAndServe(":8080", router))
	return router
}
