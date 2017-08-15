package persistencia

import(
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  Config "rlnieto.org/eventos/go-services/config"
)

type Database struct{
  Schema string
  Conn *sql.DB
}

var Db Database

//-----------------------------------------------------------------------------
// Apertura de la conexión con la BD
//
//-----------------------------------------------------------------------------
func (Db *Database) Open(){
  db, err := sql.Open("mysql", "root@/" + Config.DB_SCHEMA)
  if err != nil {
    //fmt.Println("Error conectando con la bd")
    panic("Imposible conectar con la BD")
  }

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
