package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	person "tinder"
)

// HandlerAddSinglePersonAndMatch Add a new user to the matching system and find any possible matches for the new user
func (s *Server) HandlerAddSinglePersonAndMatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p person.Person
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			s.handleError(w, err, http.StatusBadRequest)
			return
		}

		if err := normalizeAndValidatePerson(&p); err != nil {
			s.handleError(w, err, http.StatusBadRequest)
			return
		}

		persons, err := s.personService.AddPersonAndMatch(&p)
		if err != nil {
			s.handleError(w, err, http.StatusInternalServerError)
			return
		}

		s.respond(w, persons, http.StatusOK)
	}
}

// HandlerRemoveSinglePerson Remove a user from the matching system so that the user cannot be matched anymore
func (s *Server) HandlerRemoveSinglePerson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameStr := r.URL.Query().Get("name")
		if nameStr == "" {
			s.handleError(w, errors.New("missing query parameter: 'name'"), http.StatusBadRequest)
			return
		}

		normalizedName := strings.ToLower(strings.TrimSpace(nameStr))

		err := s.personService.RemovePerson(normalizedName)
		if err != nil {
			if err.Error() == person.NotFoundStr {
				s.handleError(w, err, http.StatusNotFound)
				return
			}

			s.handleError(w, err, http.StatusInternalServerError)
			return
		}

		s.respond(w, nil, http.StatusOK)
	}
}

// HandlerQuerySinglePeople Remove a user from the matching system so that the user cannot be matched anymore
func (s *Server) HandlerQuerySinglePeople() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nStr := r.URL.Query().Get("n")
		n, err := strconv.ParseInt(nStr, 10, 32)
		if err != nil {
			s.handleError(w, errors.New("missing query parameter: 'n'"), http.StatusBadRequest)
			return
		}

		persons, err := s.personService.QuerySinglePeople(int(n))
		if err != nil {
			s.handleError(w, err, http.StatusInternalServerError)
			return
		}

		s.respond(w, persons, http.StatusOK)
	}
}

func normalizeAndValidatePerson(p *person.Person) error {
	p.Name = strings.ToLower(strings.TrimSpace(p.Name))
	p.Gender = strings.ToLower(strings.TrimSpace(p.Gender))

	if err := person.Validate.Struct(p); err != nil {
		return err
	}
	return nil
}
