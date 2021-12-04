// Code generated by counterfeiter. DO NOT EDIT.
package internalfakes

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/elgohr/go-localstack/internal"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type FakeDockerClient struct {
	ContainerCreateStub        func(context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, *v1.Platform, string) (container.ContainerCreateCreatedBody, error)
	containerCreateMutex       sync.RWMutex
	containerCreateArgsForCall []struct {
		arg1 context.Context
		arg2 *container.Config
		arg3 *container.HostConfig
		arg4 *network.NetworkingConfig
		arg5 *v1.Platform
		arg6 string
	}
	containerCreateReturns struct {
		result1 container.ContainerCreateCreatedBody
		result2 error
	}
	containerCreateReturnsOnCall map[int]struct {
		result1 container.ContainerCreateCreatedBody
		result2 error
	}
	ContainerInspectStub        func(context.Context, string) (types.ContainerJSON, error)
	containerInspectMutex       sync.RWMutex
	containerInspectArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	containerInspectReturns struct {
		result1 types.ContainerJSON
		result2 error
	}
	containerInspectReturnsOnCall map[int]struct {
		result1 types.ContainerJSON
		result2 error
	}
	ContainerStartStub        func(context.Context, string, types.ContainerStartOptions) error
	containerStartMutex       sync.RWMutex
	containerStartArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 types.ContainerStartOptions
	}
	containerStartReturns struct {
		result1 error
	}
	containerStartReturnsOnCall map[int]struct {
		result1 error
	}
	ContainerStopStub        func(context.Context, string, *time.Duration) error
	containerStopMutex       sync.RWMutex
	containerStopArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 *time.Duration
	}
	containerStopReturns struct {
		result1 error
	}
	containerStopReturnsOnCall map[int]struct {
		result1 error
	}
	ImageListStub        func(context.Context, types.ImageListOptions) ([]types.ImageSummary, error)
	imageListMutex       sync.RWMutex
	imageListArgsForCall []struct {
		arg1 context.Context
		arg2 types.ImageListOptions
	}
	imageListReturns struct {
		result1 []types.ImageSummary
		result2 error
	}
	imageListReturnsOnCall map[int]struct {
		result1 []types.ImageSummary
		result2 error
	}
	ImagePullStub        func(context.Context, string, types.ImagePullOptions) (io.ReadCloser, error)
	imagePullMutex       sync.RWMutex
	imagePullArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 types.ImagePullOptions
	}
	imagePullReturns struct {
		result1 io.ReadCloser
		result2 error
	}
	imagePullReturnsOnCall map[int]struct {
		result1 io.ReadCloser
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDockerClient) ContainerCreate(arg1 context.Context, arg2 *container.Config, arg3 *container.HostConfig, arg4 *network.NetworkingConfig, arg5 *v1.Platform, arg6 string) (container.ContainerCreateCreatedBody, error) {
	fake.containerCreateMutex.Lock()
	ret, specificReturn := fake.containerCreateReturnsOnCall[len(fake.containerCreateArgsForCall)]
	fake.containerCreateArgsForCall = append(fake.containerCreateArgsForCall, struct {
		arg1 context.Context
		arg2 *container.Config
		arg3 *container.HostConfig
		arg4 *network.NetworkingConfig
		arg5 *v1.Platform
		arg6 string
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	stub := fake.ContainerCreateStub
	fakeReturns := fake.containerCreateReturns
	fake.recordInvocation("ContainerCreate", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.containerCreateMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDockerClient) ContainerCreateCallCount() int {
	fake.containerCreateMutex.RLock()
	defer fake.containerCreateMutex.RUnlock()
	return len(fake.containerCreateArgsForCall)
}

func (fake *FakeDockerClient) ContainerCreateCalls(stub func(context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, *v1.Platform, string) (container.ContainerCreateCreatedBody, error)) {
	fake.containerCreateMutex.Lock()
	defer fake.containerCreateMutex.Unlock()
	fake.ContainerCreateStub = stub
}

func (fake *FakeDockerClient) ContainerCreateArgsForCall(i int) (context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, *v1.Platform, string) {
	fake.containerCreateMutex.RLock()
	defer fake.containerCreateMutex.RUnlock()
	argsForCall := fake.containerCreateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6
}

func (fake *FakeDockerClient) ContainerCreateReturns(result1 container.ContainerCreateCreatedBody, result2 error) {
	fake.containerCreateMutex.Lock()
	defer fake.containerCreateMutex.Unlock()
	fake.ContainerCreateStub = nil
	fake.containerCreateReturns = struct {
		result1 container.ContainerCreateCreatedBody
		result2 error
	}{result1, result2}
}

func (fake *FakeDockerClient) ContainerCreateReturnsOnCall(i int, result1 container.ContainerCreateCreatedBody, result2 error) {
	fake.containerCreateMutex.Lock()
	defer fake.containerCreateMutex.Unlock()
	fake.ContainerCreateStub = nil
	if fake.containerCreateReturnsOnCall == nil {
		fake.containerCreateReturnsOnCall = make(map[int]struct {
			result1 container.ContainerCreateCreatedBody
			result2 error
		})
	}
	fake.containerCreateReturnsOnCall[i] = struct {
		result1 container.ContainerCreateCreatedBody
		result2 error
	}{result1, result2}
}

func (fake *FakeDockerClient) ContainerInspect(arg1 context.Context, arg2 string) (types.ContainerJSON, error) {
	fake.containerInspectMutex.Lock()
	ret, specificReturn := fake.containerInspectReturnsOnCall[len(fake.containerInspectArgsForCall)]
	fake.containerInspectArgsForCall = append(fake.containerInspectArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.ContainerInspectStub
	fakeReturns := fake.containerInspectReturns
	fake.recordInvocation("ContainerInspect", []interface{}{arg1, arg2})
	fake.containerInspectMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDockerClient) ContainerInspectCallCount() int {
	fake.containerInspectMutex.RLock()
	defer fake.containerInspectMutex.RUnlock()
	return len(fake.containerInspectArgsForCall)
}

func (fake *FakeDockerClient) ContainerInspectCalls(stub func(context.Context, string) (types.ContainerJSON, error)) {
	fake.containerInspectMutex.Lock()
	defer fake.containerInspectMutex.Unlock()
	fake.ContainerInspectStub = stub
}

func (fake *FakeDockerClient) ContainerInspectArgsForCall(i int) (context.Context, string) {
	fake.containerInspectMutex.RLock()
	defer fake.containerInspectMutex.RUnlock()
	argsForCall := fake.containerInspectArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeDockerClient) ContainerInspectReturns(result1 types.ContainerJSON, result2 error) {
	fake.containerInspectMutex.Lock()
	defer fake.containerInspectMutex.Unlock()
	fake.ContainerInspectStub = nil
	fake.containerInspectReturns = struct {
		result1 types.ContainerJSON
		result2 error
	}{result1, result2}
}

func (fake *FakeDockerClient) ContainerInspectReturnsOnCall(i int, result1 types.ContainerJSON, result2 error) {
	fake.containerInspectMutex.Lock()
	defer fake.containerInspectMutex.Unlock()
	fake.ContainerInspectStub = nil
	if fake.containerInspectReturnsOnCall == nil {
		fake.containerInspectReturnsOnCall = make(map[int]struct {
			result1 types.ContainerJSON
			result2 error
		})
	}
	fake.containerInspectReturnsOnCall[i] = struct {
		result1 types.ContainerJSON
		result2 error
	}{result1, result2}
}

func (fake *FakeDockerClient) ContainerStart(arg1 context.Context, arg2 string, arg3 types.ContainerStartOptions) error {
	fake.containerStartMutex.Lock()
	ret, specificReturn := fake.containerStartReturnsOnCall[len(fake.containerStartArgsForCall)]
	fake.containerStartArgsForCall = append(fake.containerStartArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 types.ContainerStartOptions
	}{arg1, arg2, arg3})
	stub := fake.ContainerStartStub
	fakeReturns := fake.containerStartReturns
	fake.recordInvocation("ContainerStart", []interface{}{arg1, arg2, arg3})
	fake.containerStartMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeDockerClient) ContainerStartCallCount() int {
	fake.containerStartMutex.RLock()
	defer fake.containerStartMutex.RUnlock()
	return len(fake.containerStartArgsForCall)
}

func (fake *FakeDockerClient) ContainerStartCalls(stub func(context.Context, string, types.ContainerStartOptions) error) {
	fake.containerStartMutex.Lock()
	defer fake.containerStartMutex.Unlock()
	fake.ContainerStartStub = stub
}

func (fake *FakeDockerClient) ContainerStartArgsForCall(i int) (context.Context, string, types.ContainerStartOptions) {
	fake.containerStartMutex.RLock()
	defer fake.containerStartMutex.RUnlock()
	argsForCall := fake.containerStartArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeDockerClient) ContainerStartReturns(result1 error) {
	fake.containerStartMutex.Lock()
	defer fake.containerStartMutex.Unlock()
	fake.ContainerStartStub = nil
	fake.containerStartReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDockerClient) ContainerStartReturnsOnCall(i int, result1 error) {
	fake.containerStartMutex.Lock()
	defer fake.containerStartMutex.Unlock()
	fake.ContainerStartStub = nil
	if fake.containerStartReturnsOnCall == nil {
		fake.containerStartReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.containerStartReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeDockerClient) ContainerStop(arg1 context.Context, arg2 string, arg3 *time.Duration) error {
	fake.containerStopMutex.Lock()
	ret, specificReturn := fake.containerStopReturnsOnCall[len(fake.containerStopArgsForCall)]
	fake.containerStopArgsForCall = append(fake.containerStopArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 *time.Duration
	}{arg1, arg2, arg3})
	stub := fake.ContainerStopStub
	fakeReturns := fake.containerStopReturns
	fake.recordInvocation("ContainerStop", []interface{}{arg1, arg2, arg3})
	fake.containerStopMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeDockerClient) ContainerStopCallCount() int {
	fake.containerStopMutex.RLock()
	defer fake.containerStopMutex.RUnlock()
	return len(fake.containerStopArgsForCall)
}

func (fake *FakeDockerClient) ContainerStopCalls(stub func(context.Context, string, *time.Duration) error) {
	fake.containerStopMutex.Lock()
	defer fake.containerStopMutex.Unlock()
	fake.ContainerStopStub = stub
}

func (fake *FakeDockerClient) ContainerStopArgsForCall(i int) (context.Context, string, *time.Duration) {
	fake.containerStopMutex.RLock()
	defer fake.containerStopMutex.RUnlock()
	argsForCall := fake.containerStopArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeDockerClient) ContainerStopReturns(result1 error) {
	fake.containerStopMutex.Lock()
	defer fake.containerStopMutex.Unlock()
	fake.ContainerStopStub = nil
	fake.containerStopReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDockerClient) ContainerStopReturnsOnCall(i int, result1 error) {
	fake.containerStopMutex.Lock()
	defer fake.containerStopMutex.Unlock()
	fake.ContainerStopStub = nil
	if fake.containerStopReturnsOnCall == nil {
		fake.containerStopReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.containerStopReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeDockerClient) ImageList(arg1 context.Context, arg2 types.ImageListOptions) ([]types.ImageSummary, error) {
	fake.imageListMutex.Lock()
	ret, specificReturn := fake.imageListReturnsOnCall[len(fake.imageListArgsForCall)]
	fake.imageListArgsForCall = append(fake.imageListArgsForCall, struct {
		arg1 context.Context
		arg2 types.ImageListOptions
	}{arg1, arg2})
	stub := fake.ImageListStub
	fakeReturns := fake.imageListReturns
	fake.recordInvocation("ImageList", []interface{}{arg1, arg2})
	fake.imageListMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDockerClient) ImageListCallCount() int {
	fake.imageListMutex.RLock()
	defer fake.imageListMutex.RUnlock()
	return len(fake.imageListArgsForCall)
}

func (fake *FakeDockerClient) ImageListCalls(stub func(context.Context, types.ImageListOptions) ([]types.ImageSummary, error)) {
	fake.imageListMutex.Lock()
	defer fake.imageListMutex.Unlock()
	fake.ImageListStub = stub
}

func (fake *FakeDockerClient) ImageListArgsForCall(i int) (context.Context, types.ImageListOptions) {
	fake.imageListMutex.RLock()
	defer fake.imageListMutex.RUnlock()
	argsForCall := fake.imageListArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeDockerClient) ImageListReturns(result1 []types.ImageSummary, result2 error) {
	fake.imageListMutex.Lock()
	defer fake.imageListMutex.Unlock()
	fake.ImageListStub = nil
	fake.imageListReturns = struct {
		result1 []types.ImageSummary
		result2 error
	}{result1, result2}
}

func (fake *FakeDockerClient) ImageListReturnsOnCall(i int, result1 []types.ImageSummary, result2 error) {
	fake.imageListMutex.Lock()
	defer fake.imageListMutex.Unlock()
	fake.ImageListStub = nil
	if fake.imageListReturnsOnCall == nil {
		fake.imageListReturnsOnCall = make(map[int]struct {
			result1 []types.ImageSummary
			result2 error
		})
	}
	fake.imageListReturnsOnCall[i] = struct {
		result1 []types.ImageSummary
		result2 error
	}{result1, result2}
}

func (fake *FakeDockerClient) ImagePull(arg1 context.Context, arg2 string, arg3 types.ImagePullOptions) (io.ReadCloser, error) {
	fake.imagePullMutex.Lock()
	ret, specificReturn := fake.imagePullReturnsOnCall[len(fake.imagePullArgsForCall)]
	fake.imagePullArgsForCall = append(fake.imagePullArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 types.ImagePullOptions
	}{arg1, arg2, arg3})
	stub := fake.ImagePullStub
	fakeReturns := fake.imagePullReturns
	fake.recordInvocation("ImagePull", []interface{}{arg1, arg2, arg3})
	fake.imagePullMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDockerClient) ImagePullCallCount() int {
	fake.imagePullMutex.RLock()
	defer fake.imagePullMutex.RUnlock()
	return len(fake.imagePullArgsForCall)
}

func (fake *FakeDockerClient) ImagePullCalls(stub func(context.Context, string, types.ImagePullOptions) (io.ReadCloser, error)) {
	fake.imagePullMutex.Lock()
	defer fake.imagePullMutex.Unlock()
	fake.ImagePullStub = stub
}

func (fake *FakeDockerClient) ImagePullArgsForCall(i int) (context.Context, string, types.ImagePullOptions) {
	fake.imagePullMutex.RLock()
	defer fake.imagePullMutex.RUnlock()
	argsForCall := fake.imagePullArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeDockerClient) ImagePullReturns(result1 io.ReadCloser, result2 error) {
	fake.imagePullMutex.Lock()
	defer fake.imagePullMutex.Unlock()
	fake.ImagePullStub = nil
	fake.imagePullReturns = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeDockerClient) ImagePullReturnsOnCall(i int, result1 io.ReadCloser, result2 error) {
	fake.imagePullMutex.Lock()
	defer fake.imagePullMutex.Unlock()
	fake.ImagePullStub = nil
	if fake.imagePullReturnsOnCall == nil {
		fake.imagePullReturnsOnCall = make(map[int]struct {
			result1 io.ReadCloser
			result2 error
		})
	}
	fake.imagePullReturnsOnCall[i] = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeDockerClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.containerCreateMutex.RLock()
	defer fake.containerCreateMutex.RUnlock()
	fake.containerInspectMutex.RLock()
	defer fake.containerInspectMutex.RUnlock()
	fake.containerStartMutex.RLock()
	defer fake.containerStartMutex.RUnlock()
	fake.containerStopMutex.RLock()
	defer fake.containerStopMutex.RUnlock()
	fake.imageListMutex.RLock()
	defer fake.imageListMutex.RUnlock()
	fake.imagePullMutex.RLock()
	defer fake.imagePullMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDockerClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ internal.DockerClient = new(FakeDockerClient)