package cache

import  pb "cache/cachepb"
// PeerPicker is a interface to locate peer that has specify key
type PeerPicker interface {
	PickPeer(key string) (peer peerGetter, ok bool)
}

// A interface must be implemented by peer
type peerGetter interface {
	Get(in *pb.Request ,out *pb.Response) error
}

