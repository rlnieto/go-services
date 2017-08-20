package persistencia

import(
  "fmt"
  "strings"
  "strconv"
)

type UsuarioEvento struct{
  IdUsuario int64
  IdEvento int64
  Confirmado string
}


//-----------------------------------------------------------------------------
// Alta de los usuarios de un evento
// Borra los usuarios enviados para evitar duplicados
//
//-----------------------------------------------------------------------------
func (dummy *UsuarioEvento)AltaEnEvento(idEvento int64, idUsuarios []string)(string){

  Db.Open()
  defer Db.Close()

  conexion := Db.Conn

  // Iniciamos la transacción
  tx, dbError := conexion.Begin()
  if dbError != nil{
    return dbError.Error()
  }

  // Borramos los usuarios que ya hay asociados
  // Le concatenamos los usuarios porque al pasarlos compo parámetro no lo hacía bien
  // (parece que borraba sólo el primero...)
  sql := "DELETE FROM usuario_eventos WHERE id_evento=$1::integer AND id_usuario IN (" + strings.Join(idUsuarios, ",") + ")"
  //sql := "DELETE FROM usuario_eventos WHERE id_evento=$1::integer AND id_usuario IN (7)"

  query, dbError := conexion.Prepare(sql)
  if dbError != nil{
    return dbError.Error()
  }

fmt.Println(sql)
fmt.Println(idEvento)

  _, dbError = tx.Stmt(query).Exec(idEvento)
  if dbError != nil{
    tx.Rollback()
    return dbError.Error()
  }

  // Damos de alta los usuarios enviados
  var usuario UsuarioEvento

  // Insertamoss las filas en la tabla
  for _, idUsuario := range(idUsuarios){

    usuario.IdEvento = idEvento
    usuario.IdUsuario, _ = strconv.ParseInt(idUsuario, 10, 32)
    usuario.Confirmado = "N"

    dbError = Db.Conn.Insert(&usuario)
    if dbError != nil{
      tx.Rollback()
      return dbError.Error()
    }
  }

  tx.Commit()

  return ""
}

//-----------------------------------------------------------------------------
// Borrado de usuarios de una evento
//
//-----------------------------------------------------------------------------
func (usuario *UsuarioEvento)BajaEnEvento(idEvento int64, idUsuarios string) (string){

  Db.Open()
  defer Db.Close()

  // Borramos los asistentes que nos han enviado
  sql := "DELETE FROM usuario_eventos WHERE id_evento = $1::integer and id_usuario in(" + idUsuarios + ")"

  query, dbError := Db.Conn.Prepare(sql)
  if dbError != nil{
    return dbError.Error()
  }

  _, dbError = query.Exec(idEvento)

  if dbError != nil{
    return dbError.Error()
  }

  return ""
}

//-----------------------------------------------------------------------------
// Consulta de una evento con sus usuarios
//
//-----------------------------------------------------------------------------
func (usuario *UsuarioEvento)ByEvento(idEvento int64) ([]UsuarioEvento, string){

  var usuarios []UsuarioEvento

  Db.Open()
  defer Db.Close()

  query := "SELECT * FROM usuario_eventos WHERE id_evento=" + strconv.Itoa(int(idEvento))
  _, dbError := Db.Conn.Query(&usuarios, query)

  if dbError != nil{
    return nil, dbError.Error()
  }

  return usuarios, ""
}

//-----------------------------------------------------------------------------
// Comprueba si existen usuarios para un evento
//
//-----------------------------------------------------------------------------
func (usuario *UsuarioEvento)NumeroUsuariosEvento(idEvento int64) (int, string){

  var usuarios int

  Db.Open()
  defer Db.Close()

  //query := "SELECT count(*) FROM usuarios_evento WHERE idevento= " + strconv.Itoa(int(idEvento))
  usuarios, dbError := Db.Conn.Model(usuario).Count()
  if dbError != nil{
    return 0, dbError.Error()
  }

  return usuarios, ""
}
