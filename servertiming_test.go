package servertiming

import (
	"testing"

	"github.com/cheekybits/is"
)

func TestString(t *testing.T) {
	is := is.New(t)

	ti := New()
	is.NotNil(ti)
}
