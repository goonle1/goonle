package main

import (
	"github.com/golang/glog"
	"os"
	"io/ioutil"
	"path/filepath"
	dp "dp/dpds"
)


// This function transfers a filesystem from
// the local FS into the file consumer.  Depending
// on the implementation of the consumer, the files
// can go anywhere from there.
func transferFilesystem(root string) error {
	consumeErrors := make(chan error, 1)
	nameDotMap  := make(map[string]*dp.Dot)
	consumer := dp.GetConsumerInstance()
	var i uint64 = 0

	// No select needed for this send, since errc is buffered.
    consumeErrors <- filepath.Walk(root, func(currentPath string, fInfo os.FileInfo, err error) error {
            if err != nil {
                return err
            }

            dot := &dp.Dot{i, 0, fInfo.Name(), ""}
            i++
            cleanPath := filepath.Clean(currentPath)
            dir, _ := filepath.Rel(root, filepath.Dir(cleanPath))
            parent := nameDotMap[dir]
            if parent != nil {
            	dot.ParentId = parent.Id
            } else {
            	dot.ParentId = 0
            }
            
            if fInfo.IsDir() {
            	consumer.Prepare()
                nameDotMap[dot.Name] = dot
			    go func(dot *dp.Dot) {
			    	consumer.Consume(dot.Id, dot.ParentId, dot.Name, "")
			    } (dot)
                return nil
            }

            if !fInfo.Mode().IsRegular() {
                return nil
            }
            consumer.Prepare()

			go func(cleanPath string, dot *dp.Dot) {
	            byteArrayValue, err := ioutil.ReadFile(cleanPath)
	            if err != nil {
	            	// Cleanup and abort..
	            	consumer.Abort()
	            	return
	            }
	
	            if byteArrayValue != nil {
	            	dot.Value = string(byteArrayValue[:])
	            }
		        consumer.Consume(dot.Id, dot.ParentId, dot.Name, dot.Value)
	        }(cleanPath, dot)

            return nil
        })
	
	consumer.Commit()

    return nil;
}

func main() {
	glog.Info("Start to dottify.")
	
	dir, err := os.Getwd();
	
	if err != nil {
		glog.Info("Current working directory is no good.")
	}

	transferFilesystem(dir)
	
	glog.Info("End dottify.")
}

