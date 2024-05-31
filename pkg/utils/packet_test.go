package utils

import (
	"time"

	mocks_utils "github.com/k8snetworkplumbingwg/sriov-cni/pkg/utils/mocks"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)



var _ = ginkgo.Describe("Packets", func() {

	ginkgo.Context("WaitForCarrier", func() {
		ginkgo.It("should wait until the link has IFF_UP flag", func() {
			ginkgo.DeferCleanup(func(old NetlinkManager) { netLinkLib = old }, netLinkLib)
			
			mockedNetLink := &mocks_utils.NetlinkManager{}
			netLinkLib = mockedNetLink
			
			fakeLink := &FakeLink{LinkAttrs: netlink.LinkAttrs{
				Index:        1000,
				Name:         "dummylink",
				RawFlags: 0,
			}}
		
			mockedNetLink.On("LinkByName", "dummylink").Return(fakeLink, nil)
		
			hasCarrier := make(chan bool)
			go func() {
				hasCarrier <- WaitForCarrier("dummylink", 5*time.Second)
			}()

			gomega.Consistently(hasCarrier, "100ms").ShouldNot(gomega.Receive())

			fakeLink.RawFlags |= unix.IFF_UP

			gomega.Eventually(hasCarrier, "300ms").Should(gomega.Receive())
		})
	})
})
