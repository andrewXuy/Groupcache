# groupcache
## Summary

This pacakage is a re-implementation of group cache.

Package groupcache provides a data loading mechanism with caching and de-duplication that works across a set of peer processes.

Each data Get first consults its local cache, otherwise delegates to the requested key's canonical owner, which then checks its cache or finally gets the data. In the common case, many concurrent cache misses across a set of peers for the same key result in just one cache fill.


For original version and examples, see https://godoc.org/github.com/golang/groupcache

## Version
 **go 1.14 darwin/amd64**
 
 **github.com/golang/protobuf v1.3.5**
 
## Demo

There is a simple demo script named run.sh. You can easily run it by typing ./run.sh

## Usage

### Loading Process

a groupcache lookup of **Get("foo")** like:
(On node #5 of a set of N node running the same code)

 1. Is the value of "foo" in local memory because it's super hot?  If so, use it.

 2. Is the value of "foo" in local memory because peer #5 (the current
    peer) is the owner of it?  If so, use it.

 3. Amongst all the peers in my set of N, am I the owner of the key
    "foo"?  (e.g. does it consistent hash to 5?)  If so, load it.  If
    other callers come in, via the same process or via RPC requests
    from peers, they block waiting for the load to finish and get the
    same answer.  If not, RPC to the peer that's the owner and get
    the answer.  If the RPC fails, just load it locally (still with
    local dup suppression).

## Implementation

### Consisitent Hash

Implement Consisitent Hash with virtual address. Solving imbalance hash distribution.

### Single Filght

Implement a functionality that response once with large scale requests. Prevent cache breakdown.

### Protocol Buffers 

Use Protobuf as the Intermediate between HTTP connection.

### Distributed LRU cache

Use the functionality above create the safe concurrent and database friendly distributed cache
