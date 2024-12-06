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

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/randysimpson/naive-bayes/model"

	"k8s.io/klog"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome!"})
}

type buildReq struct {
	Data   string   `json:"data,omitempty"`
	Tokens []string `json:"tokens,omitempty"`
}

func BuildModel(c *gin.Context) {
	var newBuildRequest buildReq
	if err := c.BindJSON(&newBuildRequest); err != nil {
		klog.Errorf("error: %+v", err)
		return
	}

	totalCount := 0
	if len(newBuildRequest.Data) > 0 {
		count, err := model.AddData(strings.Fields(newBuildRequest.Data))
		if err != nil {
			klog.Errorf("error: %+v", err)
		}
		totalCount += count
	}
	if len(newBuildRequest.Tokens) > 0 {
		count, err := model.AddData(newBuildRequest.Tokens)
		if err != nil {
			klog.Errorf("error: %+v", err)
		}
		totalCount += count
	}

	t := map[string]interface{}{
		"status": "Success",
		"size":   totalCount,
	}
	c.JSON(http.StatusCreated, t)
}

func TestEntropy(c *gin.Context) {
	var newBuildRequest buildReq
	if err := c.BindJSON(&newBuildRequest); err != nil {
		klog.Errorf("error: %+v", err)
		return
	}

	entropy := 0.0
	if len(newBuildRequest.Data) > 0 {
		e, err := model.GetEntropy(strings.Fields(newBuildRequest.Data))
		if err != nil {
			klog.Errorf("error: %+v", err)
		}
		entropy = e
	}
	if len(newBuildRequest.Tokens) > 0 {
		e, err := model.GetEntropy(newBuildRequest.Tokens)
		if err != nil {
			klog.Errorf("error: %+v", err)
		}
		entropy = e
	}

	t := map[string]interface{}{
		"entropy": entropy,
	}
	c.JSON(http.StatusOK, t)
}

func GetModel(c *gin.Context) {
	t := model.GetModel()
	c.JSON(http.StatusOK, t)
}

func Predict(c *gin.Context) {
	first := c.Param("first")

	options, err := model.GetNext(first)
	if err != nil {
		klog.Errorf("error: %+v", err)
	}

	c.JSON(http.StatusOK, options)
}

func TriPredict(c *gin.Context) {
	first := c.Param("first")
	second := c.Param("second")
	options, err := model.GetTriNext(first, second)
	if err != nil {
		klog.Errorf("error: %+v", err)
	}

	c.JSON(http.StatusOK, options)
}

func QuadPredict(c *gin.Context) {
	first := c.Param("first")
	second := c.Param("second")
	third := c.Param("third")
	options, err := model.GetQuadNext(first, second, third)
	if err != nil {
		klog.Errorf("error: %+v", err)
	}

	c.JSON(http.StatusOK, options)
}

func Clear(c *gin.Context) {
	model.ClearModel()

	t := map[string]interface{}{
		"status": "clear",
	}
	c.JSON(http.StatusOK, t)
}
