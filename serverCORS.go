package main

import (
  "log"
  "net/http"
  "github.com/semicircle/gocors"
  "./api"
)

func main() {
  //CORS
  cors := gocors.New()

  //photo
  http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("./data"))))

  //getinfo
  var gihandle api.Getinfo
  http.Handle("/getinfo/", cors.Handler(gihandle))

  //getdetail
  var gdhandle api.Getdetail
  http.Handle("/getdetail/", cors.Handler(gdhandle))

  log.Fatal(http.ListenAndServe(":65534", nil))
}
