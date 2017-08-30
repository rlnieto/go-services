package persistencia

import(
  "fmt"
  "strconv"
  "github.com/go-pg/pg"
  Mensaje "rlnieto.org/go-services/mensajes"
)

type Evento struct{
  Id int64
  Fecha string
  Hora string
  IdLocal int64  `sql:",notnull"`
  Motivo  string
  Menu  string
  PrecioEstimado int  `sql:",notnull"`
  IdOrganizador int64  `sql:",notnull"`
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
    return Mensaje.NO_HAY_DATOS_EVENTO + strconv.FormatInt(evento.Id, 10)
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
      mensajeError = Mensaje.NO_HAY_DATOS_EVENTO + strconv.FormatInt(evento.Id, 10)
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

  // Comprobamos si existe el local. Un 0 siginifica que no lo hay
  if evento.IdLocal  > 0{
    var local Restaurante
    local.Id = evento.IdLocal

    dbError := Db.Conn.Select(&local)
    if dbError != nil{
      if dbError == pg.ErrNoRows{
        return Mensaje.NO_HAY_DATOS_LOCAL + strconv.FormatInt(evento.IdLocal, 10)
      }else{
        return dbError.Error()
      }
    }
  }

  // Comprobamos si existe el organizador
  var organizador Usuario
  organizador.Id = evento.IdOrganizador

  dbError := Db.Conn.Select(&organizador)
  if dbError != nil{
    if dbError == pg.ErrNoRows{
      return Mensaje.NO_HAY_ORGANIZADOR + strconv.FormatInt(evento.IdOrganizador, 10)
    }else{
      return dbError.Error()
    }
  }

  // Creamos el evento
  dbError = Db.Conn.Insert(evento)
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
      mensajeError = Mensaje.NO_HAY_DATOS_EVENTO + strconv.FormatInt(evento.Id, 10)
    }else{
      mensajeError = dbError.Error()
    }
  }

  return mensajeError
}


//-----------------------------------------------------------------------------
// Búsqueda de los eventos a los que está asociado un usuario
//
//-----------------------------------------------------------------------------
func (evento *Evento) ByUsuario(idUsuario int64) ([]Evento, string){

  var eventos []Evento

  Db.Open()
  defer Db.Close()

  query := "SELECT eventos.* FROM eventos, usuario_eventos WHERE id = id_evento and id_usuario = " + strconv.Itoa(int(idUsuario))
  _, dbError := Db.Conn.Query(&eventos, query)

  if dbError != nil{
    return nil, dbError.Error()
  }

/*  if eventos == nil{
    return nil, NO_HAY_EVENTOS_USUARIO + strconv.Itoa(int(idUsuario))
  }*/
  return eventos, ""
}
