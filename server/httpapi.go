package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type httpApi struct {
	store MapStore
}

func NewHttpApi(store MapStore) *httpApi {
	return &httpApi{
		store: store,
	}
}

func (api *httpApi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	switch req.Method {
	case http.MethodGet:
		api.handleGetAll(w, req)
	case http.MethodPut:
		api.handlePut(w, req)
	case http.MethodDelete:
		api.handleDelete(w, req)
	default:
		w.Header().Add("Allow", http.MethodGet)
		w.Header().Set("Allow", http.MethodPut)
		w.Header().Add("Allow", http.MethodDelete)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (api *httpApi) handleGetAll(w http.ResponseWriter, req *http.Request) {
	ips, err := api.store.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read bpfMap: %s", err), http.StatusInternalServerError)
		return
	}

	jsonIPs, err := json.Marshal(ips)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal struct: %s", err), http.StatusInternalServerError)
		return
	}

	w.Write(jsonIPs)
}

func (api *httpApi) handlePut(w http.ResponseWriter, req *http.Request) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	if err := api.store.Put(string(b)); err != nil {
		http.Error(w, fmt.Sprintf("failed to write to bpfMap: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (api *httpApi) handleDelete(w http.ResponseWriter, req *http.Request) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	if err := api.store.Delete(string(b)); err != nil {
		http.Error(w, fmt.Sprintf("failed to write to bpfMap: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
