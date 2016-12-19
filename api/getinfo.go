package api

import (
  "os"
  "io"
  "strings"
  "fmt"
  "net/http"
  "encoding/json"
  "bufio"
  "strconv"
  "time"
)

type Getinfo struct {
}

type IndexInfo struct{
  Title   string
  Serial  string
  Cover   string
  Begin   string
  End     string
}

type IndexInfoPkg struct{
  Status  int
  Data    []IndexInfo
}

//index info
func (gihandle Getinfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  //GET
  if r.Method != "GET" {
    fmt.Println(w,"404 page not found")
    return
  }

  //pre process
  r.ParseForm()
  var out IndexInfoPkg
  out.Status=200
  pathS:=strings.Split(r.URL.Path[1:],"/")
  //assign page,year
  validPage,validYear:=true,true
  var pageString,yearString string
  if len(pathS)==2{
    validYear=false
    pageString=pathS[1]
  }else if len(pathS)==3{
    pageString=pathS[1]
    yearString=pathS[2]
  }else{
    validPage,validYear=false,false
    out.Status=404
  }

  //show info
  fmt.Println(time.Now())
  fmt.Println("/getinfo from",r.RemoteAddr)
  fmt.Println("Page:",pageString,"Year:",yearString)
  fmt.Println("====================================")

  //process
  if validPage {
    page,_:=strconv.Atoi(pageString)
    //get data
    f,_:=os.OpenFile("./data/info.csv",os.O_RDONLY,0777)
    defer f.Close()
    fbuf:=bufio.NewReader(f)
    for i,times:=(page-1)*13,0 ; i<page*13 ; times++ {
      info,err:=fbuf.ReadString('\n')
      //check eof
      if err==io.EOF {
        if i==(page-1)*13{
          out.Status=404
          validPage=false
        }
        break
      }

      //check year
      if validYear {
        year:=strings.TrimLeft((strings.Split((strings.Split(info,","))[3],"-")[0])," ")
        if year != yearString {
          times--
          continue
        }
      }

      //valid or not
      if i!=times {
        continue;
      }else{
        //valid -> add to out
        i++;
        infoS:=strings.Split(info,",")
        for a:=0 ; a<len(infoS) ; a++ {
          infoS[a]=strings.TrimLeft(infoS[a]," ")
        }
        data:=IndexInfo{infoS[0],infoS[1],infoS[2],infoS[3],infoS[4]}
        out.Data=append(out.Data,data)
      }
    }
  }

  //make json and send
  j,_:=json.Marshal(out)
  w.Header().Set("Content-Type", "application/json")
  w.Write(j)
}
