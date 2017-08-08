package api

import(
  _ "github.com/go-sql-driver/mysql"
  "net/http"
  "fmt"
  "strings"
  "log"
  Error "rlnieto.org/go-services/error"
)


/*------------------------------------------------------------------------------
 Alta de usuarios en la cena
 OJO: espera la lista de TODOS los usuarios y BORRA los que ya hay!!

------------------------------------------------------------------------------*/
func AltaUsuariosCena(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")
  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idUsuarios := strings.Split(r.FormValue("idusuario"), ",")
  idCena := r.FormValue("idcena")
  confirmado := r.FormValue("confirmado")

  // Validaciones
  // Campos obligatorios: fecha, hora, id organizador
  if r.FormValue("idusuario") == ""{
    error.ErrorCode = Error.GENERIC_ERROR
    error.Msg = "Faltan los ids de los usuarios"
    statusCode, response = error.Dispatch()
  }

  if idCena == ""{
    error.ErrorCode = Error.GENERIC_ERROR
    error.Msg = "Falta el id de la cena"
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

  // Iniciamos la transacción
  tx, dbError := Database.Begin()
  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  // Borramos los usuarios que ya hay asociados
  sql := "DELETE FROM usuarios_cena where idcena=?"

  query, dbError := Database.Prepare(sql)
  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  _, dbError = tx.Stmt(query).Exec(idCena)
  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)

    tx.Rollback()
    fmt.Fprintf(w, response)
    return
  }

  // Damos de alta los usuarios enviados
  sql = "INSERT INTO usuarios_cena set idusuario=?" +
  ", idcena=" + idCena + ", confirmado='" + confirmado + "'"

  queryAlta, dbError := Database.Prepare(sql)
  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  // Insertamoss las filas en la tabla
  for _, idUsuario := range(idUsuarios){
    res, dbError := tx.Stmt(queryAlta).Exec(idUsuario)
    fmt.Println(res)

    // Si hubo error hacemos rollback
    if dbError != nil{
      error.ErrorCode = Error.DB_ERROR
      error.Msg = dbError.Error()
      statusCode, response = error.Dispatch()
      http.Error(w, http.StatusText(int(statusCode)), statusCode)

      tx.Rollback()
      break
    }
  }
  if error.ErrorCode == 0{
    tx.Commit()
  }
  fmt.Fprintf(w, response)
}


/*------------------------------------------------------------------------------
 Baja de los usuarios asociados a la cena

------------------------------------------------------------------------------*/
func BajaUsuariosCena(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idCena := r.FormValue("idcena")
  idUsuarios := r.FormValue("idusuario")

  // Validaciones
  if idCena == ""{
    error.ErrorCode = Error.GENERIC_ERROR
    error.Msg = "Falta el id"
    statusCode, response = error.Dispatch()
  }

  if idUsuarios == ""{
    error.ErrorCode = Error.GENERIC_ERROR
    error.Msg = "Faltan los ids de los usuarios"
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }


  // Borramos los asistentes que nos han enviado
  sql := "DELETE FROM usuarios_cena WHERE idcena = " + idCena  +
    " and idusuario in(" + idUsuarios + ")"

  log.Println(sql)

  _, dbError := Database.Exec(sql)

  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  fmt.Fprintf(w, response)
}
