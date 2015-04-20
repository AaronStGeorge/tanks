package main

import (
	"errors"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var ErrUserNotPresent = errors.New("User not present in context")
var ErrFriendsNotPresent = errors.New("Friends not present in context")
var ErrSessionNotPresent = errors.New("No session found")

type contextKey int

// Define keys to retrieve values from context
const userKey contextKey = 0
const friendKey contextKey = 1
const sessionKey contextKey = 2

func GetFrends(r *http.Request) ([]User, error) {
	val, ok := context.GetOk(r, friendKey)
	if !ok {
		return nil, ErrFriendsNotPresent
	}

	friends, ok := val.([]User)
	if !ok {
		return nil, ErrFriendsNotPresent
	}

	return friends, nil
}

func SetFreinds(r *http.Request, val []User) {
	context.Set(r, friendKey, val)
}

func (g *Global) loadFriends(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		friends := make([]User, 0, 10)

		user, err := GetUser(r)
		if err != nil {
			log.Fatal(err)
		}

		if user.Id != -1 {

			rows, err := g.db.Query("SELECT user_id_2, Username FROM "+
				"Friends INNER JOIN Users "+
				"ON Friends.user_id_2=Users.id "+
				"WHERE user_id_1 = ?", user.Id)
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				var id int
				var userName string
				err = rows.Scan(&id, &userName)
				if err != nil {
					log.Fatal(err)
				}
				friends = append(friends, User{UserName: userName, Id: id})
			}
			rows.Close()
		}

		SetFreinds(r, friends)

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetUser(r *http.Request) (User, error) {
	val, ok := context.GetOk(r, userKey)
	if !ok {
		return User{}, ErrUserNotPresent
	}

	user, ok := val.(User)
	if !ok {
		return User{}, ErrUserNotPresent
	}

	return user, nil
}

func SetUser(r *http.Request, val User) {
	context.Set(r, userKey, val)
}

func (g *Global) loadUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		session, err := GetSession(r)
		if err != nil {
			log.Fatal(err)
		}

		user := User{}
		if session.IsNew {
			user.Id = -1
		} else {
			user.Id = session.Values["Id"].(int)
			user.UserName = session.Values["UserName"].(string)
		}
		context.Set(r, userKey, user)

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	val, ok := context.GetOk(r, sessionKey)
	if !ok {
		return nil, ErrUserNotPresent
	}

	session, ok := val.(*sessions.Session)
	if !ok {
		return nil, ErrUserNotPresent
	}

	return session, nil
}

func SetSession(r *http.Request, val *sessions.Session) {
	context.Set(r, sessionKey, val)
}

func (g *Global) loadSession(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		session, err := g.store.Get(r, "tanks-app")
		if err != nil {
			log.Fatal(err)
		}

		SetSession(r, session)

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
