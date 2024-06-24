package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prebid/openrtb/v20/openrtb2"
	"io"
	"net/http"
	"simpleSeller/middleware"
	"simpleSeller/utils"
)

type Handler struct {
	Router *mux.Router
	Http   *http.Server
}

func NewHandler() *Handler {
	h := &Handler{}
	h.Router = mux.NewRouter()
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
}

func postBid(w http.ResponseWriter, r *http.Request) {

	// Read the entire request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Parse the request body into an OpenRTB BidRequest struct
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
	fmt.Println("bid tpe is ", bidType)
	//middleware.MiddlewareLog(bidType)

	// Process the bid request and generate a BidResponse (this part is application-specific)
	bidResponse := utils.GenerateBidResponse(&bidRequest)
	// Encode the BidResponse to JSON
	responseJSON, err := json.Marshal(bidResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
		return
	}

	// Set response content type and write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)

}

func (h *Handler) InitServer() error {
	if err := h.Http.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
