package main

import (
    "os"
    "fmt"
    "archive/zip"
    "io"
    "path/filepath"
    "strings"
)

func zipIt(source, target string) error {
    // create an empty file
    zipfile, err := os.Create(target)


    if err != nil {
        return err
    }

    defer zipfile.Close()

    // create something which can create zip files
    archive := zip.NewWriter(zipfile)
    defer archive.Close()

    // get info about the source file
    info, err := os.Stat(source)
    
    if err != nil {
        return err
    }

    var baseDir string

    // source is a directory then the zip file being created will have source as its root
    if info.IsDir() {
        baseDir = filepath.Base(source)
    }

    // walk through source and all its descendents, zipping each one
    filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // every zip file has a header with info that lets it be unzipped
        header, err := zip.FileInfoHeader(info)
        if err != nil {
            return err
        }

        // if source is a directory then append source before the path of the current file
        if baseDir != "" {
            header.Name = filepath.Join(baseDir,strings.TrimPrefix(path, source))
        }

        // if the current file is a directory, it needs a slash at the end
        if info.IsDir() {
            header.Name += "/"
        } else {
            // if the current file is not a directory, the algorithm needs to know to deflate it
            header.Method = zip.Deflate
        }

        writer, err := archive.CreateHeader(header)
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        // get the current file
        file, err := os.Open(path)
        if err != nil {
            return err
        }

        defer file.Close()

        // add the current file to the zip writer
        _, err = io.Copy(writer, file)
        return err;;;;
    })
    return err  
}



func unZipIt(source, target string) error {
    reader, err := zip.OpenReader(source)
    if err != nil {
        return err
    }

    if err := os.MkdirAll(target, 0755); err != nil {
        return err
    }

    for _, file := range reader.File {
        path := filepath.Join(target, file.Name)
        if file.FileInfo().IsDir() {
            os.MkdirAll(path, file.Mode())
            continue
        }

        fileReader, err := file.Open()
        if err != nil {
            if fileReader != nil {
                fileReader.Close()
            }
            return err
        }

        targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
        if err != nil {
    fileReader.Close()

            if targetFile != nil {
                targetFile.Close()
            }

            return err
        }

        if _, err := io.Copy(targetFile, fileReader); err != nil {
            fileReader.Close()
            targetFile.Close()

            return err
        }

        fileReader.Close()
        targetFile.Close()
    }

    return nil
}


func main() {
    args := os.Args
    if len(args) < 2 {
        fmt.Printf("Please supply an argument, either zip or unzip\n")
        return
    }

    var source, dest string
    switch(os.Args[1]) {
        case "unzip":
            fmt.Printf("What file would you like to unzip?\n")
            if _, err := fmt.Scanf("%s", &source); err != nil {
                return
            } else {
            fmt.Printf("What is the destination name?\n")

            if _, err := fmt.Scanf("%s", &dest); err != nil {
                    return
                } else {
                
                    if err := unZipIt(source, dest); err != nil {
                                fmt.Printf("Error %s", err)
                            }
                        }
                    }
        case "zip":
        fmt.Printf("What file would you like to zip?\n")
            if _, err := fmt.Scanf("%s", &source); err != nil {
                return
            } else {
            fmt.Printf("What is the destination name?\n")

            if _, err := fmt.Scanf("%s", &dest); err != nil {
                    return
                } else {
                
                    if err := zipIt(source, dest); err != nil {
                                fmt.Printf("Error %s", err)
                            }
                        }
                    }
    
        default: 
            fmt.Println("Usage -- compressor zip|unzip")
    }
   
        

}