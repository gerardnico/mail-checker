package report

import (
	"github.com/kuberhealthy/kuberhealthy/v2/pkg/checks/external/checkclient"
	"log"
	"os"
)

func ToKuberHealthy(check MetaCheck) {
	// A kuberhealthy will fail only if it can't report it
	// Docs: https://github.com/kuberhealthy/kuberhealthy/blob/master/docs/CHECK_CREATION.md

	if os.Getenv("KH_REPORTING_URL") == "" {
		return
	}

	// Error
	if len(check.Errors) > 0 {

		err := checkclient.ReportFailure(check.Errors)
		if err != nil {
			log.Fatal("Unable to report the following errors to kuberhealthy", check.Errors)
		}

	}

	// Success
	err := checkclient.ReportSuccess()
	if err != nil {
		log.Fatal("Could not report kuberhealthy success")
	}

}
