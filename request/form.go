package request

import "io"
import "os"
import "time"
import "strings"
import "strconv"
import "net/http"
import "io/ioutil"

type Form struct {}

func (F *Form) Parse(r *http.Request)map[string]string {
	params := map[string]string{}

   values := r.URL.Query()
   for k, v := range values {
      if strings.Contains(k, "[") {
         rk := k[:len(k)-2]
         params[rk] = strings.Join(v, ":")
      } else {
         params[k] = v[0]
      }
   }


   contentType := r.Header.Get("Content-Type")

   if(contentType=="application/x-www-form-urlencoded") {
	   r.ParseForm()
   	if len(r.Form) > 0 {
         for k,v := range r.Form {
            if strings.Contains(k, "[") {
               rk := k[:len(k)-2]
               params[rk] = strings.Join(v, ":")
            } else {
               params[k] = v[0]
            }
         }
      }
   }

   if len(contentType)>=19 && contentType[0:19]=="multipart/form-data" {
      reader, err := r.MultipartReader()
      if err == nil {
         for {
            part, err := reader.NextPart()
            if err == io.EOF {
               break
            }

            if part.FileName() == "" {
                  data, _  := ioutil.ReadAll(part)
                  k        := part.FormName()
                  v        := string(data)
                  params[k] = v

            } else {
               k      := part.FormName()
               v      := "/tmp/"+strconv.FormatInt(time.Now().UnixNano(), 10)+".nova.upload"
               dst, _ := os.Create(v)
               defer dst.Close()
               io.Copy(dst, part)

               if strings.Contains(k, "[") {
                  rk := k[:len(k)-2]
                  if _, ok := params[rk];ok {
                     params[rk] = params[rk]+","+v
                     params[rk+"_name"] = params[rk+"_name"]+","+part.FileName()
                  } else {
                     params[rk] = v
                     params[rk+"_name"] = part.FileName()
                  }
               } else {
                  params[k] = v
                  params[k+"_name"] = part.FileName()
               }
            }
         }
      }
   }

	return params
}
