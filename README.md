# Naive-Bayes

This application was created to use the Naive-Bayes algorithm for analyzing text.  First send sanitized text into the build api call, then the predict and the entropy methods can be used to perform processing on the text.  Additional text can be sent into the build function later to also be placed into the model.  This model is ephemeral and will not retain after the application closes or crashes.

## Run

To run the container issue:

```
docker run -d -p 8080:8080 randysimpson/naive-bayes:latest
```

## Build Model

The `http://localhost:8080/api/v1/build` api method will build the model to be used in language processing.  The payload must be type `application/json` and in the form of:

```json
{
  "data": "This is a test",
}
```

### From data input

The example curl command is:

```
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -d '{"data":"This is a test"}' -X POST http://localhost:8080/api/v1/build
```

example output:

```
HTTP/1.1 201 Created
Content-Type: application/json; charset=UTF-8
Date: Sat, 20 Feb 2021 01:47:22 GMT
Content-Length: 32

{"size":"4","status":"Success"}
```

### From File

This example is using an output file from the sanitize-text app.

```sh
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -d @output.txt -X POST http://localhost:8080/api/v1/build
```

## Predict

To get the probabilities of the next word after `a` we use:

```sh
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8080/api/v1/predict/a
```

example output:

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Sat, 20 Feb 2021 01:47:58 GMT
Content-Length: 11

{"test":1}
```

This means that 100% chance of test after a. 

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

### Docker Build

To build this container using docker issue:

```
docker build -t randysimpson/naive-bayes:v1.1 .
```

### Manual Build

To manually build the source files you will need to get the external dependencies and then build the binary executable file.

```
go get k8s.io/klog
go get github.com/gorilla/mux
go build
```

# Examples

A little more exciting of a prediction:

```
rsimpson@k8server1:~/code/naive-bayes$ curl -i -H "Accept: application/json" -H "Content-Type: application/json" -d '{"data":"1 2 3 4 5 2 1 3 5 1 2"}' -X POST http://localhost:8080/api/v1/build
HTTP/1.1 201 Created
Content-Type: application/json; charset=UTF-8
Date: Sat, 20 Feb 2021 02:51:14 GMT
Content-Length: 33

{"size":"11","status":"Success"}
rsimpson@k8server1:~/code/naive-bayes$ curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8080/api/v1/predict/1
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Sat, 20 Feb 2021 02:51:23 GMT
Content-Length: 31

{"2":0.6666667,"3":0.33333334}
```

With the data we used to build the model we can see that with a 1, there is a 66% chance of a 2 next and a 33% chance of a 3.

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
