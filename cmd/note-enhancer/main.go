package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)
const PDF_READER_LOCATION = "C:/Program Files (x86)/jisupdf/JisuPdf.exe"

func OpenPDFHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var values = request.URL.Query()
	var bookName = values.Get("name")
	responseWriter.Header().Add("Content-Type", "application/json")
	responseWriter.Header().Add("Access-Control-Allow-Origin",  "*")
	responseWriter.Header().Add("Access-Control-Allow-Methods", "*")
	responseWriter.Header().Add("Access-Control-Allow-Headers", "*")
	fmt.Fprint(responseWriter, "{\"success\": " + strconv.FormatBool(searchBookAndOpen(bookName)) +"}")
	
}

func getFileList(path string) *list.List {
	fs,_:= ioutil.ReadDir(path)
	var list = list.New()
	for _,file:=range fs{
			if file.IsDir(){
					list.PushBackList(getFileList(path+file.Name()+"/"))
			}else{
					list.PushBack(path+file.Name())
			}
	}
	return list
}

func searchBookAndOpen(bookName string) bool {

	if len(bookName) == 0 {
		return false
	}

	var fileList = getFileList("d:/learning/ebook/")
	for i:= fileList.Front(); i != nil; i = i.Next() {
		var file string = i.Value.(string)
		if strings.HasSuffix(file, ".pdf") && strings.Contains(strings.ToUpper(file), strings.ToUpper(bookName)) {
			cmd := exec.Command(PDF_READER_LOCATION, file)
			err := cmd.Start()
			if err != nil {
				log.Fatal(err)
				return false
			}
		
			return true
		}
	}
	return false
}

func main() {

	// pdf打开指令
	http.HandleFunc("/openPDF", OpenPDFHandler)
	http.ListenAndServe(":12945", nil)

}