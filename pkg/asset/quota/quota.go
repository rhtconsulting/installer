package quota

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/quota/aws"
	"github.com/openshift/installer/pkg/asset/quota/gcp"
	"github.com/openshift/installer/pkg/diagnostics"
	"github.com/openshift/installer/pkg/quota"
	quotaaws "github.com/openshift/installer/pkg/quota/aws"
	quotagcp "github.com/openshift/installer/pkg/quota/gcp"
	typesaws "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	typesgcp "github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// PlatformQuotaCheck is an asset that validates the install-config platform for
// any resource requirements based on the quotas available.
type PlatformQuotaCheck struct {
}

var _ asset.Asset = (*PlatformQuotaCheck)(nil)

// Dependencies returns the dependencies for PlatformQuotaCheck
func (a *PlatformQuotaCheck) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&machines.Master{},
		&machines.Worker{},
	}
}

// Generate queries for input from the user.
func (a *PlatformQuotaCheck) Generate(dependencies asset.Parents) error {
	ic := &installconfig.InstallConfig{}
	mastersAsset := &machines.Master{}
	workersAsset := &machines.Worker{}
	dependencies.Get(ic, mastersAsset, workersAsset)

	masters, err := mastersAsset.Machines()
	if err != nil {
		return err
	}

	workers, err := workersAsset.MachineSets()
	if err != nil {
		return err
	}

	platform := ic.Config.Platform.Name()
	switch platform {
	case typesaws.Name:
		services := []string{"ec2", "vpc"}
		session, err := ic.AWS.Session(context.TODO())
		if err != nil {
			return errors.Wrap(err, "failed to load AWS session")
		}
		q, err := quotaaws.Load(context.TODO(), session, ic.AWS.Region, services...)
		if quotaaws.IsUnauthorized(err) {
			logrus.Warnf("Missing permissions to fetch Quotas and therefore will skip checking them: %v, make sure you have `servicequotas:ListAWSDefaultServiceQuotas` persmission available to the user.", err)
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "failed to load Quota for services: %s", strings.Join(services, ", "))
		}
		instanceTypes, err := aws.InstanceTypes(context.TODO(), session, ic.AWS.Region)
		if quotaaws.IsUnauthorized(err) {
			logrus.Warnf("Missing permissions to fetch instance types and therefore will skip checking Quotas: %v, make sure you have `ec2:DescribeInstanceTypes` persmission available to the user.", err)
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "failed to load instance types for %s", ic.AWS.Region)
		}
		reports, err := quota.Check(q, aws.Constraints(ic.Config, masters, workers, instanceTypes))
		if err != nil {
			return summarizeFailingReport(reports)
		}
		summarizeReport(reports)
	case typesgcp.Name:
		services := []string{"compute.googleapis.com", "iam.googleapis.com"}
		q, err := quotagcp.Load(context.TODO(), ic.Config.Platform.GCP.ProjectID, services...)
		if quotagcp.IsUnauthorized(err) {
			logrus.Warnf("Missing permissions to fetch Quotas and therefore will skip checking them: %v, make sure you have `roles/servicemanagement.quotaViewer` assigned to the user.", err)
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "failed to load Quota for services: %s", strings.Join(services, ", "))
		}
		reports, err := quota.Check(q, gcp.Constraints(ic.Config, masters, workers))
		if err != nil {
			return summarizeFailingReport(reports)
		}
		summarizeReport(reports)
	case azure.Name, baremetal.Name, libvirt.Name, none.Name, openstack.Name, ovirt.Name, vsphere.Name:
		// no special provisioning requirements to check
	default:
		err = fmt.Errorf("unknown platform type %q", platform)
	}
	return err
}

// Name returns the human-friendly name of the asset.
func (a *PlatformQuotaCheck) Name() string {
	return "Platform Quota Check"
}

// summarizeFailingReport summarizes a report when there are failing constraints.
func summarizeFailingReport(reports []quota.ConstraintReport) error {
	var notavailable []string
	var unknown []string
	for _, report := range reports {
		switch report.Result {
		case quota.NotAvailable:
			notavailable = append(notavailable, fmt.Sprintf("%s is not available in %s because %s", report.For.Name, report.For.Region, report.Message))
		case quota.Unknown:
			unknown = append(unknown, report.For.Name)
		default:
			continue
		}
	}

	if len(notavailable) == 0 && len(unknown) > 0 {
		// all quotas are missing information so warn and skip
		logrus.Warnf("Failed to find information on quotas %s", strings.Join(unknown, ", "))
		return nil
	}

	msg := strings.Join(notavailable, ", ")
	if len(unknown) > 0 {
		msg = fmt.Sprintf("%s, and could not find information on %s", msg, strings.Join(unknown, ", "))
	}
	return &diagnostics.Err{Reason: "MissingQuota", Message: msg}
}

// summarizeReport summarizes a report when there are availble.
func summarizeReport(reports []quota.ConstraintReport) {
	var low []string
	for _, report := range reports {
		switch report.Result {
		case quota.AvailableButLow:
			low = append(low, fmt.Sprintf("%s (%s)", report.For.Name, report.For.Region))
		default:
			continue
		}
	}
	if len(low) > 0 {
		logrus.Warnf("Following quotas %s are available but will be completely used pretty soon.", strings.Join(low, ", "))
	}
}
