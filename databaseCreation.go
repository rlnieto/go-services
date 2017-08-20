package main

import(
  "fmt"
  "github.com/go-pg/pg"
  Persistencia "rlnieto.org/pruebas/postgres/persistencia"
)


func altaTablas(db *pg.DB) error {
  tablas := []interface{}{&Persistencia.Evento{}, &Persistencia.UsuarioEvento{}}

  for _, model := range tablas {
    err := db.CreateTable(model, nil)
    if err != nil {
       return err
    }
    fmt.Println("Creada tabla ")
  }
  return nil
}


func CrearBd(){

  var Db Persistencia.Database

  Db.Open()
  defer Db.Close()

  error := altaTablas(Db.Conn)
  if error != nil{
    panic(error)
  }

  fmt.Println("Bd creada correctamente!")
}
