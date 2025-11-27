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
	if hasOSUpdates(params) {
		vm.SetOs(buildOSForUpdate(params))
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

func hasOSUpdates(params UpdateVMParameters) bool {
	return len(params.BootDevices()) > 0 || params.Cmdline() != nil ||
		params.CustomKernelCmdline() != nil || params.Initrd() != nil || params.Kernel() != nil
}

func buildOSForUpdate(params UpdateVMParameters) *ovirtsdk.OperatingSystem {
	osBuilder := ovirtsdk.NewOperatingSystemBuilder()

	if bootDevices := params.BootDevices(); len(bootDevices) > 0 {
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

	return osBuilder.MustBuild()
}

func (m *mockClient) UpdateVM(id VMID, params UpdateVMParameters, _ ...RetryStrategy) (VM, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.vms[id]; !ok {
		return nil, newError(ENotFound, "VM with ID %s not found", id)
	}

	vm := m.vms[id]
	vm = m.updateVMBasicFields(vm, params)
	vm = m.updateVMKernelParams(vm, params)

	m.vms[id] = vm
	return vm, nil
}

func (m *mockClient) updateVMBasicFields(vm *vm, params UpdateVMParameters) *vm {
	if name := params.Name(); name != nil {
		for _, otherVM := range m.vms {
			if otherVM.name == *name && otherVM.ID() != vm.ID() {
				return vm
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
	return vm
}

func (m *mockClient) updateVMKernelParams(vm *vm, params UpdateVMParameters) *vm {
	// Update boot devices
	if bootDevices := params.BootDevices(); len(bootDevices) > 0 {
		vm = m.updateVMBootDevices(vm, bootDevices)
	}

	// Update kernel parameters
	if cmdline := params.Cmdline(); cmdline != nil {
		vm = m.updateVMOSField(vm, func(os *vmOS) { os.cmdline = cmdline })
	}
	if customKernelCmdline := params.CustomKernelCmdline(); customKernelCmdline != nil {
		vm = m.updateVMOSField(vm, func(os *vmOS) { os.customKernelCmdline = customKernelCmdline })
	}
	if initrd := params.Initrd(); initrd != nil {
		vm = m.updateVMOSField(vm, func(os *vmOS) { os.initrd = initrd })
	}
	if kernel := params.Kernel(); kernel != nil {
		vm = m.updateVMOSField(vm, func(os *vmOS) { os.kernel = kernel })
	}

	return vm
}

func (m *mockClient) updateVMBootDevices(vm *vm, bootDevices []BootDevice) *vm {
	if vm.os == nil {
		vm.os = &vmOS{bootDevices: bootDevices}
	} else {
		newOS := *vm.os
		newOS.bootDevices = bootDevices
		vm.os = &newOS
	}
	return vm
}

func (m *mockClient) updateVMOSField(vm *vm, updateFunc func(*vmOS)) *vm {
	if vm.os == nil {
		vm.os = &vmOS{}
	}
	newOS := *vm.os
	updateFunc(&newOS)
	vm.os = &newOS
	return vm
}
