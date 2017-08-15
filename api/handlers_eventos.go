package api

import(
  "net/http"
  "fmt"
  "strconv"
  "encoding/json"
  Error "rlnieto.org/eventos/go-services/error"
  Persistencia "rlnieto.org/eventos/go-services/persistencia"
)


/*------------------------------------------------------------------------------
 Alta de un evento

------------------------------------------------------------------------------*/
func AltaEvento(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  fecha := r.FormValue("fecha")
  hora := r.FormValue("hora")
  idLocal := r.FormValue("idlocal")
  motivo := r.FormValue("motivo")
  menu := r.FormValue("menu")
  precioEstimado := r.FormValue("precioestimado")
  idOrganizador := r.FormValue("idorganizador")

  // Validaciones
  // Campos obligatorios: fecha, hora, id organizador
  if fecha == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta la fecha"
    statusCode, response = error.Dispatch()
  }

  if hora == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta la hora"
    statusCode, response = error.Dispatch()
  }

  if idOrganizador == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta el organizador"
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }

  // Comprobamos los campos optativos => los numéricos vienen con un blanco
  if motivo == ""{
    motivo = "Porque hoy es hoy"
  }
  if idLocal == ""{
    idLocal = "0"
  }

  if precioEstimado == ""{
    precioEstimado = "0"
  }

  var evento Persistencia.Evento
  evento.Fecha = fecha
  evento.Hora = hora
  evento.IdLocal, _ = strconv.ParseInt(idLocal, 10, 32)
  evento.Motivo = motivo
  evento.Menu = menu
  evento.PrecioEstimado, _ = strconv.Atoi(precioEstimado)
  evento.IdOrganizador, _ = strconv.ParseInt(idOrganizador, 10, 32)

  dbError := evento.Alta()
  if dbError != ""{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  fmt.Fprintf(w, response)
}


/*------------------------------------------------------------------------------
 Modificación de un evento
 Esperamos que lleguen todos los datos!
 Si alguno de los campos optativos viene a blanco, lo actualizamos con un blanco

------------------------------------------------------------------------------*/
func ModificarEvento(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idEvento := r.FormValue("idevento")
  fecha := r.FormValue("fecha")
  hora := r.FormValue("hora")
  idLocal := r.FormValue("idlocal")
  motivo := r.FormValue("motivo")
  menu := r.FormValue("menu")
  precioEstimado := r.FormValue("precioestimado")
  idOrganizador := r.FormValue("idorganizador")

  // Validaciones
  // Campos obligatorios: fecha, hora, id organizador
  if idEvento == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta el id"
    statusCode, response = error.Dispatch()
  }

  if fecha == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta la fecha"
    statusCode, response = error.Dispatch()
  }

  if hora == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta la hora"
    statusCode, response = error.Dispatch()
  }

  if idOrganizador == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta el organizador"
    statusCode, response = error.Dispatch()
  }


  // Comprobamos los campos optativos => los numéricos vienen con un blanco
  if idLocal == ""{
    idLocal = "0"
  }

  if precioEstimado == ""{
    precioEstimado = "0"
  }

  if error.ErrorCode !=0{
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }

  var evento Persistencia.Evento
  evento.Id, _ = strconv.ParseInt(idEvento, 10, 32)
  evento.Fecha = fecha
  evento.Hora = hora
  evento.IdLocal, _ = strconv.ParseInt(idLocal, 10, 32)
  evento.Motivo = motivo
  evento.Menu = menu
  evento.PrecioEstimado, _ = strconv.Atoi(precioEstimado)
  evento.IdOrganizador, _ = strconv.ParseInt(idOrganizador, 10, 32)

  dbError := evento.Actualizar()
  if dbError != ""{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  fmt.Fprintf(w, response)
}


/*------------------------------------------------------------------------------
 Baja de un evento y sus usuarios asociados

------------------------------------------------------------------------------*/
func BorrarEvento(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idEvento := r.FormValue("idevento")

  // Validaciones
  if idEvento == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta el id"
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }

  var evento Persistencia.Evento
  evento.Id, _ = strconv.ParseInt(idEvento, 10, 32)

  dbError := evento.Borrar()
  if dbError != ""{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  fmt.Fprintf(w, response)
}


/*------------------------------------------------------------------------------
 Consulta de un evento con sus usuarios por clave

------------------------------------------------------------------------------*/
func ConsultarEvento(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idEvento := r.FormValue("idevento")

  // Validaciones
  if idEvento == ""{
    error.ErrorCode = Error.PARAMETRO_INCORRECTO
    error.Msg = "Falta el id"
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }


  // Datos a devolver. Se usa para transformar la respuesta a json, no es
  // necesario que sea global
  type datosEvento struct{
    Evento Persistencia.Evento
    Usuarios []Persistencia.UsuarioEvento
  }

  var salida datosEvento

  // Buscamos el evento por clave
  var evento Persistencia.Evento
  evento.Id, _ = strconv.ParseInt(idEvento, 10, 64)

  dbError := evento.ById()
  if dbError != ""{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }

  // Buscamos los usuarios
  var usuario Persistencia.UsuarioEvento

  usuarios, dbError := usuario.ByEvento(evento.Id)
  if dbError != ""{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }

  salida.Evento = evento
  salida.Usuarios = usuarios

  output, _ := json.Marshal(salida)
  fmt.Fprintf(w, string(output))
}
