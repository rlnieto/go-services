package main

import(
  _ "github.com/go-sql-driver/mysql"
  Api "rlnieto.org/go-services/api"
)

/*------------------------------------------------------------------------------
 main


------------------------------------------------------------------------------*/
func main(){

  api := Api.ApiServer{}
  api.Start()

}
