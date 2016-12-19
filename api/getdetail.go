package api

import (
  "os"
  "io"
  "strings"
  "fmt"
  "net/http"
  "encoding/json"
  "bufio"
  "io/ioutil"
  "time"
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
    f,_:=os.OpenFile("./data/"+serialString+"/info.txt",os.O_RDONLY,0777)
    defer f.Close()
    fbuf:=bufio.NewReader(f)
    for i:=0 ; i<5 ; i++ {
      artInfo,_:=fbuf.ReadString('\n')
      switch i {
        case 0:
          out.Data.Title=strings.TrimRight(artInfo,"\n")
        case 1:
          out.Data.Begin=strings.TrimRight(artInfo,"\n")
          out.Data.Year=strings.Split(strings.TrimRight(artInfo,"\n"),"-")[0]
        case 2:
          out.Data.End=strings.TrimRight(artInfo,"\n")
        case 3:
          out.Data.People=strings.TrimRight(artInfo,"\n")
        case 4:
          out.Data.Place=strings.TrimRight(artInfo,"\n")
      }
    }
    desInfo,_:=ioutil.ReadFile("./data/"+serialString+"/description.txt")
    out.Data.Description=string(desInfo)

    //get Item data
    f,_=os.OpenFile("./data/"+serialString+"/item.csv",os.O_RDONLY,0777)
    defer f.Close()
    fbuf=bufio.NewReader(f)
    for true {
      itemInfo,err:=fbuf.ReadString('\n')
      //check eof
      if err==io.EOF {
        if len(out.Data.Items)==0 {
          out.Status=404
        }
        break
      }

      //split data
      itemInfoS:=strings.Split(itemInfo,",")
      for i:=0 ; i<len(itemInfoS) ; i++ {
        itemInfoS[i]=strings.TrimLeft(itemInfoS[i]," ")
      }

      //assign data
      item:=Item{itemInfoS[0],itemInfoS[1]}
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
