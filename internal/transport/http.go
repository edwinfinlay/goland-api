package transport

import (
	"encoding/json"
	"log"
	"my-first-api/internal/todo"
	"net/http"
)

type TodoItem struct {
	//ID   int    `json:"id"`
	Item string `json:"item"`
}

type Server struct {
	mux *http.ServeMux
}

func NewServer(todoSvc *todo.Service) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		todoItems, err := todoSvc.GetAll()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(todoItems)
		if err != nil {
			log.Println(err)
		}

		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}

	})

	mux.HandleFunc("POST /todo", func(writer http.ResponseWriter, request *http.Request) {
		var t TodoItem
		err := json.NewDecoder(request.Body).Decode(&t)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		//todos[t.ID] = t.Item
		err = todoSvc.Add(t.Item)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusCreated)
		return
	})

	//mux.HandleFunc("DELETE /todo/{id}", func(writer http.ResponseWriter, request *http.Request) {
	//	id := request.PathValue("id")
	//	fmt.Printf("id: %s\n", id)
	//	intID, err := strconv.Atoi(id)
	//	if err != nil {
	//		log.Println(err)
	//		writer.WriteHeader(http.StatusBadRequest)
	//		return
	//	}
	//	delete(todos, intID)
	//	writer.WriteHeader(http.StatusNoContent)
	//	return
	//})

	mux.HandleFunc("GET /search", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query().Get("q")
		if query == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		results, err := todoSvc.Search(query)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(results)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
	})

	return &Server{
		mux: mux,
	}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(":8181", s.mux)
}
