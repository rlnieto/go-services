package main

import(
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"
  "encoding/xml"
  "fmt"
  "log"
  "strconv"
  "strings"
  "rlnieto.org/cenasapp/go-services/specifications"
)

var database *sql.DB
var Format string

type Users struct{
  Users []User `json:"users"`
}

type User struct{
  ID int "json:id"
  Name string "json:username"
  Email string "json:email"
  First string "json:first"
  Last string "json:last"
}


type UpdateResponse struct {
  Error string "json:error"
  ErrorCode int "json:code"
}


// Para los errores
type CreateResponse struct {
  Error string "json:error"
  ErrorCode int "json:code"
}

type ErrMsg struct{
    ErrCode int
    StatusCode int
    Msg string
}

// Recupera el formato de salida de los datos de los servicios (json, rss, txt...)
// Si no lo envían asumimos json
func GetFormat(r *http.Request){
  if r.URL.Query()["format"] != nil{
    Format = r.URL.Query()["format"][0]
  }else{
    Format = "json"
  }
}

// Convierte los datos de salida al formato pedido
func SetFormat(data interface{}) []byte{
  var apiOutput []byte

  if Format == "json" || Format == ""{
    output, _ := json.Marshal(data)
    apiOutput = output
    log.Println("Json marshalling")
  }else if Format == "xml"{
    output, _ := xml.Marshal(data)
    apiOutput = output
  }

  return apiOutput
}

/*------------------------------------------------------------------------------
 Gestión de errores

------------------------------------------------------------------------------*/
func dbErrorParse(err string) (string, int64){
  Parts := strings.Split(err, ":")
  errorMessage := Parts[1]
  Code := strings.Split(Parts[0], " ")
  errorCode, _ := strconv.ParseInt(Code[1], 10, 32)

  return errorMessage, errorCode
}

func ErrorMessages(err int64) (ErrMsg){
  var em ErrMsg

  errorMessage := ""
  statusCode := 200;
  errorCode := 0

  switch(err){
    case 1062:
      errorMessage = "Duplicate entry"
      errorCode = 10
      statusCode = 409
    default:
      errorMessage = http.StatusText(int(err))
      errorCode = 0
      statusCode = int(err)
  }


  em.ErrCode = errorCode
  em.StatusCode = statusCode
  em.Msg = errorMessage

  return em
}


/*------------------------------------------------------------------------------
 Consulta de n usuarios

------------------------------------------------------------------------------*/
func UsersRetrieve(w http.ResponseWriter, r *http.Request){
  GetFormat(r)   // el formato de salida viene indicado en un parámetro

  start := 0
  limit := 10
  next := start + limit;

  w.Header().Set("Pragma", "no-cache")
  // Enviamos un link en la cabecera con la llamada para consultar los 10 siguientes
  w.Header().Set("Link", "http://localhost:8080/api/users?start=" + strconv.Itoa(next) + ";rel=\"next\"")

  rows, _ := database.Query("select * from users LIMIT 10")
  Response := Users{}

  for rows.Next(){
    user := User{}
    rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)

    Response.Users = append(Response.Users, user)
  }

  //output, _ := json.Marshal(Response)
  output := SetFormat(Response)

  fmt.Println(string(output))

}


/*------------------------------------------------------------------------------
 Consulta de usuario por clave

------------------------------------------------------------------------------*/
func UserRetrieve(w http.ResponseWriter, r *http.Request){
  urlParams := mux.Vars(r)
  id := urlParams["id"]
  ReadUser := User{}

  w.Header().Set("Pragma", "no-cache")

  err := database.QueryRow("select * from users where user_id=?", id).Scan(&ReadUser.ID, &ReadUser.Name, &ReadUser.First, &ReadUser.Last, &ReadUser.Email)
  switch{
  case err == sql.ErrNoRows:
    fmt.Fprintf(w, "No such user")
  case err != nil:
    log.Fatal(err)
    fmt.Fprintf(w, "Error en la consulta")
  default:
    output, _ := json.Marshal(ReadUser)
    fmt.Fprintf(w, string(output))
  }
}

