package persistencia

import(
  "strconv"
  "strings"
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
func (usuario *UsuarioEvento)AltaEnEvento(idEvento int64, idUsuarios []string)(string){

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
  sql := "DELETE FROM usuarios_evento WHERE idevento=? AND idusuario IN (" + strings.Join(idUsuarios, ",") + ")"
  //sql := "DELETE FROM usuarios_evento WHERE idevento=? AND idusuario IN (" + "7" + ")"

  query, dbError := conexion.Prepare(sql)
  if dbError != nil{
    return dbError.Error()
  }

  _, dbError = tx.Stmt(query).Exec(idEvento)
  if dbError != nil{
    tx.Rollback()
    return dbError.Error()
  }

  // Damos de alta los usuarios enviados
  sql = "INSERT INTO usuarios_evento set idusuario=?" +
  ", idevento=?, confirmado='N'"

  queryAlta, dbError := conexion.Prepare(sql)
  if dbError != nil{
    tx.Rollback()

    return dbError.Error()
  }

  // Insertamoss las filas en la tabla
  for _, idUsuario := range(idUsuarios){
    _, dbError := tx.Stmt(queryAlta).Exec(idUsuario, idEvento)

    // Si hubo error hacemos rollback y salimos
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
  sql := "DELETE FROM usuarios_evento WHERE idevento = ? and idusuario in(" + idUsuarios + ")"

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

  query := "SELECT * FROM usuarios_evento WHERE idevento=" + strconv.Itoa(int(idEvento))
  rows, dbError := Db.Conn.Query(query)

  if dbError != nil{
    return nil, dbError.Error()
  }

  for rows.Next(){
    usuario := UsuarioEvento{}
    rows.Scan(&usuario.IdUsuario, &usuario.IdEvento, &usuario.Confirmado)

    usuarios = append(usuarios, usuario)
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

  query := "SELECT count(*) FROM usuarios_evento WHERE idevento= " + strconv.Itoa(int(idEvento))
  dbError := Db.Conn.QueryRow(query).Scan(&usuarios)

  if dbError != nil{
    return 0, dbError.Error()
  }

  return usuarios, ""
}
