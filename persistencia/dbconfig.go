package persistencia

import(
  "github.com/go-pg/pg"
  Config "rlnieto.org/go-services/config"
)

type Database struct{
  Schema string
  Conn *pg.DB
}

var Db Database

//-----------------------------------------------------------------------------
// Apertura de la conexión con la BD
//
//-----------------------------------------------------------------------------
func (Db *Database) Open() {

  db := pg.Connect(&pg.Options{
    User: "guest",
    Password: "guest",
    Database: "eventos",
  })

  Db.Schema = Config.DB_SCHEMA
  Db.Conn = db

}


//-----------------------------------------------------------------------------
// Cierre de la conexión con la BD
//
//-----------------------------------------------------------------------------
func (Db *Database)Close(){
    if Db != nil && Db.Conn != nil{
      Db.Conn.Close()
    }
}
