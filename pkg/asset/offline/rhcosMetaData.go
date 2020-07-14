package offline

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/validate"
)

type rhcosReleaseMetaData struct {
	distributionURL string
}

var _ asset.Asset = (*rhcosReleaseMetaData)(nil)

// Dependencies returns no dependencies.
func (a *rhcosReleaseMetaData) Dependencies() []asset.Asset {
	return []asset.Asset{
		//&platform{},
	}
}

// Generate queries for the cluster name from the user.
func (a *rhcosReleaseMetaData) Generate(parents asset.Parents) error {

	validator := survey.Required

	validator = survey.ComposeValidators(validator, func(ans interface{}) error {
		return validate.URI(ans.(string))
	})

	return survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "URL for RHCOS Release",
				Help:    "The URL link to download the RHCOS release. This will be used to download the RHCOS release image for the platform. \n\nFor AWS, the URL  http://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/latest/latest/rhcos-latest-x86_64-aws-x86_64.vmdk.gz .",
			},
			Validate: validator,
		},
	}, &a.distributionURL)
}

// Name returns the human-friendly name of the asset.
func (a *rhcosReleaseMetaData) Name() string {
	return "RHCOS Release Metadata"
}
