package error

import(
  "encoding/json"
  "strings"
  "strconv"
  "log"
)

// Constantes
const GENERIC_ERROR = 1      // ErrorCode => error genérico
const DB_ERROR = 2           // ErrorCode => error de BD
const NOT_FOUND = 3
const PARAMETRO_INCORRECTO = 4

// Respresenta un error:
//   - ErrorCode => código específico dentro de la aplicación
//   - StatusCode => http status code
//   - Msg => texto describiendo el error
type ErrorMsg struct{
    ErrorCode int64
    StatusCode int
    Msg string
}

// Datos que devolveremos al cliente en caso de error
type Response struct {
  Error string "json:error"
  ErrorCode int64 "json:code"
}


// Response para cuando todo fue bien
func (error *ErrorMsg) OkResponse() (string){
  response := Response{}
  response.Error = ""
  response.ErrorCode = 200;

  output, _ := json.Marshal(response)
  return string(output)
}

/*------------------------------------------------------------------------------
 Gestión de errores

------------------------------------------------------------------------------*/
//------------------------------------------------------------------------------
// Separa en código y mensaje los errores de BD
// TODO: hay que darle una vuelta porque es un lío... Hay que ver también
// cómo montamos la gestión de errores para postgres
// Para MySql buscamos errores con el siguiente formato:
//   "Error 1062: Duplicate entry '1-15' for key 1"
//------------------------------------------------------------------------------
func dbErrorParse(err string) (string, int64){

  log.Println(err)

  var errorMessage string
  var errorCode int64

  // Si no aparece la palabra "Error" no buscamos el código numérico
  if strings.Contains(err, "Error"){
    Parts := strings.Split(err, ":")
    Code := strings.Split(Parts[0], " ")

    errorMessage = Parts[1]
    errorCode, _ = strconv.ParseInt(Code[1], 10, 32)
  }else{
    errorMessage = err
    errorCode = -1
  }

  return errorMessage, errorCode
}


//------------------------------------------------------------------------------
// Devuelve la estructura de error en formato json
//------------------------------------------------------------------------------
func (error *ErrorMsg) Dispatch() (int, string){

  // Si es un error de bd extraemos el código de mysql para decidir el status code
  if error.ErrorCode == DB_ERROR{
    error.Msg, error.ErrorCode = dbErrorParse(error.Msg)
  }

  var errorMessage string
  var statusCode int
  var errorCode int64

  switch(error.ErrorCode){
    case GENERIC_ERROR:
      errorMessage = error.Msg
      errorCode = error.ErrorCode
      statusCode = 409
    case NOT_FOUND:
      errorMessage = error.Msg
      errorCode = error.ErrorCode
      statusCode = 404
    case PARAMETRO_INCORRECTO:
      errorMessage = error.Msg
      errorCode = error.ErrorCode
      statusCode = 412
    case 1062:
      errorMessage = "Duplicate entry"
      errorCode = error.ErrorCode
      statusCode = 409
    default:
      errorMessage = error.Msg
      errorCode = 1
      statusCode = 409
  }

  // Creamos la respuesta a enviar al cliente y la devolvemos en formato json
  response := Response{}
  response.Error = errorMessage
  response.ErrorCode = errorCode

  output, _ := json.Marshal(response)
  return statusCode, string(output)

}
