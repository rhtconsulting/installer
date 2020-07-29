package offline


import (
        "fmt"
        "bytes"
        "log"
        "os/exec"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"

)



type mirrorReleaseMetaData struct {
}

var _ asset.Asset = (*mirrorReleaseMetaData)(nil)

// Dependencies returns no dependencies.
func (a *mirrorReleaseMetaData) Dependencies() []asset.Asset {
	return nil 
}

// Generate queries for the cluster name from the user.
func (a *mirrorReleaseMetaData) Generate(parents asset.Parents) error {
        return nil
}

// Name returns the human-friendly name of the asset.
func (a *mirrorReleaseMetaData) Name() string {
	return "Mirror Release Disconnected"
}


func (a *mirrorReleaseMetaData) pullMirrorImages(config *types.OfflineConfig) bool {
	fmt.Println("Pulling mirror images for OCP 4 mirror")

	source := "--from=" + config.Ocpmirror.Src
	dest   := "--to-dir=" + config.Ocpmirror.Dest

	command := exec.Command(config.Ocpmirror.Ocbin, "adm", "-a", config.Ocpmirror.Pullsecret, "release", "mirror", source, "--to=file://openshift/release", dest)

	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		fmt.Println("Error occurred")
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	fmt.Printf("Command Output: %s\n", out.String())
	return true
}

func (a *mirrorReleaseMetaData) extractInstaller(config *types.OfflineConfig) bool {
	fmt.Println("Extracting OCP 4 installer from mirror images")

	source := "--from=" + config.Ocpmirror.Src
	dest   := "--to=" + config.Ocpmirror.Dest

	command := exec.Command(config.Ocpmirror.Ocbin, "adm", "-a", config.Ocpmirror.Pullsecret, "release", "extract", "--command=openshift-install", source, dest)

	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		fmt.Println("Error occurred")
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	fmt.Printf("Command Output: %s\n", out.String())
	fmt.Println("extractInstaller called!")
	return true
}
