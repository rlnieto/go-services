package api

import(
  "net/http"
  "fmt"
  "strings"
  "strconv"
  "log"
  Error "rlnieto.org/eventos/go-services/error"
  Persistencia "rlnieto.org/eventos/go-services/persistencia"
)


/*------------------------------------------------------------------------------
 Alta de usuarios en un evento

------------------------------------------------------------------------------*/
func AltaUsuariosEvento(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")
  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idUsuarios := strings.Split(r.FormValue("idusuario"), ",")
  idEvento := r.FormValue("idevento")
  confirmado := r.FormValue("confirmado")

  // Validaciones
  // Campos obligatorios: fecha, hora, id organizador
  if r.FormValue("idusuario") == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Faltan los ids de los usuarios"
    statusCode, response = error.Dispatch()
  }

  if idEvento == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta el id de la evento"
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }

  // Comprobamos los campos optativos
  if confirmado == ""{
    confirmado = "N"
  }

  // Comprobamos que exista el evento
  var evento Persistencia.Evento
  idEventoInt, _ := strconv.ParseInt(idEvento, 10, 32)
  evento.Id = idEventoInt

  dbError := evento.ById()
  if dbError != ""{
    error.ErrorCode = Error.NOT_FOUND
    error.Msg = "No se encontró el evento con id " + idEvento
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  // Damos de alta los usuarios
  if dbError == ""{
    var usuario Persistencia.UsuarioEvento

    dbError = usuario.AltaEnEvento(idEventoInt, idUsuarios)
    if dbError != ""{
      error.ErrorCode = Error.DB_ERROR
      error.Msg = dbError
      statusCode, response = error.Dispatch()
      http.Error(w, http.StatusText(int(statusCode)), statusCode)
    }
  }

  fmt.Fprintf(w, response)
}


/*------------------------------------------------------------------------------
 Baja de los usuarios asociados al evento

------------------------------------------------------------------------------*/
func BajaUsuariosEvento(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idEvento := r.FormValue("idevento")
  idUsuarios := r.FormValue("idusuario")

  // Validaciones
  if idEvento == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta el id"
    statusCode, response = error.Dispatch()
  }

  if idUsuarios == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Faltan los ids de los usuarios"
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }

  // Comprobamos que haya usuarios para el evento
  var usuario Persistencia.UsuarioEvento
  idEventoInt, _ := strconv.ParseInt(idEvento, 10, 32)

  numUsuarios, dbError := usuario.NumeroUsuariosEvento(idEventoInt)
  if dbError != ""{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }

  log.Println(numUsuarios)

  if numUsuarios == 0{
    error.ErrorCode = Error.NOT_FOUND
    error.Msg = "No hay usuarios para el evento con id " + idEvento
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }

  // Hacemos la baja
  dbError = usuario.BajaEnEvento(idEventoInt, idUsuarios)

  if dbError != ""{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  fmt.Fprintf(w, response)
}
