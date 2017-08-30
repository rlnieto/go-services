package main

import(
  "fmt"
  "encoding/json"
  "github.com/go-pg/pg"
  Persistencia "rlnieto.org/go-services/persistencia"
)

//-----------------------------------------------------------------------------
// Alta de las tablas
//-----------------------------------------------------------------------------
func altaTablas(db *pg.DB) error {
  tablas := []interface{}{&Persistencia.Evento{}, &Persistencia.UsuarioEvento{}, &Persistencia.Usuario{}, &Persistencia.Restaurante{}}

  for _, model := range tablas {
    err := db.DropTable(model, nil)
    if err != nil {
      fmt.Println(err)
       return err
    }

    err = db.CreateTable(model, nil)
    if err != nil {
      fmt.Println(err)
       return err
    }
    fmt.Println("Creadas las tablas")
  }
  return nil
}


//-----------------------------------------------------------------------------
// Carga inicial de datos
//-----------------------------------------------------------------------------
func altaDatos(db *pg.DB){

  var datosLocales = []byte(`[
    {"nombre":"Cúrcuma", "tipoCocina":"moderna", "direccion":"Rúa Marconi", "numero":4, "poblacion":"A Coruña", "provincia":"A Coruña", "telefono":646146642, "precioMedio":15, "valoracion":8, "latitud":0, "longitud":0},
    {"nombre":"Art&sushi", "tipoCocina":"oriental", "direccion":"Juan Flórez", "numero":52, "poblacion":"A Coruña", "provincia":"A Coruña", "telefono":881948605, "precioMedio":22, "valoracion":8, "latitud":0, "longitud":0}
  ]`)

  var datosUsuarios = []byte(`[
    {"nick":"pepe", "nombre":"pepe", "apellido1":"perez", "apellido2":"perez", "email":"pepe@gmail.com", "telefono":555102030, "fechaalta":"2017-05-30"},
    {"nick":"perico", "nombre":"perico", "apellido1":"perez", "apellido2":"lopez", "email":"ppalotes@gmail.com", "telefono":555203040, "fechaalta":"2017-08-12"},
    {"nick":"pepa", "nombre":"pepa", "apellido1":"ruibarbo", "apellido2":"cadillo", "email":"pepa@gmail.com", "telefono":555123456, "fechaalta":"2017-08-23"},
    {"nick":"perica", "nombre":"perica", "apellido1":"coton", "apellido2":"taina", "email":"perika@gmail.com", "telefono":555789012, "fechaalta":"2017-03-10"},
    {"nick":"juan", "nombre":"juan", "apellido1":"DelCorral", "apellido2":"Quilado", "email":"juan@gmail.com", "telefono":555666777, "fechaalta":"2017-08-12"},
    {"nick":"josefa", "nombre":"josefa", "apellido1":"rodriguez", "apellido2":"tirado", "email":"jozefa@gmail.com", "telefono":555000000, "fechaalta":"2017-08-15"}
  ]`)


  // Locales
  var locales []Persistencia.Restaurante
  error := json.Unmarshal(datosLocales, &locales)
  if error != nil{
    fmt.Println("Error parseando json locales: " + error.Error())
  }

  _, dbError := db.Exec("DELETE FROM restaurantes")

  for _, local := range(locales){
    if db.Insert(&local) != nil {
      fmt.Println(dbError.Error())
    }
  }

  fmt.Println("Cargados los datos de locales...")


  // Usuarios
  var usuarios []Persistencia.Usuario
  error = json.Unmarshal(datosUsuarios, &usuarios)
  if error != nil{
    fmt.Println("Error parseando json usuarios: " + error.Error())
  }

  _, dbError = db.Exec("DELETE FROM usuarios")

  for _, usuario := range(usuarios){
    if db.Insert(&usuario) != nil {
      fmt.Println(dbError.Error())
    }
  }

  fmt.Println("Cargados los datos de usuarios...")

}


//-----------------------------------------------------------------------------
// Punto de entrada
//
//-----------------------------------------------------------------------------
func CrearBd(){

  var Db Persistencia.Database

  Db.Open()
  defer Db.Close()

  error := altaTablas(Db.Conn)
  if error != nil{
    fmt.Println(error)
  }

  altaDatos(Db.Conn)


  fmt.Println("Bd creada correctamente!")
}
