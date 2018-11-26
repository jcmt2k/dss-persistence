package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"encoding/hex"
	"encoding/json"

	"github.com/gorilla/mux"
)

type Document struct {
	Id   string
	Name string
	Size int
}

func hashFileMd5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	var docs []Document
	path := "./repositorio/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fPath := fmt.Sprintf(path, f.Name)
		println(fPath)
		md5sum := hashFileMd5(fPath)
		docs = append(docs,
			Document{Id: md5sum, Name: f.Name(), Size: int(f.Size())})
	}

	//	docs = append(docs,
	//		Document{Id: "doc-1", Name: "Report.docx", Size: 1500})
	//	docs = append(docs,
	//		Document{Id: "doc-2", Name: "Sheet.xlsx", Size: 5000})
	//	docs = append(docs,
	//		Document{Id: "doc-3", Name: "Container.tar", Size: 50000})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}
