package offline


import (
        "fmt"
        "log"
        "os"
	"os/exec"
	"bytes"
//	"path/filepath"
//	iso9660 "github.com/iso9660"
	"github.com/openshift/installer/pkg/asset"
)



type isoImageMetaData struct {
}

var _ asset.Asset = (*isoImageMetaData)(nil)

// Dependencies returns no dependencies.
func (a *isoImageMetaData) Dependencies() []asset.Asset {
	return nil 
}

// Generate queries for the cluster name from the user.
func (a *isoImageMetaData) Generate(parents asset.Parents) error {
        return nil
}

// Name returns the human-friendly name of the asset.
func (a *isoImageMetaData) Name() string {
	return "ISO Create  Disconnected"
}


func (a *isoImageMetaData) createISOImage(src string, fileName string) bool { 
	fmt.Println("Creating ISO Image for Disconnected")
	err := os.Chdir(src)

	if err != nil {
		log.Println(err)
	}

	outfile := "-o" + fileName

	command := exec.Command("/usr/bin/genisoimage", "-J", "-joliet-long", "-allow-limited-size", "-r", outfile, src)

	var out bytes.Buffer
	command.Stdout = &out
	err = command.Run()
	if err != nil {
		fmt.Println("Error occurred")
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	fmt.Printf("Command Output: %s\n", out.String())

	return true
}
	//genisoimage -J -joliet-long -quiet -r -o openshift-offline-installer.iso ./installer


// func (a *isoImageMetaData) createISOImage(src string, fileName string) bool { 
// 	fmt.Println("Creating ISO Image for Disconnected")

// 	// First let's create the ISO file
// 	writer, err := iso9660.NewWriter()
// 	if err != nil {
// 		log.Fatalf("failed to create writer: %s", err)
// 		return false
// 	}

// 	// Set the volume
// 	writer.Primary.VolumeIdentifier = "TISC Disconnected Volume"

// 	// Now let's walk through the directories and add the files
// 	err = os.Chdir(src)
// 	if err != nil {
// 		log.Println(err)
// 	} else {
		
// 		err = filepath.Walk(".",
// 			func(path string, info os.FileInfo, err error) error {
// 				if err != nil {
// 					return err
// 				}
// 				if !info.IsDir() {
// 					fmt.Printf("Adding file [%v] \n", path)
// 					err = writer.AddLocalFile(path, path)
// 					if err != nil {
// 						//log.Fatalf("failed to add file: %s", err)
// 						log.Printf("failed to add file [%s]: %s\n", path, err)
// 						//return err
// 					}
// 				} else {
// 					fmt.Printf("Skipping [%v] \n", path)
// 				}
// 				return nil
// 			})
// 		if err != nil {
// 			log.Println(err)
// 			return false
// 		}
// 	}
// 	fmt.Printf("Ready to write file [%v]\n", fileName)
// 	outputFile, err := os.OpenFile(fileName, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0644)
// 	if err != nil {
// 		log.Fatalf("failed to create file: %s", err)
// 	}
	
// 	err = writer.WriteTo(outputFile)
// 	if err != nil {
// 		log.Fatalf("failed to write ISO image: %s", err)
// 	}
	
// 	err = outputFile.Close()
// 	if err != nil {
// 		log.Fatalf("failed to close output file: %s", err)
// 	}
	
// 	return true
// }


