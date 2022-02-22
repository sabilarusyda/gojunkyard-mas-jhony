package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromEnvFile(t *testing.T) {
	os.Clearenv()

	// Since we didn't load env file, http port must be empty
	httpPort := os.Getenv("devcode.xeemore.com/systech/gojunkyard_HTTP_PORT")
	assert.Empty(t, httpPort)

	// Load from file
	Load("fixtures/default.env")

	httpPort = os.Getenv("devcode.xeemore.com/systech/gojunkyard_HTTP_PORT")
	assert.NotEmpty(t, httpPort)
	assert.Equal(t, "8080", httpPort)
}

func TestLoadFromEnvVar(t *testing.T) {
	os.Clearenv()

	// Let's say we have an environment variable
	os.Setenv("KEY_1", "1")

	// Load from file
	Load("fixtures/default.env")

	key1 := os.Getenv("KEY_1")
	assert.NotEmpty(t, key1)
	assert.Equal(t, "1", key1)
}

func TestLoadFromEnvFileNoOverride(t *testing.T) {
	os.Clearenv()

	// Let's say we have an environment variable
	os.Setenv("devcode.xeemore.com/systech/gojunkyard_HTTP_PORT", "9090")

	// Load from file that also have environment variable
	Load("fixtures/default.env")

	// The value must not overriden by Load() function
	httpPort := os.Getenv("devcode.xeemore.com/systech/gojunkyard_HTTP_PORT")
	assert.NotEmpty(t, httpPort)
	assert.Equal(t, "9090", httpPort)
}

type Specification struct {
	HTTPPort string `envconfig:"HTTP_PORT"`
}

func TestParse(t *testing.T) {
	os.Clearenv()

	// Let's say we have an environment variable
	os.Setenv("devcode.xeemore.com/systech/gojunkyard_HTTP_PORT", "9090")

	var spec Specification

	Parse("devcode.xeemore.com/systech/gojunkyard", &spec)

	assert.Equal(t, "9090", spec.HTTPPort)
}

func TestParseNoPrefix(t *testing.T) {
	os.Clearenv()

	// Let's say we have an environment variable
	os.Setenv("HTTP_PORT", "9090")

	var spec Specification

	Parse("", &spec)

	assert.Equal(t, "9090", spec.HTTPPort)
}

type NestedSpecification struct {
	Server ServerSpecification `envconfig:"SERVER"`
}

type ServerSpecification struct {
	HTTPPort string `envconfig:"HTTP_PORT"`
}

func TestParseNested(t *testing.T) {
	os.Clearenv()

	// Let's say we have environment variables
	os.Setenv("devcode.xeemore.com/systech/gojunkyard_SERVER_HTTP_PORT", "9090")
	os.Setenv("devcode.xeemore.com/systech/gojunkyard_DATABASE_MASTER_DRIVER", "mysql")

	var nestedSpec NestedSpecification

	Parse("devcode.xeemore.com/systech/gojunkyard", &nestedSpec)

	assert.Equal(t, "9090", nestedSpec.Server.HTTPPort)
}

func TestLoadAndParse(t *testing.T) {
	os.Clearenv()

	var spec Specification

	err := LoadAndParse("devcode.xeemore.com/systech/gojunkyard", &spec, "fixtures/default.env")
	assert.Nil(t, err)
	assert.Equal(t, "8080", spec.HTTPPort)
}
