package services

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/assets"
)

var _ = Describe("Service Broker Lifecycle", func() {
	var broker ServiceBroker

	BeforeEach(func() {
		broker = NewServiceBroker(generator.RandomName(), assets.NewAssets().ServiceBroker, context)
	})

	It("confirms correct behavior in the lifecycle of a service broker", func() {
		broker.Push()
		broker.Configure()

		// Adding the service broker
		broker.Create()

		// Confirming the plans are not yet public
		plans := cf.Cf("marketplace").Wait(DEFAULT_TIMEOUT)
		Expect(plans).To(Exit(0))
		output := plans.Out.Contents()
		Expect(output).NotTo(ContainSubstring(broker.Service.Name))
		Expect(output).NotTo(ContainSubstring(broker.Plans()[0].Name))

		broker.PublicizePlans()

		// Confirming plans show up in the marketplace
		plans = cf.Cf("marketplace").Wait(DEFAULT_TIMEOUT)
		Expect(plans).To(Exit(0))
		output = plans.Out.Contents()
		Expect(output).To(ContainSubstring(broker.Service.Name))
		Expect(output).To(ContainSubstring(broker.Plans()[0].Name))

		// Changing the catalog on the broker
		oldServiceName := broker.Service.Name
		oldPlanName := broker.SyncPlans[0].Name
		broker.Service.Name = generator.RandomName()
		broker.SyncPlans[0].Name = generator.RandomName()
		broker.Configure()
		broker.Update()

		// Confirming the changes to the broker show up in the marketplace
		plans = cf.Cf("marketplace").Wait(DEFAULT_TIMEOUT)
		Expect(plans).To(Exit(0))
		output = plans.Out.Contents()
		Expect(output).NotTo(ContainSubstring(oldServiceName))
		Expect(output).NotTo(ContainSubstring(oldPlanName))
		Expect(output).To(ContainSubstring(broker.Service.Name))
		Expect(output).To(ContainSubstring(broker.Plans()[0].Name))

		// Deleting the service broker and confirming the plans no longer display
		broker.Delete()
		plans = cf.Cf("marketplace").Wait(DEFAULT_TIMEOUT)
		Expect(plans).To(Exit(0))
		output = plans.Out.Contents()
		Expect(output).NotTo(ContainSubstring(oldServiceName))
		Expect(output).NotTo(ContainSubstring(oldPlanName))
		Expect(output).NotTo(ContainSubstring(broker.Service.Name))
		Expect(output).NotTo(ContainSubstring(broker.Plans()[0].Name))

		broker.Destroy()
	})
})
