package main

import(
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "os"
  "os/exec"
  "strings"
)

func main(){
  f,_:=os.OpenFile("info.csv",os.O_WRONLY | os.O_CREATE, 0777)
  defer f.Close()
  for a:=1;a<=8;a++{
    url:=fmt.Sprintf("http://hall.lib.nctu.edu.tw/api/galleryList/%d",a)
    resp,_:=http.Get(url)
    body,_:=ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    var u map[string]interface{}
    json.Unmarshal(body,&u)
    if(u["status"].(float64)!=200){
      continue
    }
    mo:=u["data"].([]interface{})
    for b:=0 ; b<len(mo) ; b++{
      m:=mo[b].(map[string]interface{})
      url2:=fmt.Sprintf("http://hall.lib.nctu.edu.tw/api/artData/%s",m["serial"])
      resp2,_:=http.Get(url2)
      body2,_:=ioutil.ReadAll(resp2.Body)
      defer resp2.Body.Close()
      var u2 map[string]interface{}
      json.Unmarshal(body2,&u2)
      mm2:=u2["data"].(map[string]interface{})
      mm:=u2["data"].(map[string]interface{})["items"].([]interface{})
      os.Mkdir(m["serial"].(string),1775)
      ff,_:=os.OpenFile("./"+m["serial"].(string)+"/info.csv",os.O_WRONLY | os.O_CREATE, 0777)
      fmt.Fprintf(ff,"\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",cd(mm2["title"]),cd(mm2["begin"]),cd(mm2["end"]),cd(mm2["people"]),cd(mm2["place"]),cd(mm2["description"]))
      ff.Close()
      ff,_=os.OpenFile("./"+m["serial"].(string)+"/item.csv",os.O_WRONLY | os.O_CREATE, 0777)
      defer ff.Close()
      var cover string
      for c:=0;c<len(mm);c++{
        mmm:=mm[c].(map[string]interface{})
        o:=exec.Command("wget","-P","./"+m["serial"].(string),mmm["imgPath"].(string))
        o.Run()
        fmt.Fprintf(ff,"\"%s\",\"%s\"\n",cd(mmm["title"]),"http://140.113.66.249:65534/data/"+m["serial"].(string)+"/"+(strings.Split(mmm["imgPath"].(string),"/"))[4])
        if c==0 {
          cover="http://140.113.66.249:65534/data/"+m["serial"].(string)+"/"+(strings.Split(mmm["imgPath"].(string),"/"))[4]
        }
      }
      fmt.Fprintf(f,"\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",cd(m["title"]),m["serial"],cover,m["begin"],m["end"])
    }
  }
}

func cd(s interface{})string{
  sn:=s.(string)
  if strings.Contains(sn,"\""){
    sS:=strings.Split(sn,"\"")
    sn=sS[0]
    for b:=1;b<len(sS);b++{
      sn+=("\"\""+sS[b])
    }
  }
  return sn
}
