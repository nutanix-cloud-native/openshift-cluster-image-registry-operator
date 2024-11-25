// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// ImageRegistryConfigProxyApplyConfiguration represents a declarative configuration of the ImageRegistryConfigProxy type for use
// with apply.
type ImageRegistryConfigProxyApplyConfiguration struct {
	HTTP    *string `json:"http,omitempty"`
	HTTPS   *string `json:"https,omitempty"`
	NoProxy *string `json:"noProxy,omitempty"`
}

// ImageRegistryConfigProxyApplyConfiguration constructs a declarative configuration of the ImageRegistryConfigProxy type for use with
// apply.
func ImageRegistryConfigProxy() *ImageRegistryConfigProxyApplyConfiguration {
	return &ImageRegistryConfigProxyApplyConfiguration{}
}

// WithHTTP sets the HTTP field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the HTTP field is set to the value of the last call.
func (b *ImageRegistryConfigProxyApplyConfiguration) WithHTTP(value string) *ImageRegistryConfigProxyApplyConfiguration {
	b.HTTP = &value
	return b
}

// WithHTTPS sets the HTTPS field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the HTTPS field is set to the value of the last call.
func (b *ImageRegistryConfigProxyApplyConfiguration) WithHTTPS(value string) *ImageRegistryConfigProxyApplyConfiguration {
	b.HTTPS = &value
	return b
}

// WithNoProxy sets the NoProxy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NoProxy field is set to the value of the last call.
func (b *ImageRegistryConfigProxyApplyConfiguration) WithNoProxy(value string) *ImageRegistryConfigProxyApplyConfiguration {
	b.NoProxy = &value
	return b
}
