---
title: tm_cidrhost - Functions - Configuration Language
description: |-
  The tm_cidrhost function calculates a full host IP address within a given
  IP network address prefix.
---

# `tm_cidrhost` Function

`tm_cidrhost` calculates a full host IP address for a given host number within
a given IP network address prefix.

```hcl
tm_cidrhost(prefix, hostnum)
```

`prefix` must be given in CIDR notation, as defined in
[RFC 4632 section 3.1](https://tools.ietf.org/html/rfc4632#section-3.1).

`hostnum` is a whole number that can be represented as a binary integer with
no more than the number of digits remaining in the address after the given
prefix. For more details on how this function interprets CIDR prefixes and
populates host numbers, see the worked example for
[`tm_cidrsubnet`](./tm_cidrsubnet.md).

Conventionally host number zero is used to represent the address of the
network itself and the host number that would fill all the host bits with
binary 1 represents the network's broadcast address. These numbers should
generally not be used to identify individual hosts except in unusual
situations, such as point-to-point links.

This function accepts both IPv6 and IPv4 prefixes, and the result always uses
the same addressing scheme as the given prefix.

-> **Note:** As a historical accident, this function interprets IPv4 address
octets that have leading zeros as decimal numbers, which is contrary to some
other systems which interpret them as octal. We have preserved this behavior
for backward compatibility, but recommend against relying on this behavior.

## Examples

```
tm_cidrhost("10.12.112.0/20", 16)
10.12.112.16
tm_cidrhost("10.12.112.0/20", 268)
10.12.113.12
tm_cidrhost("fd00:fd12:3456:7890:00a2::/72", 34)
fd00:fd12:3456:7890::22
```

## Related Functions

* [`tm_cidrsubnet`](./tm_cidrsubnet.md) calculates a subnet address under a given
  network address prefix.