/*------------------------------------------------------------------------------
 Alta de usuario

------------------------------------------------------------------------------*/
func UserCreate(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  NewUser := User{}

  NewUser.Name = r.FormValue("name")
  NewUser.Email = r.FormValue("email")
  NewUser.First = r.FormValue("first")
  NewUser.Last = r.FormValue("last")

  _, err := json.Marshal(NewUser)
  if err != nil{
    fmt.Println("Error al generar el json")
  }

  Response := CreateResponse{}

  sql := "INSERT INTO users set user_nickname='" + NewUser.Name +
  "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last +
  "', user_email='" + NewUser.Email + "'"

  q, err := database.Exec(sql)

  if err != nil{
    errorMessage, errorCode := dbErrorParse(err.Error())
    fmt.Println(errorMessage)

    error := ErrorMessages(errorCode)
    Response.Error = error.Msg
    Response.ErrorCode = error.ErrCode
    http.Error(w, "Conflict", error.StatusCode)
  }

  fmt.Println(q)
  createOutput, _ := json.Marshal(Response)
  fmt.Fprintf(w, string(createOutput))
}

/*------------------------------------------------------------------------------
 Actualización de usuario

------------------------------------------------------------------------------*/
func UsersUpdate(w http.ResponseWriter, r *http.Request){
  GetFormat(r)

  Response := UpdateResponse{}
  params := mux.Vars(r)
  uid := params["id"]
  email := r.FormValue("email")

  error :=  ErrMsg{}

  var userCount int
  err := database.QueryRow("SELECT COUNT(user_id) FROM users WHERE user_id = ?", uid).Scan(&userCount)

  if userCount == 0 {
    error = ErrorMessages(404)
    log.Println(error.ErrCode)
    log.Println(w, error.Msg, error.StatusCode)

    Response.Error = error.Msg
    Response.ErrorCode = error.ErrCode

    http.Error(w, error.Msg, error.StatusCode)
  }else if err != nil {
    log.Println(error.Msg)
  }else{
    _, uperr := database.Exec("UPDATE users SET user_email=? WHERE user_id = ?", email, uid)
    if uperr != nil{
      _, errorCode := dbErrorParse(uperr.Error())
      error = ErrorMessages(errorCode)
      Response.Error = error.Msg
      Response.ErrorCode = error.ErrCode
      http.Error(w, error.Msg, error.StatusCode)
    }else{
      Response.Error = "success"
      Response.ErrorCode = 0
      output := SetFormat(Response)
      fmt.Fprintln(w, string(output))
    }
  }

}

/*------------------------------------------------------------------------------
 Información acerca de los endpoints

------------------------------------------------------------------------------*/
type DocMethod interface{
}

func UsersInfo(w http.ResponseWriter, r *http.Request){
  GetFormat(r)

  w.Header().Set("Allow", "DELETE, GET, HEAD, OPTIONS, POST, PUT")

  UserDocumentation := []DocMethod{}

  UserDocumentation = append(UserDocumentation, specifications.UserGET)
  UserDocumentation = append(UserDocumentation, specifications.UserPOST)
  UserDocumentation = append(UserDocumentation, specifications.UserOPTIONS)

  output := SetFormat(UserDocumentation)

  fmt.Fprintln(w, string(output))
}


/*------------------------------------------------------------------------------
 main


------------------------------------------------------------------------------*/
func main(){
  db, err := sql.Open("mysql", "root:kalimotxo@/golang")
  if err != nil {
    fmt.Println("Error conectando con la bd")
  }

  database = db

  routes := mux.NewRouter()
  routes.HandleFunc("/api/users", UserCreate).Methods("POST")
  routes.HandleFunc("/api/users", UsersRetrieve).Methods("GET")
  routes.HandleFunc("/api/users/{id:[0-9]+}", UserRetrieve).Methods("GET")
  routes.HandleFunc("/api/users/{id:[0-9]+}", UsersUpdate).Methods("PUT")
  routes.HandleFunc("/api/users", UsersInfo).Methods("OPTIONS")

  http.Handle("/", routes)
  http.ListenAndServe(":8080", nil)
}
