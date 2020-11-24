package check

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
	"k8s.io/klog/v2"
	"k8s.io/legacy-cloud-providers/vsphere"
)

// CheckFolderList tests that OCP has permissions to list volumes in Datastore.
// This is necessary to create volumes.
func CheckTaskPermissions(pctx context.Context, vmConfig *vsphere.VSphereConfig, vmClient *vim25.Client, kubeClient KubeClient) error {
	klog.V(4).Infof("CheckTaskPermissions started")

	ctx, cancel := context.WithTimeout(pctx, *Timeout)
	defer cancel()

	mgr := view.NewManager(vmClient)
	view, err := mgr.CreateTaskView(ctx, vmClient.ServiceContent.TaskManager)
	if err != nil {
		return fmt.Errorf("error creating task view: %s", err)
	}

	taskCount := 0
	ctx, cancel = context.WithTimeout(pctx, *Timeout)
	defer cancel()
	err = view.Collect(ctx, func(tasks []types.TaskInfo) {
		for _, task := range tasks {
			klog.V(4).Infof("Found task %s", task.Name)
			taskCount++
		}
	})

	if err != nil {
		return fmt.Errorf("error collecting tasks: %s", err)
	}
	klog.Infof("CheckTaskPermissions succeeded, %d tasks found", taskCount)
	return nil
}
