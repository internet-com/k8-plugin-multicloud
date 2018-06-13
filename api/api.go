/*
Copyright 2018 Intel Corporation.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	pkgerrors "github.com/pkg/errors"
)

// Start VNFInstance WebService
func Start() {
	service, err := NewVNFInstanceService("")
	if err != nil {
		werr := pkgerrors.Wrap(err, "Creation of a service error")
		log.Panic(werr)
	}

	router := mux.NewRouter()
	router.HandleFunc("/v1/vnf_instances", service.CreateVNF).Methods("POST")
	router.HandleFunc("/v1/vnf_instances", service.ListVNF).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Println("[INFO] Started Kubernetes Multicloud API")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter)) // Remove hardcode.
}