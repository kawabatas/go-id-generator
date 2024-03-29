# go-id-generator
This package generates IDs based on the [Twitter Snowflake](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake) ID format.

```
// # A Twitter Snowflake ID is composed of
//
//	 | 0 | 00000000000000000000000000000000000000000 | 00000 | 00000 | 000000000000 |
//  1.unused       2.timestamp (millisecond)       3.datacenterID     5.sequenceNumber
//                                                          4.machineID
```

- The first bit is unused. In total, a generated ID is 63 bits, so this package returns an `int64` value, not a `uint64` value.
- The next 41 bits are used for a timestamp (in `milliseconds`), representing the offset of the current time relative to a certain time (by default this package uses `2024-01-01 00:00:00 UTC`, which differs from the Twitter Snowflake Epoch). This allows for a maximum lifetime of 2^41 - 1 milliseconds, which is approximately 69 years.
- The following 5 bits are used for a datacenter ID, with a maximum value of 2^5 - 1 = 31.
- The next 5 bits are used for a machine ID, with a maximum value of 2^5 - 1 = 31.
- The last 12 bits are used for a sequence number, with a maximum value of 2^12 - 1 = 4095.

# Installation
```
go get github.com/kawabatas/go-id-generator
```

# Usage
```
TODO
```
