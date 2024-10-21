package snowflake

import (
	"errors"
	"time"
)

var ErrNegativeID = errors.New("id is negative")

// ElapsedTime returns the elapsed time when the given Snowflake ID was generated.
func ElapsedTime(id int64) time.Duration {
	return time.Duration(elapsedTime(id) * SnowflakeTimeUnit)
}

func elapsedTime(id int64) int64 {
	return id >> (MachineIDBitLen + SequenceBitLen)
}

// SequenceNumber returns the sequence number of a Snowflake ID.
func SequenceNumber(id int64) int64 {
	return id & SequenceBitMask
}

// MachineID returns the machine ID of a Snowflake ID.
func MachineID(id int64) int64 {
	const mask = int64(MachineIDBitMask << SequenceBitLen)
	return id & mask >> SequenceBitLen
}

// Decompose returns a set of Snowflake ID parts.
func Decompose(id int64) (map[string]int64, error) {
	if id < 0 {
		return nil, ErrNegativeID
	}
	t := elapsedTime(id)
	sequence := SequenceNumber(id)
	machineID := MachineID(id)
	return map[string]int64{
		"id":         id,
		"time":       t,
		"sequence":   sequence,
		"machine-id": machineID,
	}, nil
}
