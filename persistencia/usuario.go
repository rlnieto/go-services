package persistencia

import(
  "strconv"
  "github.com/go-pg/pg"
  Mensaje "rlnieto.org/go-services/mensajes"
)

type Usuario struct{
  Id int64
  Nick string
  Nombre string
  Apellido1 string
  Apellido2 string
  Email string
  Telefono int64  `sql:",notnull"`
  FechaAlta string
}



//-----------------------------------------------------------------------------
// Búsqueda por clave
//
//-----------------------------------------------------------------------------
func (usuario *Usuario)ById() (string){

  Db.Open()
  defer Db.Close()

  dbError := Db.Conn.Select(usuario)
  if dbError != nil{
    if dbError == pg.ErrNoRows{
      return Mensaje.NO_HAY_USUARIO + strconv.FormatInt(usuario.Id, 10)
    }else{
      return dbError.Error()
    }
  }

  return ""
}
