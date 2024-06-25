package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"simpleSeller/middleware"
	"simpleSeller/ratelimiter"
	"simpleSeller/utils"

	"github.com/gorilla/mux"
	"github.com/prebid/openrtb/v20/openrtb2"
)

type Handler struct {
	Router      *mux.Router
	Http        *http.Server
	Ratelimiter *ratelimiter.RedisRateLimiter
}

func NewHandler() *Handler {
	h := &Handler{
		Router:      mux.NewRouter(),
		Ratelimiter: ratelimiter.NewRedisRateLimiter(),
	}

	h.mapRoutes()
	h.Http = &http.Server{
		Addr:    ":8080",
		Handler: h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/post-bid", postBid).Methods("POST")
	h.Router.Use(middleware.Validate())
	h.Router.Use(ratelimiter.RunLimit(h.Ratelimiter))

}

func postBid(w http.ResponseWriter, r *http.Request) {

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Parse the request body into an OpenRTB BidRequest
	var bidRequest openrtb2.BidRequest
	err = json.Unmarshal(body, &bidRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}
	var bidType string
	if bidRequest.Site != nil {
		bidType = "site"
	} else {
		bidType = "app"
	}
	fmt.Println("bid type is ", bidType)

	bidResponse := utils.GenerateBidResponse(&bidRequest)
	// Encode the BidResponse to JSON
	responseJSON, err := json.Marshal(bidResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
		return
	}

	// write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)

}

func (h *Handler) InitServer() error {
	if err := h.Http.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
