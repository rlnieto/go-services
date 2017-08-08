package persistencia

import(
  "strconv"
  "log"
)

type Cena struct{
  Id int64
  Fecha string
  Hora string
  IdLocal int64
  Motivo  string
  Menu  string
  PrecioEstimado int64
  IdOrganizador int64
}


//-----------------------------------------------------------------------------
// Búsqueda por clave
//
//-----------------------------------------------------------------------------
func (cena *Cena) ById() (string){
  mensajeError := ""

  Db.Open()

  log.Println("Abierta la BD")
  query := "SELECT * FROM cena WHERE id=" + strconv.Itoa(int(cena.Id))
  dbError := Db.Conn.QueryRow(query).Scan(&cena.Id,&cena.Fecha,&cena.Hora,&cena.IdLocal,&cena.Motivo,&cena.Menu,&cena.PrecioEstimado,&cena.IdOrganizador)

  if dbError != nil{
    mensajeError = dbError.Error()
  }

  Db.Close()
  return mensajeError
}
