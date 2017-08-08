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


// Respresenta un error:
//   - ErrorCode => código específico dentro de la aplicación
//   - StatusCode => http status code
//   - Msg => texto describiendo el error
type ErrorMsg struct{
    ErrorCode int64
    StatusCode int
    Msg string
}

// Datos que devolveremos al cliente
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
//------------------------------------------------------------------------------
func dbErrorParse(err string) (string, int64){
  Parts := strings.Split(err, ":")
  errorMessage := Parts[1]
  Code := strings.Split(Parts[0], " ")
  errorCode, _ := strconv.ParseInt(Code[1], 10, 32)

  return errorMessage, errorCode
}

//------------------------------------------------------------------------------
// Devuelve la estructura de error en formato json
//------------------------------------------------------------------------------
func (error *ErrorMsg) Dispatch() (int, string){

  // Si es un error de bd extraemos el código de mysql para decidir el status code
  if error.ErrorCode == DB_ERROR{
    error.Msg, error.ErrorCode = dbErrorParse(error.Msg)
    log.Println(error.Msg)
    log.Println(error.ErrorCode)
  }

  var errorMessage string
  var statusCode int
  var errorCode int64

  switch(error.ErrorCode){
    case GENERIC_ERROR:
      errorMessage = error.Msg
      errorCode = error.ErrorCode
      statusCode = 409
    // Puerto en uso => salta cuando el mysql no está disponible
    case 3306:
      errorMessage = "Database not available"
      errorCode = error.ErrorCode
      statusCode = 404
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
