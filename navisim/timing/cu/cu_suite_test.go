package cu

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -write_package_comment=false -package=$GOPACKAGE -destination=mock_akita_test.go gitlab.com/akita/akita Port,Engine
//go:generate mockgen -write_package_comment=false -package=$GOPACKAGE -destination=mock_pipelining_test.go gitlab.com/akita/util/pipelining Pipeline
//go:generate mockgen -write_package_comment=false -package=$GOPACKAGE -destination=mock_util_test.go gitlab.com/akita/util Buffer
//go:generate mockgen -source subcomponent.go -destination mock_subcomponent_test.go -package $GOPACKAGE
//go:generate mockgen -source wfdispatcher.go -destination mock_wfdispatcher_test.go -package $GOPACKAGE

func TestSimulator(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "CU Suite")
}
