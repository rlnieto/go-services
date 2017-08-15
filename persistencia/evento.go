package persistencia

import(
  "strconv"
)

type Evento struct{
  Id int64
  Fecha string
  Hora string
  IdLocal int64
  Motivo  string
  Menu  string
  PrecioEstimado int
  IdOrganizador int64
}

//-----------------------------------------------------------------------------
// Baja
//
//-----------------------------------------------------------------------------
func (evento *Evento) Borrar() (string){
  mensajeError := ""

  Db.Open()
  defer Db.Close()

  conexion := Db.Conn

  // Iniciamos la transacción
  tx, dbError := conexion.Begin()
  if dbError != nil{
    mensajeError = dbError.Error()
    return mensajeError
  }

  // Borramos los asistentes
  sql := "DELETE FROM usuarios_evento WHERE idevento = ?"

  queryBaja, dbError := conexion.Prepare(sql)
  if dbError != nil{
    mensajeError = dbError.Error()
    return mensajeError
  }

  _, dbError = tx.Stmt(queryBaja).Exec(evento.Id)
  if dbError != nil{
    mensajeError = dbError.Error()
    tx.Rollback()
    return mensajeError
  }

  // Borramos el evento
  sql = "DELETE FROM evento WHERE id = ?"

  queryBaja, dbError = conexion.Prepare(sql)
  if dbError != nil{
    mensajeError = dbError.Error()
    tx.Rollback()
    return mensajeError
  }

  _, dbError = tx.Stmt(queryBaja).Exec(evento.Id)
  if dbError != nil{
    mensajeError = dbError.Error()
    tx.Rollback()
    return mensajeError
  }

  tx.Commit()
  return mensajeError

}




//-----------------------------------------------------------------------------
// Actualización
//
//-----------------------------------------------------------------------------
func (evento *Evento) Actualizar() (string){
  mensajeError := ""

  Db.Open()
  defer Db.Close()

  statement, dbError := Db.Conn.Prepare("UPDATE evento SET fecha=?, hora=?, idlocal=?, motivo=?, menu=?, precio_estimado=?, idorganizador=? WHERE id=?")
	if dbError != nil {
    mensajeError = dbError.Error()
    return mensajeError
  }

  result, dbError := statement.Exec(evento.Fecha,evento.Hora,evento.IdLocal,evento.Motivo,evento.Menu,evento.PrecioEstimado,evento.IdOrganizador,evento.Id)
  if dbError != nil {
    mensajeError = dbError.Error()
    return mensajeError
  }
  filas, _ := result.RowsAffected()
  if filas == 0{
    return "No se encontraron filas"
  }

  return mensajeError
}


//-----------------------------------------------------------------------------
// Alta
//
//-----------------------------------------------------------------------------
func (evento *Evento) Alta() (string){
  mensajeError := ""

  Db.Open()
  defer Db.Close()

  statement, dbError := Db.Conn.Prepare("INSERT INTO evento VALUES( ?, ?, ?, ?, ?, ?, ?, ? )")
	if dbError != nil {
    mensajeError = dbError.Error()
    return mensajeError
  }

  _, dbError = statement.Exec(0,evento.Fecha,evento.Hora,evento.IdLocal,evento.Motivo,evento.Menu,evento.PrecioEstimado,evento.IdOrganizador)
  if dbError != nil {
    mensajeError = dbError.Error()
    return mensajeError
  }

  return mensajeError
}



//-----------------------------------------------------------------------------
// Búsqueda por clave
//
//-----------------------------------------------------------------------------
func (evento *Evento) ById() (string){
  mensajeError := ""

  Db.Open()
  defer Db.Close()

  query := "SELECT * FROM evento WHERE id=" + strconv.Itoa(int(evento.Id))
  dbError := Db.Conn.QueryRow(query).Scan(&evento.Id,&evento.Fecha,&evento.Hora,&evento.IdLocal,&evento.Motivo,&evento.Menu,&evento.PrecioEstimado,&evento.IdOrganizador)

  // Cuando no encuentra filas devuelve el mensaje de error: "sql: no rows in result set"
  if dbError != nil{
    mensajeError = dbError.Error()
  }

  return mensajeError
}
