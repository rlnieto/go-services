package api

import(
  "net/http"
  "fmt"
  "strconv"
  "encoding/json"
  Mensaje "rlnieto.org/go-services/mensajes"
  Persistencia "rlnieto.org/go-services/persistencia"
)

/*------------------------------------------------------------------------------
 Alta de un evento

------------------------------------------------------------------------------*/
func AltaEvento(w http.ResponseWriter, r *http.Request){

  var error = Mensaje.ErrorMsg{}

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
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_FECHA
    statusCode, response = error.Dispatch()
  }

  if hora == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_HORA
    statusCode, response = error.Dispatch()
  }

  if idOrganizador == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_ORGANIZADOR
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, response, statusCode)
    return
  }

  // Comprobamos los campos optativos => los numéricos vienen con un blanco
  if motivo == ""{
    motivo = Mensaje.MOTIVO_EVENTO_POR_DEFECTO
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

  idLocalConvertido, errorConversion := strconv.ParseInt(idLocal, 10, 32)
  if errorConversion != nil{
    error.ErrorCode = Mensaje.DB_ERROR
    error.Msg = Mensaje.ERROR_FORMATO_ID_LOCAL
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  evento.IdLocal = idLocalConvertido
  evento.Motivo = motivo
  evento.Menu = menu
  evento.PrecioEstimado, _ = strconv.Atoi(precioEstimado)
  evento.IdOrganizador, _ = strconv.ParseInt(idOrganizador, 10, 32)

  dbError := evento.Alta()
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
 Modificación de un evento
 Esperamos que lleguen todos los datos!
 Si alguno de los campos optativos viene a blanco, lo actualizamos con un blanco

------------------------------------------------------------------------------*/
func ModificarEvento(w http.ResponseWriter, r *http.Request){

  var error = Mensaje.ErrorMsg{}

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
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_ID
    statusCode, response = error.Dispatch()
  }

  if fecha == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_FECHA
    statusCode, response = error.Dispatch()
  }

  if hora == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_HORA
    statusCode, response = error.Dispatch()
  }

  if idOrganizador == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_ORGANIZADOR
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
    http.Error(w, response, statusCode)
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
    error.ErrorCode = Mensaje.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  fmt.Fprintf(w, response)
}


/*------------------------------------------------------------------------------
 Baja de un evento y sus usuarios asociados

------------------------------------------------------------------------------*/
func BorrarEvento(w http.ResponseWriter, r *http.Request){

  var error = Mensaje.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idEvento := r.FormValue("idevento")

  // Validaciones
  if idEvento == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_ID
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, response, statusCode)
    return
  }

  var evento Persistencia.Evento
  evento.Id, _ = strconv.ParseInt(idEvento, 10, 32)

  dbError := evento.Borrar()
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
 Consulta de un evento con sus usuarios por clave de evento

------------------------------------------------------------------------------*/
func ConsultarEvento(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  var error = Mensaje.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idEvento := r.FormValue("idevento")

  // Validaciones
  if idEvento == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_ID
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, response, statusCode)
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
    error.ErrorCode = Mensaje.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  // Buscamos los usuarios
  var usuario Persistencia.UsuarioEvento

  usuarios, dbError := usuario.ByEvento(evento.Id)
  if dbError != ""{
    error.ErrorCode = Mensaje.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  salida.Evento = evento
  salida.Usuarios = usuarios

  output, _ := json.Marshal(salida)
  fmt.Fprintf(w, string(output))
  fmt.Fprintf(w, response)
}




/*------------------------------------------------------------------------------
 Consulta de los eventos de un usuario

------------------------------------------------------------------------------*/
func EventosUsuario(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  type DatosSalida struct{
    IdEvento int64
    Fecha string
    Hora string
    Local string
    Organizador string
    Motivo string
    Confirmado string
  }

  var error = Mensaje.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idUsuario := r.FormValue("idusuario")

  // Validaciones
  if idUsuario == ""{
    error.ErrorCode = Mensaje.PARAMETRO_INCORRECTO
    error.Msg = Mensaje.ERROR_NO_HAY_ID
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, response, statusCode)
    return
  }

  fmt.Println("Recuperando los eventos del usuario...")
  fmt.Println(idUsuario)


  // Buscamos los eventos
  var salida []Persistencia.Evento
  var evento Persistencia.Evento

  usuarioInt, _ := strconv.ParseInt(idUsuario, 10, 64)

  salida, dbError := evento.ByUsuario(usuarioInt)
  if dbError != ""{
    error.ErrorCode = Mensaje.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  // Comprobamos si había datos
  if len(salida) == 0{
    error.ErrorCode = Mensaje.DB_ERROR
    error.Msg = Mensaje.NO_HAY_EVENTOS_USUARIO + idUsuario
    statusCode, response = error.Dispatch()
    http.Error(w, response, statusCode)
    return
  }

  datosSalida := []DatosSalida{}
  var item DatosSalida

  for _, evento := range salida{
    item.IdEvento = evento.Id
    item.Fecha = evento.Fecha
    item.Hora = evento.Hora

    // Recuperamos el nombre del local
    var local Persistencia.Restaurante
    local.Id = evento.IdLocal
    dbError = local.ById()
    if dbError != ""{
      error.ErrorCode = Mensaje.DB_ERROR
      error.Msg = dbError
      statusCode, response = error.Dispatch()
      http.Error(w, response, statusCode)
      return
    }
    if local.Nombre == "" {
      item.Local = Mensaje.DESCRIPCION_LOCAL_POR_DEFECTO
    }else{
      item.Local = local.Nombre
    }


    // Recuperamos el nombre del organizador
    var usuario Persistencia.Usuario
    usuario.Id = evento.IdOrganizador
    dbError = usuario.ById()
    if dbError != ""{
      error.ErrorCode = Mensaje.DB_ERROR
      error.Msg = dbError
      statusCode, response = error.Dispatch()
      http.Error(w, response, statusCode)
      return
    }
    item.Organizador = usuario.Nick

    item.Motivo = evento.Motivo

    // Y comprobamos si confirmó o no la asistencia
    var usuarioEvento Persistencia.UsuarioEvento
    usuarioEvento.IdUsuario = usuarioInt
    usuarioEvento.IdEvento = evento.Id
    dbError = usuarioEvento.ById()
    if dbError != ""{
      error.ErrorCode = Mensaje.DB_ERROR
      error.Msg = dbError
      statusCode, response = error.Dispatch()
      http.Error(w, response, statusCode)
      return
    }
    item.Confirmado = usuarioEvento.Confirmado

    datosSalida = append(datosSalida, item)
  }

  output, _ := json.Marshal(datosSalida)
  fmt.Fprintf(w, string(output))
}
