package persistencia

import(
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

const DB_SCHEMA = "cenas"

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
  db, err := sql.Open("mysql", "root@/" + DB_SCHEMA)
  if err != nil {
    //fmt.Println("Error conectando con la bd")
    panic("Imposible conectar con la BD")
  }

  Db.Schema = DB_SCHEMA
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
