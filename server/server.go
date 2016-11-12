package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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
	r.HandleFunc("/history", verifySecret(history)).Methods("GET")
	r.HandleFunc("/command", verifySecret(logCommand)).Methods("POST")
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

func verifySecret(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)

		u := user.Current(c)
		if u != nil && u.Admin {
			fn(w, r)
			return
		}

		secret, err := fetchSecret(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(r.Header["Secret"]) != 1 || r.Header["Secret"][0] != secret.SecretKey {
			c.Debugf("Secret key mismatch")
			http.Error(w, "Secret key invalid", http.StatusUnauthorized)
			return
		}

		fn(w, r)
	}
}

func logCommand(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if command, err := cmd.NewCommand(r); err != nil || !command.IsValid() {
		if err != nil {
			c.Debugf("Failed with error: %v", err.Error())
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

func history(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	limit := 10000
	q := datastore.NewQuery("HistoryLine").Ancestor(historyKey(c)).Order("-Timestamp").Limit(limit).EventualConsistency()
	logLines := make([]cmd.Command, 0, limit)
	if _, err := q.GetAll(c, &logLines); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Debugf("Found Log Lines: %v", len(logLines))
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(logLines); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func historyKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "HistoryLog", "default_history_log", 0, nil)
}

func saveCommand(command *cmd.Command, c appengine.Context) error {
	key := datastore.NewIncompleteKey(c, "HistoryLine", historyKey(c))
	_, err := datastore.Put(c, key, command)
	return err
}
