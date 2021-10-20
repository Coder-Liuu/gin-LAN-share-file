package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type File struct {
	Name  string
	Size  string
	IsDir bool
}

type Headline_Path struct {
    Path string 
    Url string
}

func getFiles(search string) ([]File, string) {
	fmt.Println("search", search)
	str, _ := os.Getwd()
	if search != "" {
		str = search
	}

	result := make([]File, 0)
	files, err := ioutil.ReadDir(str)
	fmt.Println("当前文件下的目录")
	if err == nil {
		for _, file := range files {
			name := file.Name()
			size := fmt.Sprintf("%d B", file.Size())
			isdir := file.IsDir()
			result = append(result, File{name, size, isdir})
		}
	}
	return result, str
}

func download(c *gin.Context) {
    /* 文件下载入口 */
	filePath := c.Query("file")
	fmt.Println(filePath)
	paths := strings.Split(filePath, "/")
	name := paths[len(paths)-1]

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filePath)
}

func delete_str(s string) string{
    var pos = 0
    for index, value := range s {
        if(value == '/') {
            pos = index
        }
    }
    return s[:pos]
}

func deal_headlines(path_list []string) []string{
    res := []string{}
    path := ""
    for _, value := range path_list {
        path += "/" + value 

        if value == "" {
            path = ""
        }

        fmt.Println("value = ", path)
        res = append(res,path)
    }
    return res
}

func get_headlines(path string) []Headline_Path{
    headline_path := strings.Split("root"+path, "/")
    headline_url := deal_headlines(strings.Split(path,"/"))
    headlines := make([]Headline_Path, 0)
    for index, value := range headline_path {
        headlines = append(headlines,Headline_Path{value,headline_url[index]})
    }
    return headlines
}


func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		search := c.Query("path")
		fmt.Println(search)
		files, path := getFiles(search)
        fmt.Println(get_headlines(path))

         
		c.HTML(200, "home.html", gin.H{
			"files": files,
			"current_path": path,
            "pre_path" : delete_str(path),
            "headlines" : get_headlines(path),
		})

	})

	r.GET("/download", download)
	r.Run(":8888")
}

// /home/liuyang/code/go//src
