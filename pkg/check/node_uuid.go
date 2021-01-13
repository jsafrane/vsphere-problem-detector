package check

import (
	"fmt"

	"github.com/vmware/govmomi/vim25/mo"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

// CheckNodeDiskUUID makes sure that all nodes have disk.enableUUID=TRUE.
type CheckNodeDiskUUID struct{}

var _ NodeCheck = &CheckNodeDiskUUID{}

func (c *CheckNodeDiskUUID) Name() string {
	return "CheckNodeDiskUUID"
}

func (c *CheckNodeDiskUUID) StartCheck() error {
	return nil
}

func (c *CheckNodeDiskUUID) CheckNode(ctx *CheckContext, node *v1.Node, vm *mo.VirtualMachine) error {
	if vm.Config.Flags.DiskUuidEnabled == nil {
		return fmt.Errorf("the node has empty disk.enableUUID")
	}
	if *vm.Config.Flags.DiskUuidEnabled == false {
		return fmt.Errorf("the node has disk.enableUUID = FALSE")
	}
	klog.V(4).Infof("... the node has correct disk.enableUUID")

	if node.Labels != nil {
		if _, found := node.Labels["node-role.kubernetes.io/worker"]; found {
			klog.Info("Injecting dummy error")
			return fmt.Errorf("Mock error")
		}
		klog.Info("No dummy error")
	}
	return nil
}

func (c *CheckNodeDiskUUID) FinishCheck() {
	return
}
