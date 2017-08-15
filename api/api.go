package api

import(
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
)

var Database *sql.DB

type ApiServer struct{}

// Para los errores
type CreateResponse struct {
  Error string "json:error"
  ErrorCode int "json:code"
}

/*------------------------------------------------------------------------------
 main


------------------------------------------------------------------------------*/
func (api *ApiServer) Start(){

  routes := mux.NewRouter()
  routes.HandleFunc("/api/eventos", ConsultarEvento).Methods("GET")
  routes.HandleFunc("/api/eventos", AltaEvento).Methods("POST")
  routes.HandleFunc("/api/eventos", ModificarEvento).Methods("PUT")
  routes.HandleFunc("/api/eventos", BorrarEvento).Methods("DELETE")
  routes.HandleFunc("/api/usuariosevento", AltaUsuariosEvento).Methods("POST")
  routes.HandleFunc("/api/usuariosevento", BajaUsuariosEvento).Methods("DELETE")

  fmt.Println("Listening on port 8080...")
/*
  routes.HandleFunc("/api/users", UsersRetrieve).Methods("GET")
  routes.HandleFunc("/api/users/{id:[0-9]+}", UserRetrieve).Methods("GET")
  routes.HandleFunc("/api/users/{id:[0-9]+}", UsersUpdate).Methods("PUT")
  routes.HandleFunc("/api/users", UsersInfo).Methods("OPTIONS")

  routes.HandleFunc("/authorize", ApplicationAuthorize).Methods("POST")
  routes.HandleFunc("/authorize", ApplicationAuthenticate).Methods("GET")

  routes.HandleFunc("/authorize/{service:[a-z]+}", ServiceAuthorize).Methods("GET")
  routes.HandleFunc("/connect/{service:[a-z]+}", ServiceConnect).Methods("GET")
  routes.HandleFunc("/oauth/token", CheckCredentials).Method("POST")
*/
  http.Handle("/", routes)
  http.ListenAndServe(":8080", nil)
}
