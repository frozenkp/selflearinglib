package main

import(
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "os"
  //"os/exec"
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
    for b:=0 ; b<13 ; b++{
      m:=u["data"].([]interface{})[b].(map[string]interface{})
      url2:=fmt.Sprintf("http://hall.lib.nctu.edu.tw/api/artData/%s",m["serial"])
      resp2,_:=http.Get(url2)
      body2,_:=ioutil.ReadAll(resp2.Body)
      defer resp2.Body.Close()
      var u2 map[string]interface{}
      json.Unmarshal(body2,&u2)
      mm2:=u2["data"].(map[string]interface{})
      mm:=u2["data"].(map[string]interface{})["items"].([]interface{})
      //os.Mkdir(m["serial"].(string),1775)
      ff,_:=os.OpenFile("./"+m["serial"].(string)+"/info.txt",os.O_WRONLY | os.O_CREATE, 0777)
      fmt.Fprintf(ff,"%s\n%s\n%s\n%s\n%s\n",mm2["title"],mm2["begin"],mm2["end"],mm2["people"],mm2["place"])
      ff.Close()
      ff,_=os.OpenFile("./"+m["serial"].(string)+"/description.txt",os.O_WRONLY | os.O_CREATE, 0777)
      fmt.Fprintf(ff,"%s\n",mm2["description"])
      ff.Close()
      ff,_=os.OpenFile("./"+m["serial"].(string)+"/item.csv",os.O_WRONLY | os.O_CREATE, 0777)
      var cover string
      for b:=0;b<len(mm);b++{
        mmm:=mm[b].(map[string]interface{})
        /*o:=exec.Command("wget","-P","./"+m["serial"].(string),mmm["imgPath"].(string))
        o.Run()*/
        fmt.Fprintf(ff,"%s, %s, \n",mmm["title"],"http://140.113.66.249:65534/data/"+m["serial"].(string)+"/"+(strings.Split(mmm["imgPath"].(string),"/"))[4])
        if b==0 {
          cover="http://140.113.66.249:65534/data/"+m["serial"].(string)+"/"+(strings.Split(mmm["imgPath"].(string),"/"))[4]
        }
      }
      fmt.Fprintf(f,"%s, %s, %s, %s, %s, \n",m["title"],m["serial"],cover,m["begin"],m["end"])
    }
  }
}
