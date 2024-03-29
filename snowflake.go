// Package idgenerator provides a simple and easy-to-use Twitter Snowflake generator.
//
// # A Twitter Snowflake ID is composed of
//
//	 | 0 | 00000000000000000000000000000000000000000 | 00000 | 00000 | 000000000000 |
//	1.unused       2.timestamp (millisecond)       3.datacenterID     5.sequenceNumber
//	                                                         4.machineID
//	1st	 1 bit  is  unused
//	2nd 41 bits are used for timestamp (millisecond)
//	3rd	 5 bits are used for a datacenter id
//	4th	 5 bits are used for a machine id
//	5th 12 bits are used for a sequence number
package idgenerator

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	timestampBitRange   = 41
	datacenterBitRange  = 5
	machineBitRange     = 5
	sequenceNumBitRange = 12
)

var (
	timestampBitShift  = datacenterBitRange + machineBitRange + sequenceNumBitRange
	datacenterBitShift = machineBitRange + sequenceNumBitRange
	machineBitShift    = sequenceNumBitRange

	defaultBaseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

var (
	ErrOverLifeTime          = errors.New("over the maximum lifetime")
	ErrInvalidTimestamp      = errors.New("invalid timestamp")
	ErrInvalidDatacenterID   = errors.New("invalid datacenter ID")
	ErrInvalidMachineID      = errors.New("invalid machine ID")
	ErrInvalidSequenceNumber = errors.New("invalid sequence number")
)

type snowflake struct {
	timestamp      int64
	datacenterID   int
	machineID      int
	sequenceNumber int

	baseTime time.Time
	random   bool

	mutex sync.Mutex
}

type option func(*snowflake) error

// NewSnowflakeID returns a new generated Snowflake ID.
func NewSnowflakeID(opts ...option) (int64, error) {
	s := &snowflake{}

	s.mutex.Lock()
	for _, f := range opts {
		if err := f(s); err != nil {
			return 0, err
		}
	}

	ts, err := s.getElapsedTimestamp()
	if err != nil {
		return 0, err
	}
	s.timestamp = ts

	if s.random {
		if s.datacenterID == 0 {
			s.datacenterID = rand.Intn(2 ^ datacenterBitRange - 1)
		}
		if s.machineID == 0 {
			s.machineID = rand.Intn(2 ^ machineBitRange - 1)
		}
		if s.sequenceNumber == 0 {
			s.sequenceNumber = rand.Intn(2 ^ sequenceNumBitRange - 1)
		}
	}
	s.mutex.Unlock()

	generatedID := s.timestamp<<timestampBitShift | int64(s.datacenterID)<<datacenterBitShift | int64(s.machineID)<<machineBitShift | int64(s.sequenceNumber)
	return generatedID, nil
}

// WithTimestamp specifies the timestamp of Snowflake ID.
func WithTimestamp(v time.Time) option {
	return func(s *snowflake) error {
		s.timestamp = v.UnixMilli()
		return nil
	}
}

// WithDatacenterID specifies the datacenter ID of Snowflake ID.
func WithDatacenterID(v int) option {
	return func(s *snowflake) error {
		if v < 0 || v > int(math.Pow(2, datacenterBitRange))-1 {
			return ErrInvalidDatacenterID
		}
		s.datacenterID = v
		return nil
	}
}

// WithMachineID specifies the machine ID of Snowflake ID.
func WithMachineID(v int) option {
	return func(s *snowflake) error {
		if v < 0 || v > int(math.Pow(2, machineBitRange))-1 {
			return ErrInvalidMachineID
		}
		s.machineID = v
		return nil
	}
}

// WithSequenceNumber specifies the sequence number of Snowflake ID.
func WithSequenceNumber(v int) option {
	return func(s *snowflake) error {
		if v < 0 || v > int(math.Pow(2, sequenceNumBitRange))-1 {
			return ErrInvalidSequenceNumber
		}
		s.sequenceNumber = v
		return nil
	}
}

// WithBaseTime changes the Snowflake base time from the default.
func WithBaseTime(v time.Time) option {
	return func(s *snowflake) error {
		s.baseTime = v
		return nil
	}
}

// WithRandomEnabled enables picking a random value for unset datacenter ID, machine ID, and sequence number.
func WithRandomEnabled() option {
	return func(s *snowflake) error {
		s.random = true
		return nil
	}
}

func (s *snowflake) getElapsedTimestamp() (int64, error) {
	at := time.Now().UTC()
	if s.timestamp > 0 {
		at = time.UnixMilli(s.timestamp)
	}

	baseTime := defaultBaseTime
	if !s.baseTime.IsZero() {
		baseTime = s.baseTime
	}

	diffMilli := at.Sub(baseTime).Milliseconds()
	if diffMilli <= 0 {
		return 0, ErrInvalidTimestamp
	} else if diffMilli > int64(math.Pow(2, timestampBitRange))-1 {
		return 0, ErrOverLifeTime
	}
	return diffMilli, nil
}
