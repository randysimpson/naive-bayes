/*MIT License

Copyright (©) 2019 - Randall Simpson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.*/
package api
 
import (
    "encoding/json"
    "fmt"
    "net/http"
    "naive-bayes/model"
    "io/ioutil"
    "io"
    "k8s.io/klog"
    "github.com/gorilla/mux"
)
 
func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Welcome!")
}
 
func BuildModel(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    klog.Errorf("error: %+v", err)
  }
  
  if err := r.Body.Close(); err != nil {
    klog.Errorln(err)
  }
  
  //convert json into map
  var jsonBody interface{}
  json.Unmarshal(body, &jsonBody)
  
  m := jsonBody.(map[string]interface{})

  s := m["data"].(string)
  count, err := model.AddData(s)
  if err != nil {
    klog.Errorf("error: %+v", err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
  
  t := map[string]interface{}{
    "status": "Success",
    "size":  fmt.Sprintf("%d", count),
  }
  if err := json.NewEncoder(w).Encode(t); err != nil {
    klog.Errorln(err)
  }
}

func TestEntropy(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    klog.Errorf("error: %+v", err)
  }
  
  if err := r.Body.Close(); err != nil {
    klog.Errorln(err)
  }
  
  //convert json into map
  var jsonBody interface{}
  json.Unmarshal(body, &jsonBody)
  
  m := jsonBody.(map[string]interface{})

  s := m["data"].(string)
  
  entropy, err := model.GetEntropy(s)
  if err != nil {
    klog.Errorf("error: %+v", err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  
  t := map[string]interface{}{
    "entropy": fmt.Sprintf("%f", entropy),
  }
  if err := json.NewEncoder(w).Encode(t); err != nil {
    klog.Errorln(err)
  }  
}

func GetModel(w http.ResponseWriter, r *http.Request) {
  count := model.GetCount()

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  
  t := map[string]interface{}{
    "count": fmt.Sprintf("%d", count),
  }
  if err := json.NewEncoder(w).Encode(t); err != nil {
    klog.Errorln(err)
  }
}

func Predict(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  word := vars["word"]
  options, err := model.GetNext(word)
  if err != nil {
    klog.Errorf("error: %+v", err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  
  enc := json.NewEncoder(w)
  enc.SetEscapeHTML(false)
  //if err := json.NewEncoder(w).Encode(options); err != nil {
  if err := enc.Encode(options); err != nil {
    klog.Errorln(err)
  }
}

func TriPredict(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  first := vars["first"]
  second := vars["second"]
  options, err := model.GetTriNext(first, second)
  if err != nil {
    klog.Errorf("error: %+v", err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)

  enc := json.NewEncoder(w)
  enc.SetEscapeHTML(false)
  //if err := json.NewEncoder(w).Encode(options); err != nil {
  if err := enc.Encode(options); err != nil {
    klog.Errorln(err)
  }
}

func QuadPredict(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  first := vars["first"]
  second := vars["second"]
  third := vars["third"]
  options, err := model.GetQuadNext(first, second, third)
  if err != nil {
    klog.Errorf("error: %+v", err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)

  enc := json.NewEncoder(w)
  enc.SetEscapeHTML(false)
  //if err := json.NewEncoder(w).Encode(options); err != nil {
  if err := enc.Encode(options); err != nil {
    klog.Errorln(err)
  }
}
