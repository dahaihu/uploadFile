package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Handle("GET", "/", func(c *gin.Context) {
		pr, pw := io.Pipe()
		writer := zip.NewWriter(pw)
		c.Header("Content-Type", "application/zip")
		c.Header(
			"Content-Disposition",
			fmt.Sprintf("attachment; filename=\"%s.zip\"", "test"),
		)
		c.Writer.Header().Del("Content-Length")
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			defer pw.Close()
			defer writer.Close()
			for i := 0; i < 100; i++ {
				filename := fmt.Sprintf("%d.txt", i)
				fmt.Println("start sending file", i)
				f, err := writer.Create(filename)
				if err != nil {
					log.Fatal(err)
				}
				readFile, err := os.Open(filename)
				if err != nil {
					log.Fatal(err)
				}
				buf := make([]byte, 1024)
				for {
					n, err := readFile.Read(buf)
					if err != nil {
						break
					}
					f.Write(buf[:n])
				}
			}
		}()

		go func() {
			defer wg.Done()
			for {
				dataRead := make([]byte, 1024)
				n, err := pr.Read(dataRead)
				if err != nil {
					return
				}
				c.Writer.Write(dataRead[:n])
			}
		}()
		wg.Wait()
	})
	router.Run()
}
