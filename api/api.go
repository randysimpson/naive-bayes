/*
MIT License

# Copyright (©) 2024 - Randall Simpson

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
*/
package api

import "github.com/gin-gonic/gin"

func HandleRequests() {
	router := gin.Default()

	router.GET("/", Index)
	router.POST("/api/v1/build", BuildModel)
	router.POST("/api/v1/entropy", TestEntropy)
	router.GET("/api/v1/model", GetModel)
	router.GET("/api/v1/predict/:first/:second/:third", QuadPredict)
	router.GET("/api/v1/predict/:first/:second", TriPredict)
	router.GET("/api/v1/predict/:first", Predict)
	router.POST("/api/v1/clear", Clear)

	router.Run("0.0.0.0:8080")
}
