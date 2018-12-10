package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"verisart-api/internal/models"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// VerisartAPI represent the high level API
type VerisartAPI struct {
	certificateModel *models.CertificateModel
	userModel        *models.UserModel
	ownerID          string
}

func main() {
	log.Print("Started Verisart API backend.")

	userModel := models.NewUserModel()
	certificateModel := models.NewCertificateModel(userModel)

	api := &VerisartAPI{
		certificateModel: certificateModel,
		userModel:        userModel,
	}

	api.start()
}

func accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s ", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// Certificates endpoint require an authenticated owner, authMiddleware validates that incoming user
func (v *VerisartAPI) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		ownerID := r.Header.Get("X-Owner-ID")

		if len(ownerID) != 0 {
			log.Printf("Verifying X-Owner-ID user %s", ownerID)
			_, err = v.userModel.GetUserByID(ownerID)
			if err == nil {
				v.ownerID = ownerID
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}

func (v *VerisartAPI) start() {
	router := mux.NewRouter()
	// Printing Access Logs
	router.Use(accessLogMiddleware)

	users := router.PathPrefix("/users").Subrouter()
	users.HandleFunc("/", v.ListUsers).Methods("GET")
	users.HandleFunc("/", v.CreateUser).Methods("POST")
	users.HandleFunc("/{id}/certificates/", v.GetUserCertificates).Methods("GET")

	certRoutes := router.PathPrefix("/certificates").Subrouter()
	certRoutes.Use(v.authMiddleware)

	certRoutes.HandleFunc("/", v.CreateCertificate).Methods("POST")
	certRoutes.HandleFunc("/{id}/", v.UpdateCertificate).Methods("PUT")
	certRoutes.HandleFunc("/{id}/", v.DeleteCertificate).Methods("DELETE")
	certRoutes.HandleFunc("/{id}/transfers/", v.CreateTransfer).Methods("POST")
	certRoutes.HandleFunc("/{id}/transfers/", v.AcceptTransfer).Methods("PUT")

	// CORS
	// Unsure where from we should limit the calls so we go with * ?
	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "X-Owner-ID"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	err := http.ListenAndServe(":8000", handlers.CORS(origins, headers, methods)(router))
	if err != nil {
		log.Fatalf("Failed to start HTTP server on at localhost:8000")
	}
}

// CreateCertificate handles the POST endpoint for /certificates, it creates a valid Verisart certificate
func (v *VerisartAPI) CreateCertificate(w http.ResponseWriter, r *http.Request) {
	var certificate *models.Certificate
	_ = json.NewDecoder(r.Body).Decode(&certificate)

	certificate.OwnerID = v.ownerID
	certificate.CreatedAt = time.Now()

	err := v.certificateModel.CreateCertificate(certificate)
	if err != nil {
		log.Printf("An error occured creating certificate, err: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(certificate)
}

// UpdateCertificate handles the PUT endpoint /certificates/:id, it updates an existing Verisart certificate
func (v *VerisartAPI) UpdateCertificate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	certID := vars["id"]

	var certificate *models.Certificate
	_ = json.NewDecoder(r.Body).Decode(&certificate)

	if v.ownerID != certificate.OwnerID {
		log.Printf("[CODE RED] Owner mismatch on certificate update, expected %s, got %s", certificate.OwnerID, v.ownerID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if certID != certificate.ID {
		log.Printf("[CODE RED] Certificate ID mismatch on certificate update, expected %s, got %s", certID, certificate.ID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err := v.certificateModel.UpdateCertificate(certificate)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(certificate)
}

// DeleteCertificate handles the DELETE endpoint /certificates/:id, it deletes an existing Verisart certificate
func (v *VerisartAPI) DeleteCertificate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	certID := vars["id"]

	ok, err := v.certificateModel.CheckCertificateOwnership(certID, v.ownerID)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	if ok {
		v.certificateModel.DeleteCertificate(certID)
	} else {
		log.Printf("[CODE RED] OwnerID mistmatch on certificate deletion. Owner %s attempted to delete certificate %s", v.ownerID, certID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CreateTransfer handles the POST endpoint /certificates/:id/transfers, it initiates the transfer of a certificate
func (v *VerisartAPI) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	certID := vars["id"]

	var transfer models.Transfer
	_ = json.NewDecoder(r.Body).Decode(&transfer)

	ok, err := v.certificateModel.CheckCertificateOwnership(certID, v.ownerID)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if ok {
		transfer.Status = models.Pending
		v.certificateModel.CreateTransfer(certID, transfer)
	} else {
		log.Printf("[CODE RED] OwnerID mistmatch on certificate transfer. Owner %s attempted to transfer certificate %s", v.ownerID, certID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// AcceptTransfer handles the PUT endpoint /certificates/:id/transfers, it confirms the transfer of a certificate
func (v *VerisartAPI) AcceptTransfer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	certID := vars["id"]

	err := v.certificateModel.AcceptTransfer(certID)
	if err != nil {
		log.Printf("An error occured confirming a transfer of certificate, err: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//ListUsers handles the GET endpoint /users
func (v *VerisartAPI) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := v.userModel.GetUsers()

	json.NewEncoder(w).Encode(users)
}

// CreateUser handles the POST endpoint /users/
func (v *VerisartAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	err := v.userModel.CreateUser(&user)
	if err != nil {
		log.Printf("Failed to create user %s, err: %s", user, err)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

// GetUserCertificates handles the GET endpoint for
func (v *VerisartAPI) GetUserCertificates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	certificates := v.certificateModel.GetCertificatesByOwnerID(vars["id"])

	json.NewEncoder(w).Encode(certificates)
}
