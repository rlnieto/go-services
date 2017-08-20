package persistencia

import(
  "fmt"
  "strconv"
  "github.com/go-pg/pg"
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

  // Borramos el evento
  sql := "DELETE FROM eventos WHERE id = $1::integer"

  queryBaja, dbError := conexion.Prepare(sql)
  if dbError != nil{
    mensajeError = dbError.Error()
    tx.Rollback()
    return mensajeError
  }

  res, dbError := tx.Stmt(queryBaja).Exec(evento.Id)
  if dbError != nil{
    mensajeError = dbError.Error()
    tx.Rollback()
    return mensajeError
  }

  // Si no existía el evento salimos con error
  if res.RowsAffected() == 0{
    tx.Rollback()
    return "No se encontraron datos para el evento " + strconv.FormatInt(evento.Id, 10)
  }



  // Borramos los asistentes
  sql = "DELETE FROM usuario_eventos WHERE id_evento = $1::integer"

  queryBaja, dbError = conexion.Prepare(sql)
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

  dbError := Db.Conn.Update(evento)
	if dbError != nil {
    if dbError == pg.ErrNoRows{
      mensajeError = "No hay datos para el evento " + strconv.FormatInt(evento.Id, 10)
    }else{
      mensajeError = dbError.Error()
    }
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

  dbError := Db.Conn.Insert(evento)
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

  query := "SELECT * FROM eventos WHERE id=" + strconv.Itoa(int(evento.Id))
  fmt.Println(query)
  _, dbError := Db.Conn.QueryOne(evento, query)

  if dbError != nil{
    if dbError == pg.ErrNoRows{
      mensajeError = "No hay datos para el evento " + strconv.FormatInt(evento.Id, 10)
    }else{
      mensajeError = dbError.Error()
    }
  }

  return mensajeError
}
