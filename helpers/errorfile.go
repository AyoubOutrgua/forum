package helpers

import "net/http"
  



func Errorhandler(w http.ResponseWriter , errors string , er int  ){


http.Error(w,errors,er)

}