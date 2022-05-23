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

   contentType := r.Header.Get("Content-Type")

   if(contentType=="application/x-www-form-urlencoded") {
	   r.ParseForm()
   	if len(r.Form) > 0 {
         for k,v := range r.Form {
            params[k] = v[0]
         }
      }
   }

   if(strings.Contains(contentType, "multipart/form-data")) {
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
               v      := "/tmp/"+strconv.FormatInt(time.Now().UnixNano(), 10)+"_"+part.FileName()
               dst, _ := os.Create(v)
               defer dst.Close()
               io.Copy(dst, part)

               if strings.Contains(k, "[") {
                  rk := k[:len(k)-2]
                  if _, ok := params[rk];ok {
                     params[rk] = params[rk]+":"+v
                  } else {
                     params[rk] = v
                  }
               } else {
                  params[k] = v
               }
            }
         }
      }
   }

	return params
}
