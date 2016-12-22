package api

import (
  "os"
  "strings"
  "fmt"
  "net/http"
  "encoding/json"
  "time"
  "encoding/csv"
)

type Getdetail struct {
}

type Item struct{
  Title   string
  ImgPath string
}

type ArtInfo struct{
  Title       string
  Begin       string
  End         string
  People      string
  Place       string
  Description string
  Year        string
  Items       []Item
}

type ArtInfoPkg struct{
  Status  int
  Data    ArtInfo
}

//getdetail
func (gdhandle Getdetail) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  //GET
  if r.Method!="GET" {
    fmt.Println("404 page not found")
    return
  }

  //pre process
  r.ParseForm()
  var out ArtInfoPkg
  out.Status=200
  pathS:=strings.Split(r.URL.Path[1:],"/")
  serialString,serialValid:="",true
  if len(pathS)==2 {
    serialString=pathS[1]
  }else{
    out.Status=404
    serialValid=false
  }
  //check Serial content valid
  if serialString=="" {
    out.Status=404
    serialValid=false
  }else if !exist("./data/"+serialString) {
    out.Status=404
    serialValid=false
  }

  //show info
  fmt.Println(time.Now())
  fmt.Println("/getdetail from",r.RemoteAddr)
  fmt.Println("Serial:",serialString)
  fmt.Println("====================================")

  //process
  if serialValid {
    //get ArtInfo data
    f,_:=os.OpenFile("./data/"+serialString+"/info.csv",os.O_RDONLY,0777)
    defer f.Close()
    r:=csv.NewReader(f)
    result,_:=r.Read()
    out.Data.Title=result[0]
    out.Data.Begin=result[1]
    out.Data.Year=strings.Split(result[1],"-")[0]
    out.Data.End=result[2]
    out.Data.People=result[3]
    out.Data.Place=result[4]
    out.Data.Description=result[5]

    //get Item data
    f,_=os.OpenFile("./data/"+serialString+"/item.csv",os.O_RDONLY,0777)
    defer f.Close()
    r=csv.NewReader(f)
    results,_:=r.ReadAll()
    for row:=0;row<len(results);row++{
      item:=Item{results[row][0],results[row][1]}
      out.Data.Items=append(out.Data.Items,item)
    }
  }

  //make json and send
  j,_:=json.Marshal(out)
  w.Header().Set("Content-Type", "application/json")
  w.Write(j)
}

func exist(path string)bool{
  _,err:=os.Stat(path)
  if err == nil {
    return true
  }else{
    return false
  }
}
