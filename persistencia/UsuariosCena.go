package persistencia

import(
  "strconv"
)

type UsuarioCena struct{
  IdUsuario int64
  IdCena int64
  Confirmado string
}


func (usuario *UsuarioCena)ByCena(idCena int64) ([]UsuarioCena, string){

  var usuarios []UsuarioCena

  Db.Open()

  query := "SELECT * FROM usuarios_cena WHERE idcena=" + strconv.Itoa(int(idCena))
  rows, dbError := Db.Conn.Query(query)

  if dbError != nil{
    return nil, dbError.Error()
  }

  for rows.Next(){
    usuario := UsuarioCena{}
    rows.Scan(&usuario.IdUsuario, &usuario.IdCena, &usuario.Confirmado)

    usuarios = append(usuarios, usuario)
  }

  Db.Close()

  return usuarios, ""

}
