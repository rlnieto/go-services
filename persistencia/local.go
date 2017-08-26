package persistencia

import(
  "github.com/go-pg/pg"
)

type Restaurante struct{
  Id int64
  Nombre string
  TipoCocina string
  Direccion string
  Numero int  `sql:",notnull"`
  Poblacion string
  Provincia string
  Telefono int64  `sql:",notnull"`
  PrecioMedio float32
  Valoracion int64  `sql:",notnull"`
  Longitud float32  `sql:",notnull"`
  Latitud float32  `sql:",notnull"`
}

//-----------------------------------------------------------------------------
// Búsqueda por clave
//
//-----------------------------------------------------------------------------
func (local *Restaurante)ById() (string){

  Db.Open()
  defer Db.Close()

  dbError := Db.Conn.Select(local)
  if dbError != nil{
    if dbError == pg.ErrNoRows{
      //return NO_HAY_DATOS_LOCAL + strconv.FormatInt(local.Id, 10)
      return ""
    }else{
      return dbError.Error()
    }
  }

  return ""
}
