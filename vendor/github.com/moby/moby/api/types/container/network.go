package container

// PortRangeProto is a string containing port number and protocol in the format "80/tcp".
type PortProto string

// PortSet is a collection of structs indexed by [PortProto].
type PortSet = map[PortProto]struct{}

// PortBinding represents a binding between a Host IP address and a Host Port.
type PortBinding struct {
	// HostIP is the host IP Address
	HostIP string `json:"HostIp"`
	// HostPort is the host port number
	HostPort string
}

// PortMap is a collection of [PortBinding] indexed by [PortProto].
type PortMap = map[PortProto][]PortBinding
