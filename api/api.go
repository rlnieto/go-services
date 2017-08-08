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
  routes.HandleFunc("/api/cenas", ConsultarCena).Methods("GET")
  routes.HandleFunc("/api/cenas", AltaCena).Methods("POST")
  routes.HandleFunc("/api/cenas", ModificarCena).Methods("PUT")
  routes.HandleFunc("/api/cenas", BorrarCena).Methods("DELETE")
  routes.HandleFunc("/api/usuarioscena", AltaUsuariosCena).Methods("POST")
  routes.HandleFunc("/api/usuarioscena", BajaUsuariosCena).Methods("DELETE")

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
