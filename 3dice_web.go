package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"wojones.com/src/dicegame"

	"github.com/flosch/pongo2"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// var tpl = template.Must(template.New("index.html").Funcs(template.FuncMap{"JoinStrings": joinem}).ParseFiles("static/index.html"))
// var ptpl, err = pongo2.FromString("<h1>hello {{name}}</h1>")
var ptpl = pongo2.Must(pongo2.FromFile("static/index.html"))

var router *chi.Mux

func setupRoutes() error {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)

	router.Get("/", indexHandler)
	router.Route("/games", func(r chi.Router) {
		r.Route("/{gameID}", func(r chi.Router) {
			r.Use(gameCtx)
			r.Get("/", getGame)
		})
	})

	fs := http.FileServer(http.Dir("static/assets"))
	router.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	return nil
}

func errcatch(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
func webListen() error {
	err := http.ListenAndServe(":3002", router)
	errcatch(err)
	return err
}

func parseargs() pongo2.Context {
	return pongo2.Context{"name": "jack", "dicegame": tdg, "start": starttime.Format(time.DateTime)}
}

func gameCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "gameID")
		fmt.Printf("gameCtx: game %s\n", gameID)
		// FOR NOW: There's only one game ...
		game := &tdg
		ctx := context.WithValue(r.Context(), "game", game)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getGame(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("getGame: %s\n", r.URL.String())
	gp, ok := r.Context().Value("game").(*dicegame.DiceGame)

	if !ok {
		fmt.Printf("Type assertion failure (%v)\n", gp)
		http.Error(w, http.StatusText(422), 422)
		return
	}
	fmt.Printf("getGame: got game %s\n", gp.ID)
	e_err := ptpl.ExecuteWriter(parseargs(), w)
	if e_err != nil {
		http.Error(w, e_err.Error(), http.StatusInternalServerError)
	}
}

func dumpGame(foo *dicegame.DiceGame) {
	fmt.Printf("Game: %s\n", foo.ID)
}

/*
func gamesHandler(w http.ResponseWriter, r *http.Request) {
	// FOR NOW: there's only one game ...
	fmt.Printf("gameHandler: ")
}
*/

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Sending index.html ...\n")
	fmt.Printf("Handling index (url: %s)\n", r.URL.String())
	err := ptpl.ExecuteWriter(parseargs(), w)
	fmt.Printf("THWACK! Writer executed!\n")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// tpl.Execute(w, tdg)
	//w.Write([]byte("<h1>Talking shit?!!</h1>"))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handling search (url: %s)\n", r.URL.String())

	/* This is for when it was a GET request
	u, err := url.Parse(r.URL.String())
	if err != nil {
		fmt.Printf("Error decoding url %s: %s\n", r.URL.String, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := u.Query()
	command := params.Get("move")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}
	*/
	if err := r.ParseForm(); err != nil {
		fmt.Printf("error parsing form: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	command := r.FormValue("move")
	page := r.FormValue("page")
	if page == "" {
		page = "1"
	}

	fmt.Println("Command is: ", command)
	fmt.Println("Page is: ", page)

	runcmd(&tdg, command)
	//tpl.Execute(w, tdg)

	e_err := ptpl.ExecuteWriter(parseargs(), w)
	fmt.Printf("Written!\n")
	if e_err != nil {
		http.Error(w, e_err.Error(), http.StatusInternalServerError)
	}
}
