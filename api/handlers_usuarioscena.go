package api

import(
  "net/http"
  "fmt"
  "strings"
  "strconv"
  "log"
  Mensaje "rlnieto.org/go-services/mensajes"
  Persistencia "rlnieto.org/go-services/persistencia"
)

/*------------------------------------------------------------------------------
 Alta de usuarios en un evento

------------------------------------------------------------------------------*/
func AltaUsuariosEvento(w http.ResponseWriter, r *http.Request){

  var error = Mensaje.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idUsuarios := strings.Split(r.FormValue("idusuario"), ",")
  idEvento := r.FormValue("idevento")
  confirmado := r.FormValue("confirmado")

  // Validaciones
  // Campos obligatorios: fecha, hora, id organizador
  if r.FormValue("idusuario") == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_IDS_USUARIOS
    statusCode, response = error.Dispatch()
  }

  if idEvento == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_ID_EVENTO
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, response, statusCode)
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
    error.ErrorCode = Mensaje.NOT_FOUND
    error.Msg = Mensaje.ERROR_EVENTO_NO_ENCONTRADO + idEvento
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  // Comprobamos que existan los usuarios
  var usuario Persistencia.Usuario

  for _, id := range(idUsuarios){
    usuario.Id, _ = strconv.ParseInt(id, 10, 32);
    dbError = usuario.ById()

    if dbError != ""{
      error.ErrorCode = Mensaje.NOT_FOUND
      error.Msg = dbError
      statusCode, response = error.Dispatch()
      http.Error(w, response, statusCode)
      return
    }
    fmt.Println(usuario.Nick)
  }


  // Damos de alta los usuarios
  if dbError == ""{
    var usuario Persistencia.UsuarioEvento

    dbError = usuario.AltaEnEvento(idEventoInt, idUsuarios)
    if dbError != ""{
      error.ErrorCode = Mensaje.DB_ERROR
      error.Msg = dbError
      statusCode, response = error.Dispatch()
      http.Error(w, response, statusCode)
      return
    }
  }

  fmt.Fprintf(w, response)
}


/*------------------------------------------------------------------------------
 Baja de los usuarios asociados al evento

------------------------------------------------------------------------------*/
func BajaUsuariosEvento(w http.ResponseWriter, r *http.Request){

  var error = Mensaje.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idEvento := r.FormValue("idevento")
  idUsuarios := r.FormValue("idusuario")

  // Validaciones
  if idEvento == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_ID
    statusCode, response = error.Dispatch()
  }

  if idUsuarios == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_IDS_USUARIOS
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, response, statusCode)
    return
  }

  // Comprobamos que haya usuarios para el evento
  var usuario Persistencia.UsuarioEvento
  idEventoInt, _ := strconv.ParseInt(idEvento, 10, 32)

  numUsuarios, dbError := usuario.NumeroUsuariosEvento(idEventoInt)
  if dbError != ""{
    error.ErrorCode = Mensaje.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  log.Println(numUsuarios)

  if numUsuarios == 0{
    error.ErrorCode = Mensaje.NOT_FOUND
    error.Msg = Mensaje.ERROR_NO_HAY_USUARIOS_EVENTO + idEvento
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  // Hacemos la baja
  dbError = usuario.BajaEnEvento(idEventoInt, idUsuarios)

  if dbError != ""{
    error.ErrorCode = Mensaje.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  fmt.Fprintf(w, response)
}


/*------------------------------------------------------------------------------
 Cambia el estado de la asistencia de un usuario a un evento

------------------------------------------------------------------------------*/
func ModificarUsuarioEvento(w http.ResponseWriter, r *http.Request){

  var error = Mensaje.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idEvento := r.FormValue("idevento")
  idUsuario := r.FormValue("idusuario")
  confirmado := r.FormValue("confirmado")

  // Validaciones
  if idEvento == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_ID_EVENTO
    statusCode, response = error.Dispatch()
  }

  if idUsuario == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_ID_USUARIO
    statusCode, response = error.Dispatch()
  }

  if confirmado == "" && confirmado != "S" && confirmado != "N"{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_IND_ASISTENCIA
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, response, statusCode)
    return
  }


  // Comprobamos que exista entrada en la tabla para ese usuario y evento
  var usuario Persistencia.UsuarioEvento

  idUsuarioInt, _ := strconv.ParseInt(idUsuario, 10, 32)
  idEventoInt, _ := strconv.ParseInt(idEvento, 10, 32)

  usuario.IdUsuario = idUsuarioInt
  usuario.IdEvento = idEventoInt

  dbError := usuario.ById()
  if dbError != ""{
    error.ErrorCode = Mensaje.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  // Si no había datos tambien salimos con error
  if(usuario.IdUsuario == 0){
    error.ErrorCode = Mensaje.NOT_FOUND
    error.Msg = Mensaje.ERROR_NO_HAY_DATOS_GENERICO
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  // Modificamos el indicador con el valor recibido
  usuario.Confirmado = confirmado
  dbError = usuario.Actualizar()

  if dbError != ""{
    error.ErrorCode = Mensaje.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  fmt.Fprintf(w, response)
}
