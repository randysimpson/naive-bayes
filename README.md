# Naive-Bayes

This application was created to use the Naive-Bayes algorithm for analyzing text.  First send sanitized text into the build api call, then the predict and the entropy methods can be used to perform processing on the text.  Additional text can be sent into the build function later to also be placed into the model.  This model is ephemeral and will not retain after the application closes or crashes.

## Build Model

The `http://localhost:8080/api/v1/build` api method will build the model to be used in language processing.  The payload must be type `application/json` and in the form of:

```json
{
  "status": "Success",
  "size": 4,
  "data": "This is a test",
}
```

### From File

This example is using an output file from the sanitize-text app.

```sh
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -d @output.txt -X POST http://localhost:8080/api/v1/build
```

## Predict

To get the probabilities of the next word of `this` we use:

```sh
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8080/api/v1/predict/this
```

Use url encoding if we want to find `<error>`:

```sh
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8080/api/v1/predict/%3Cerror%3E
```

## Entropy

check completely different data

```sh
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -d '{"data":"Somthing is very wrong here"}' -X POST http://localhost:8080/api/v1/entropy
```

check from sanitized output:

```sh
cat output2.curl | jq -r '.data' | sed -e 's/<\/line> /<\/line>\n/g' | while read LINE ; do echo '{"data":"'$LINE'"}';  done | while read LINE ; do echo 'curl -i -H "Accept: application/json" -H "Content-Type: application/json" -d '"'"''$LINE''"'"' -X POST http://localhost:8080/api/v1/entropy';  done
```

## Installation

This is a microservice which has been written based on REST-API to allow for deployment from a docker container.

### Manual Build

To manually build the source files you will need to get the external dependencies and then build the binary executable file.

```
go get k8s.io/klog
go get github.com/gorilla/mux
go build
```

# Licence

MIT License

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
SOFTWARE.