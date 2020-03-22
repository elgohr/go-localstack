package localstack

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/elgohr/go-localstack/internal/internalfakes"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInstance_Start_Fails(t *testing.T) {
	for _, tt := range [...]struct {
		name  string
		given func() *Instance
		then  func(err error)
	}{
		{
			name: "can't restart localstack when already running",
			given: func() *Instance {
				fakePool := &internalfakes.FakePool{}
				fakePool.PurgeReturns(errors.New("can't start"))
				return &Instance{
					pool:     fakePool,
					resource: &dockertest.Resource{},
				}
			},
			then: func(err error) {
				assert.Equal(t, "localstack: can't stop an already running instance: can't start", err.Error())
			},
		},
		{
			name: "can't start container",
			given: func() *Instance {
				fakePool := &internalfakes.FakePool{}
				fakePool.RunWithOptionsReturns(nil, errors.New("can't start container"))
				return &Instance{
					pool: fakePool,
				}
			},
			then: func(err error) {
				assert.Equal(t, "localstack: could not start container: can't start container", err.Error())
			},
		},
		{
			name: "fails during waiting on startup",
			given: func() *Instance {
				fakePool := &internalfakes.FakePool{}
				fakePool.RetryReturns(errors.New("can't wait"))
				return &Instance{
					pool: fakePool,
				}
			},
			then: func(err error) {
				assert.Equal(t, "localstack: could not start environment: can't wait", err.Error())
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt.then(tt.given().Start())
		})
	}
}

func TestInstance_Stop_Fails(t *testing.T) {
	fakePool := &internalfakes.FakePool{}
	fakePool.PurgeReturns(errors.New("can't stop"))
	i := &Instance{
		pool:     fakePool,
		resource: &dockertest.Resource{},
	}

	assert.Equal(t, "can't stop", i.Stop().Error())
}

func TestInstance_isAvailable_Session_Fails(t *testing.T) {
	empty := &aws.Config{}
	if err := os.Setenv("AWS_STS_REGIONAL_ENDPOINTS", "FAILURE"); err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv("AWS_STS_REGIONAL_ENDPOINTS")
	assert.Error(t, isAvailable(empty))
}
