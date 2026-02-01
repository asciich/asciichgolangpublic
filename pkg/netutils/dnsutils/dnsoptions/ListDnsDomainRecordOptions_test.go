package dnsoptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/dnsutils/dnsoptions"
)

func Test_DnsDomainRecordOptions_GetRecordNameOrEmptyString(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		options := &dnsoptions.ListDnsDomainRecordOptions{}
		recordName := options.GetRecordNameOrEmptyStringIfUnset()
		require.EqualValues(t, "", recordName)
	})

	t.Run("only record set", func(t *testing.T) {
		options := &dnsoptions.ListDnsDomainRecordOptions{
			Name: "example",
		}
		recordName := options.GetRecordNameOrEmptyStringIfUnset()
		require.EqualValues(t, "example", recordName)
	})
}
