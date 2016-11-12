package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcnnghm/cmdtrack/cmd"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// Secret holds the application secret key
type Secret struct {
	SecretKey string
}

// fetchSecret gets or creates a new secret in appengine, and returns it
func fetchSecret(c appengine.Context) (*Secret, error) {
	secretKey := datastore.NewKey(c, "Secret", "secret", 0, nil)
	s := new(Secret)
	if err := datastore.Get(c, secretKey, s); err != nil {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err == nil {
			key := hex.EncodeToString(b)
			s.SecretKey = key
			if _, err := datastore.Put(c, secretKey, s); err != nil {
				return nil, errors.New("Unable to save secret")
			}
		} else {
			return nil, errors.New("Unable to generate secret")
		}
	}
	return s, nil
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/secret", secret).Methods("GET")
	r.HandleFunc("/command", logCommand).Methods("POST")
	http.Handle("/", r)
}

// secret gets or creates a secret.  auth needs to be added to this endpoint.
func secret(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	u := user.Current(c)
	if u == nil || !u.Admin {
		http.Error(w, "Admin login only", http.StatusUnauthorized)
		return
	}

	if s, err := fetchSecret(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, s.SecretKey)
	}
}

func logCommand(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if command, err := cmd.NewCommand(r); err != nil || !command.IsValid() {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Invalid command data", http.StatusInternalServerError)
		}
	} else {
		c.Debugf("Received Command: %#v", command)
		err := saveCommand(command, c)
		if err != nil {
			c.Debugf("Failed with error: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func saveCommand(command *cmd.Command, c appengine.Context) error {
	parentKey := datastore.NewKey(c, "HistoryLog", "default_history_log", 0, nil)
	key := datastore.NewIncompleteKey(c, "HistoryLine", parentKey)
	_, err := datastore.Put(c, key, command)
	return err
}
