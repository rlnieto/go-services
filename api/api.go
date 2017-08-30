package api

import(
  "database/sql"
  "github.com/gorilla/handlers"
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
  routes.HandleFunc("/api/eventosusuario", EventosUsuario).Methods("GET")
  routes.HandleFunc("/api/usuariosevento", AltaUsuariosEvento).Methods("POST")
  routes.HandleFunc("/api/usuariosevento", BajaUsuariosEvento).Methods("DELETE")
  routes.HandleFunc("/api/usuariosevento", ModificarUsuarioEvento).Methods("PUT")

  // Configuramos CORS
  cabecerasPermitidas := handlers.AllowedHeaders([]string{"Content-Type"})
  origenesPermitidos := handlers.AllowedOrigins([]string{"http://localhost:90"})
  metodosPermitidos := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})

  fmt.Println("Listening on port 8080...")
  http.Handle("/", routes)
  http.ListenAndServe(":8080", handlers.CORS(cabecerasPermitidas, origenesPermitidos, metodosPermitidos)(routes))
  //http.ListenAndServe(":8080", nil)

}
