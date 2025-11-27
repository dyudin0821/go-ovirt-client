package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) UpdateVM(
	id VMID,
	params UpdateVMParameters,
	retries ...RetryStrategy,
) (result VM, err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))

	vm := &ovirtsdk.Vm{}
	vm.SetId(string(id))
	if name := params.Name(); name != nil {
		if *name == "" {
			return nil, newError(EBadArgument, "name must not be empty for VM update")
		}
		vm.SetName(*name)
	}
	if comment := params.Comment(); comment != nil {
		vm.SetComment(*comment)
	}
	if description := params.Description(); description != nil {
		vm.SetDescription(*description)
	}

	// Handle OS parameters including boot devices and kernel parameters
	if bootDevices := params.BootDevices(); len(bootDevices) > 0 || params.Cmdline() != nil ||
		params.CustomKernelCmdline() != nil || params.Initrd() != nil || params.Kernel() != nil {
		osBuilder := ovirtsdk.NewOperatingSystemBuilder()

		if len(bootDevices) > 0 {
			sdkBootDevices := make([]ovirtsdk.BootDevice, len(bootDevices))
			for i, device := range bootDevices {
				sdkBootDevices[i] = ovirtsdk.BootDevice(device)
			}
			bootBuilder := ovirtsdk.NewBootBuilder().DevicesOfAny(sdkBootDevices...)
			osBuilder.BootBuilder(bootBuilder)
		}

		// Set kernel parameters
		if cmdline := params.Cmdline(); cmdline != nil {
			osBuilder.Cmdline(*cmdline)
		}
		if customKernelCmdline := params.CustomKernelCmdline(); customKernelCmdline != nil {
			osBuilder.CustomKernelCmdline(*customKernelCmdline)
		}
		if initrd := params.Initrd(); initrd != nil {
			osBuilder.Initrd(*initrd)
		}
		if kernel := params.Kernel(); kernel != nil {
			osBuilder.Kernel(*kernel)
		}

		vm.SetOs(osBuilder.MustBuild())
	}

	err = retry(
		fmt.Sprintf("updating vm %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().VmsService().VmService(string(id)).Update().Vm(vm).Send()
			if err != nil {
				return wrap(err, EUnidentified, "failed to update VM")
			}
			vm, ok := response.Vm()
			if !ok {
				return newError(EFieldMissing, "missing VM in VM update response")
			}
			result, err = convertSDKVM(vm, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert VM",
				)
			}
			return nil
		})
	return result, err
}

func (m *mockClient) UpdateVM(id VMID, params UpdateVMParameters, _ ...RetryStrategy) (VM, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.vms[id]; !ok {
		return nil, newError(ENotFound, "VM with ID %s not found", id)
	}

	vm := m.vms[id]
	if name := params.Name(); name != nil {
		for _, otherVM := range m.vms {
			if otherVM.name == *name && otherVM.ID() != vm.ID() {
				return nil, newError(EConflict, "A VM with the name \"%s\" already exists.", *name)
			}
		}
		vm = vm.withName(*name)
	}
	if comment := params.Comment(); comment != nil {
		vm = vm.withComment(*comment)
	}
	if description := params.Description(); description != nil {
		vm = vm.withDescription(*description)
	}

	// Update OS parameters including boot devices and kernel parameters
	if bootDevices := params.BootDevices(); len(bootDevices) > 0 {
		if vm.os == nil {
			vm.os = &vmOS{
				bootDevices: bootDevices,
			}
		} else {
			newOS := *vm.os
			newOS.bootDevices = bootDevices
			vm.os = &newOS
		}
	}

	// Update kernel parameters
	if cmdline := params.Cmdline(); cmdline != nil {
		if vm.os == nil {
			vm.os = &vmOS{}
		}
		newOS := *vm.os
		newOS.cmdline = cmdline
		vm.os = &newOS
	}
	if customKernelCmdline := params.CustomKernelCmdline(); customKernelCmdline != nil {
		if vm.os == nil {
			vm.os = &vmOS{}
		}
		newOS := *vm.os
		newOS.customKernelCmdline = customKernelCmdline
		vm.os = &newOS
	}
	if initrd := params.Initrd(); initrd != nil {
		if vm.os == nil {
			vm.os = &vmOS{}
		}
		newOS := *vm.os
		newOS.initrd = initrd
		vm.os = &newOS
	}
	if kernel := params.Kernel(); kernel != nil {
		if vm.os == nil {
			vm.os = &vmOS{}
		}
		newOS := *vm.os
		newOS.kernel = kernel
		vm.os = &newOS
	}

	m.vms[id] = vm

	return vm, nil
}
