package api

import(
  _ "github.com/go-sql-driver/mysql"
  "net/http"
  "fmt"
  "strconv"
  "encoding/json"
  Error "rlnieto.org/go-services/error"
  Persistencia "rlnieto.org/go-services/persistencia"
)


/*------------------------------------------------------------------------------
 Alta de una cena

------------------------------------------------------------------------------*/
func AltaCena(w http.ResponseWriter, r *http.Request){
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
    error.ErrorCode = Error.GENERIC_ERROR
    error.Msg = "Falta la fecha"
    statusCode, response = error.Dispatch()
  }

  if hora == ""{
    error.ErrorCode = Error.GENERIC_ERROR
    error.Msg = "Falta la hora"
    statusCode, response = error.Dispatch()
  }

  if idOrganizador == ""{
    error.ErrorCode = Error.GENERIC_ERROR
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

  // Insertamos la fila en la tabla
  sql := "INSERT INTO cena set fecha='" + fecha +
  "', hora='" + hora + "', idlocal=" + idLocal +
  ", motivo='" + motivo +
  "', menu='" + menu + "', precio_estimado=" + precioEstimado +
  ", idorganizador=" + idOrganizador

  _, dbError := Database.Exec(sql)

  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  fmt.Fprintf(w, response)
}


/*------------------------------------------------------------------------------
 Modificación de una cena
 Esperamos que lleguen todos los datos!
 Si alguno de los campos optativos viene a blanco, lo actualizamos con un blanco

------------------------------------------------------------------------------*/
func ModificarCena(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idCena := r.FormValue("idcena")
  fecha := r.FormValue("fecha")
  hora := r.FormValue("hora")
  idLocal := r.FormValue("idlocal")
  motivo := r.FormValue("motivo")
  menu := r.FormValue("menu")
  precioEstimado := r.FormValue("precioestimado")
  idOrganizador := r.FormValue("idorganizador")

  // Validaciones
  // Campos obligatorios: fecha, hora, id organizador
  if idCena == ""{
    error.ErrorCode = Error.GENERIC_ERROR
    error.Msg = "Falta el id"
    statusCode, response = error.Dispatch()
  }

  if fecha == ""{
    error.ErrorCode = Error.GENERIC_ERROR
    error.Msg = "Falta la fecha"
    statusCode, response = error.Dispatch()
  }

  if hora == ""{
    error.ErrorCode = Error.GENERIC_ERROR
    error.Msg = "Falta la hora"
    statusCode, response = error.Dispatch()
  }

  if idOrganizador == ""{
    error.ErrorCode = Error.GENERIC_ERROR
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

  // Actualizamos la fila
  sql := "UPDATE cena set fecha='" + fecha +
  "', hora='" + hora + "', idlocal=" + idLocal +
  ", motivo='" + motivo +
  "', menu='" + menu + "', precio_estimado=" + precioEstimado +
  ", idorganizador=" + idOrganizador +
  " WHERE id = " + idCena

  _, dbError := Database.Exec(sql)

  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  fmt.Fprintf(w, response)
}


/*------------------------------------------------------------------------------
 Baja de una cena y sus usuarios asociados

------------------------------------------------------------------------------*/
func BorrarCena(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idCena := r.FormValue("idcena")

  // Validaciones
  if idCena == ""{
    error.ErrorCode = Error.GENERIC_ERROR
    error.Msg = "Falta el id"
    statusCode, response = error.Dispatch()
  }

  if error.ErrorCode !=0{
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }

  tx, dbError := Database.Begin()
  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  // Borramos los asistentes
  sql := "DELETE FROM usuarios_cena WHERE idcena = ?"

  queryBaja, dbError := Database.Prepare(sql)
  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  _, dbError = tx.Stmt(queryBaja).Exec(idCena)
  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)

    tx.Rollback()
  }

  // Borramos la cena
  sql = "DELETE FROM cena WHERE id = ?"

  queryBaja, dbError = Database.Prepare(sql)
  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
  }

  _, dbError = tx.Stmt(queryBaja).Exec(idCena)
  if dbError != nil{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError.Error()
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)

    tx.Rollback()
  }

  if error.ErrorCode == 0{
    tx.Commit()
  }
  fmt.Fprintf(w, response)

}


/*------------------------------------------------------------------------------
 Consulta de una cena con sus usuarios por clave

------------------------------------------------------------------------------*/
func ConsultarCena(w http.ResponseWriter, r *http.Request){
  //w.Header().Set("Access-Control-Allow-Origin","http://localhost:90")

  var error = Error.ErrorMsg{}

  // Inicializamos la respuesta asumiendo que todo fue ok
  response := error.OkResponse()
  var statusCode int

  idCena := r.FormValue("idcena")

  // Validaciones
  if idCena == ""{
    error.ErrorCode = Error.GENERIC_ERROR
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
  type datosCena struct{
    Cena Persistencia.Cena
    Usuarios []Persistencia.UsuarioCena
  }

  var salida datosCena

  // Buscamos la cena por clave
  var cena Persistencia.Cena
  cena.Id, _ = strconv.ParseInt(idCena, 10, 32)

  dbError := cena.ById()
  if dbError != ""{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }


  // Buscamos los usuarios
  var usuario Persistencia.UsuarioCena

  usuarios, dbError := usuario.ByCena(cena.Id)
  if dbError != ""{
    error.ErrorCode = Error.DB_ERROR
    error.Msg = dbError
    statusCode, response = error.Dispatch()
    http.Error(w, http.StatusText(int(statusCode)), statusCode)
    fmt.Fprintf(w, response)
    return
  }


  salida.Cena = cena
  salida.Usuarios = usuarios

  output, _ := json.Marshal(salida)
  fmt.Fprintf(w, string(output))
}
