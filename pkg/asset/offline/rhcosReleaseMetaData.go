package offline

import (
	survey "gopkg.in/AlecAivazis/survey.v1"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"io"
	"path"
	"log"
	"time"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/validate"
)


type rhcosReleaseMetaData struct {
	distributionURL string      
}

var _ asset.Asset = (*rhcosReleaseMetaData)(nil)

// Dependencies returns no dependencies.
func (a *rhcosReleaseMetaData) Dependencies() []asset.Asset {
        //panic("Dependencies")
	return nil 
}

func (a *rhcosReleaseMetaData) CreateOfflinePackage(src string) string {

	fmt.Printf("Downloading %s\n", src)

	res, err := http.Get(src)
	check(err)
	defer res.Body.Close()

	u, err := url.Parse(src)
	check(err)

	out, err := os.Create(path.Base(u.Path))
	defer out.Close()

	size := res.ContentLength
	bar := &Progbar{total: int(size)}
	written := make(chan int, 500)

	go func() {
		copied := 0
		c := 0
		tick := time.Tick(interval)

		for {
			select {
			case c = <-written:
				copied += c
			case <-tick:
				bar.PrintProg(copied)
			}
		}
	}()

	buf := make([]byte, 32*1024)
	for {
		rc, re := res.Body.Read(buf)
		if rc > 0 {
			wc, we := out.Write(buf[0:rc])
			check(we)

			if wc != rc {
				log.Fatal("Read and Write count mismatch")
			}

			if wc > 0 {
				written <- wc
			}
		}
		if re == io.EOF {
			break
		}
		check(re)
	}
	bar.PrintComplete()




	return ""
}


// Generate queries for the cluster name from the user.
func (a *rhcosReleaseMetaData) Generate(parents asset.Parents) error {
	validator := survey.Required

	validator = survey.ComposeValidators(validator, func(ans interface{}) error {
		return validate.URI(ans.(string))
	})

        return nil

}

// Name returns the human-friendly name of the asset.
func (a *rhcosReleaseMetaData) Name() string {
	fmt.Println("RHCOS Name")
	return "rhcosReleaseMetaData"
}

