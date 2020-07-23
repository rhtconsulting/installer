package offline


import (
        "fmt"
        "bytes"
        "log"
        "os/exec"

	"github.com/openshift/installer/pkg/asset"
//	"github.com/openshift/installer/pkg/types"

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


func (a *mirrorReleaseMetaData) pullMirrorImages() bool { //config *types.OfflineConfig) bool {
	fmt.Println("Pulling mirror images for OCP 4 mirror")

	command := exec.Command("/usr/local/bin/oc", "adm", "-a", "/opt/registry/auth/local-secret.txt", "release", "mirror", "--from=quay.io/openshift-release-dev/ocp-release:4.4.0-x86_64", "--to=file://openshift/release", "--to-dir=/ocp-images")
	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		fmt.Println("Error occurred")
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	fmt.Printf("Command Output: %s\n", out.String())
	fmt.Println("pullMirrorImages called!")
	return true
}

func (a *mirrorReleaseMetaData) extractInstaller() bool { //config *types.OfflineConfig) bool {
	fmt.Println("Extracting OCP 4 installer from mirror images")

	command := exec.Command("/usr/local/bin/oc", "adm", "-a", "/opt/registry/auth/local-secret.txt", "release", "extract", "--command=openshift-install", "--from=quay.io/openshift-release-dev/ocp-release:4.4.0-x86_64", "--to=/ocp-images")
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

