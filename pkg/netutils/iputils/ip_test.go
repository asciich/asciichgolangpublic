package iputils_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/asciich/asciichgolangpublic/pkg/contextutils"
    "github.com/asciich/asciichgolangpublic/pkg/netutils/iputils"
    "github.com/asciich/asciichgolangpublic/pkg/testutils"
    "github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func getCtx() context.Context {
    return contextutils.ContextVerbose()
}

func TestIsValidIP(t *testing.T) {
    tests := []struct {
        ip           string
        expectedValid bool
    }{
        {"192.168.1.1", true},
        {"::1", true},
        {"256.1.1.1", false},
        {"invalid", false},
        {"", false},
        {" 192.168.1.1", false},
        {"192.168.1.1 ", false},
        {" 192.168.1.1 ", false},
        {" ::1", false},
        {"::1 ", false},
    }

    for _, tt := range tests {
        t.Run(
            testutils.MustFormatAsTestname(tt),
            func(t *testing.T) {
                isValid, err := iputils.IsValidIP(getCtx(), tt.ip)
                if tt.ip == "" {
                    require.Error(t, err)
                    require.False(t, isValid)
                } else {
                    require.NoError(t, err)
                    require.Equal(t, tt.expectedValid, isValid)
                }
            },
        )
    }
}

func TestCheckValidIP(t *testing.T) {
    tests := []struct {
        ip          string
        expectError bool
        errType     error
    }{
        {"192.168.1.1", false, nil},
        {"::1", false, nil},
        {"256.1.1.1", true, iputils.ErrInvalidIP},
        {"invalid", true, iputils.ErrInvalidIP},
        {"", true, tracederrors.ErrTracedErrorEmptyString},
        {" 192.168.1.1", true, iputils.ErrInvalidIP},
        {"192.168.1.1 ", true, iputils.ErrInvalidIP},
        {" ::1", true, iputils.ErrInvalidIP},
    }

    for _, tt := range tests {
        t.Run(
            testutils.MustFormatAsTestname(tt),
            func(t *testing.T) {
                err := iputils.CheckValidIP(getCtx(), tt.ip)
                if tt.expectError {
                    require.Error(t, err)
                    if tt.errType != nil {
                        require.ErrorIs(t, err, tt.errType)
                    }
                } else {
                    require.NoError(t, err)
                }
            },
        )
    }
}

func TestIsValidIPv4(t *testing.T) {
    tests := []struct {
        ip           string
        expectedValid bool
    }{
        {"192.168.1.1", true},
        {"::1", false},
        {"256.1.1.1", false},
        {"", false},
        {" 192.168.1.1", false},
        {"192.168.1.1 ", false},
        {" 192.168.1.1 ", false},
    }

    for _, tt := range tests {
        t.Run(
            testutils.MustFormatAsTestname(tt),
            func(t *testing.T) {
                isValid, err := iputils.IsValidIPv4(getCtx(), tt.ip)
                if tt.ip == "" {
                    require.Error(t, err)
                    require.False(t, isValid)
                } else {
                    require.NoError(t, err)
                    require.Equal(t, tt.expectedValid, isValid)
                }
            },
        )
    }
}

func TestCheckValidIPv4(t *testing.T) {
    tests := []struct {
        ip          string
        expectError bool
        errType     error
    }{
        {"192.168.1.1", false, nil},
        {"::1", true, iputils.ErrInvalidIPv4},
        {"256.1.1.1", true, iputils.ErrInvalidIPv4},
        {"", true, tracederrors.ErrTracedErrorEmptyString},
        {" 192.168.1.1", true, iputils.ErrInvalidIPv4},
        {"192.168.1.1 ", true, iputils.ErrInvalidIPv4},
    }

    for _, tt := range tests {
        t.Run(
            testutils.MustFormatAsTestname(tt),
            func(t *testing.T) {
                err := iputils.CheckValidIPv4(getCtx(), tt.ip)
                if tt.expectError {
                    require.Error(t, err)
                    if tt.errType != nil {
                        require.ErrorIs(t, err, tt.errType)
                    }
                } else {
                    require.NoError(t, err)
                }
            },
        )
    }
}

func TestIsValidIPv6(t *testing.T) {
    tests := []struct {
        ip           string
        expectedValid bool
    }{
        {"::1", true},
        {"192.168.1.1", false},
        {"invalid", false},
        {"", false},
        {" ::1", false},
        {"::1 ", false},
        {" ::1 ", false},
    }

    for _, tt := range tests {
        t.Run(
            testutils.MustFormatAsTestname(tt),
            func(t *testing.T) {
                isValid, err := iputils.IsValidIPv6(getCtx(), tt.ip)
                if tt.ip == "" {
                    require.Error(t, err)
                    require.False(t, isValid)
                } else {
                    require.NoError(t, err)
                    require.Equal(t, tt.expectedValid, isValid)
                }
            },
        )
    }
}

func TestCheckValidIPv6(t *testing.T) {
    tests := []struct {
        ip          string
        expectError bool
        errType     error
    }{
        {"::1", false, nil},
        {"2001:db8::1", false, nil},
        {"192.168.1.1", true, iputils.ErrInvalidIPv6},
        {"invalid", true, iputils.ErrInvalidIPv6},
        {"", true, tracederrors.ErrTracedErrorEmptyString},
        {" ::1", true, iputils.ErrInvalidIPv6},
        {"::1 ", true, iputils.ErrInvalidIPv6},
    }

    for _, tt := range tests {
        t.Run(
            testutils.MustFormatAsTestname(tt),
            func(t *testing.T) {
                err := iputils.CheckValidIPv6(getCtx(), tt.ip)
                if tt.expectError {
                    require.Error(t, err)
                    if tt.errType != nil {
                        require.ErrorIs(t, err, tt.errType)
                    }
                } else {
                    require.NoError(t, err)
                }
            },
        )
    }
}

func TestErrorHelpers(t *testing.T) {
    t.Run("IsInvalidIPError", func(t *testing.T) {
        require.True(t, iputils.IsInvalidIPError(iputils.ErrInvalidIP))
        require.True(t, iputils.IsInvalidIPError(iputils.ErrInvalidIPv4))
        require.True(t, iputils.IsInvalidIPError(iputils.ErrInvalidIPv6))
        require.False(t, iputils.IsInvalidIPError(tracederrors.ErrTracedErrorEmptyString))
    })

    t.Run("IsInvalidIPv4Error", func(t *testing.T) {
        require.True(t, iputils.IsInvalidIPv4Error(iputils.ErrInvalidIPv4))
        require.False(t, iputils.IsInvalidIPv4Error(iputils.ErrInvalidIP))
        require.False(t, iputils.IsInvalidIPv4Error(iputils.ErrInvalidIPv6))
        require.False(t, iputils.IsInvalidIPv4Error(tracederrors.ErrTracedErrorEmptyString))
    })

    t.Run("IsInvalidIPv6Error", func(t *testing.T) {
        require.True(t, iputils.IsInvalidIPv6Error(iputils.ErrInvalidIPv6))
        require.False(t, iputils.IsInvalidIPv6Error(iputils.ErrInvalidIP))
        require.False(t, iputils.IsInvalidIPv6Error(iputils.ErrInvalidIPv4))
        require.False(t, iputils.IsInvalidIPv6Error(tracederrors.ErrTracedErrorEmptyString))
    })
}
