package knuu

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

// Destroy destroys the instance
// This function can only be called in the state 'Started' or 'Destroyed'
func (i *Instance) Destroy() error {
	if i.state == Destroyed {
		return nil
	}

	if !i.IsInState(Started, Stopped, Destroyed) {
		return ErrDestroyingNotAllowed.WithParams(i.state.String())
	}

	// TODO: receive context from the user in the breaking refactor
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := i.destroyPod(ctx); err != nil {
		return ErrDestroyingPod.WithParams(i.k8sName).Wrap(err)
	}
	if err := i.destroyResources(ctx); err != nil {
		return ErrDestroyingResourcesForInstance.WithParams(i.k8sName).Wrap(err)
	}

	err := applyFunctionToInstances(i.sidecars, func(sidecar Instance) error {
		logrus.Debugf("Destroying sidecar resources from '%s'", sidecar.k8sName)
		return sidecar.destroyResources(ctx)
	})
	if err != nil {
		return ErrDestroyingResourcesForSidecars.WithParams(i.k8sName).Wrap(err)
	}

	i.state = Destroyed
	setStateForSidecars(i.sidecars, Destroyed)
	logrus.Debugf("Set state of instance '%s' to '%s'", i.k8sName, i.state.String())

	return nil
}

// BatchDestroy destroys a list of instances.
func BatchDestroy(instances ...*Instance) error {
	if os.Getenv("KNUU_SKIP_CLEANUP") == "true" {
		logrus.Info("Skipping cleanup")
		return nil
	}

	for _, instance := range instances {
		if instance == nil {
			continue
		}
		if err := instance.Destroy(); err != nil {
			return err
		}
	}
	return nil
}
