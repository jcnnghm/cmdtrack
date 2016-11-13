package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
	"github.com/jcnnghm/cmdtrack/cmd"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

// Secret holds the application secret key
type Secret struct {
	SecretKey string
}

// fetchSecret gets or creates a new secret in appengine, and returns it
func fetchSecret(dc context.Context) (s *Secret, err error) {
	err = datastore.RunInTransaction(dc, func(tc context.Context) error {
		secretKey := datastore.NewKey(tc, "Secret", "secret", 0, nil)
		s = new(Secret)
		if err := datastore.Get(tc, secretKey, s); err != nil {
			if err != datastore.ErrNoSuchEntity {
				return err
			}
			b := make([]byte, 16)
			if _, err := rand.Read(b); err == nil {
				key := hex.EncodeToString(b)
				s.SecretKey = key
				if _, err := datastore.Put(tc, secretKey, s); err != nil {
					return errors.New("Unable to save secret")
				}
			} else {
				return errors.New("Unable to generate secret")
			}
		}
		return nil
	}, &datastore.TransactionOptions{})
	return
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
			log.Debugf(c, "Secret key mismatch")
			http.Error(w, "Secret key invalid", http.StatusUnauthorized)
			return
		}

		fn(w, r)
	}
}

func logCommand(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	log.Debugf(c, "Log Command Started")

	if command, err := cmd.NewCommand(r); err != nil || !command.IsValid() {
		if err != nil {
			log.Debugf(c, "Failed with error: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Invalid command data", http.StatusInternalServerError)
		}
	} else {
		log.Debugf(c, "Received Command: %#v", command)
		err := saveCommand(command, c)
		if err != nil {
			log.Debugf(c, "Failed with error: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func history(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	limit := 10000
	log.Debugf(c, "Starting History Query")
	q := datastore.NewQuery("HistoryLine").Order("-Timestamp").Limit(limit).EventualConsistency()
	logLines := make([]cmd.Command, 0, limit)
	if _, err := q.GetAll(c, &logLines); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Debugf(c, "Found Log Lines: %v", len(logLines))
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(logLines); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func historyKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "HistoryLog", "default_history_log", 0, nil)
}

func saveCommand(command *cmd.Command, c context.Context) error {
	key := datastore.NewIncompleteKey(c, "HistoryLine", nil)
	_, err := datastore.Put(c, key, command)
	log.Debugf(c, "Put Command Complete")
	return err
}
