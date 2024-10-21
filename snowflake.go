package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	SnowflakeBitLen  = 63
	TimeBitLen       = 41
	TimeBitMask      = 1<<TimeBitLen - 1
	MachineIDBitLen  = 12
	MachineIDBitMask = 1<<MachineIDBitLen - 1
	SequenceBitLen   = 10
	SequenceBitMask  = 1<<SequenceBitLen - 1

	SnowflakeTimeUnit = 1e6 // msec in nanosec
)

type Settings struct {
	StartTime time.Time
	MachineID string
}

type Snowflake struct {
	mutex       *sync.Mutex
	startTime   int64
	elapsedTime int64
	sequence    uint16
	machineID   uint16
}

var (
	ErrStartTimeAhead   = errors.New("start time is ahead of now")
	ErrZeroStartTime    = errors.New("no start time is provided")
	ErrInvalidMachineID = errors.New("invalid machine id")
	ErrOverTimeLimit    = errors.New("timestamp is over the capacity of Snowflake")
)

func New(st Settings) (*Snowflake, error) {
	err := validateSettings(st)
	if err != nil {
		return nil, err
	}

	sf := &Snowflake{
		mutex:       new(sync.Mutex),
		startTime:   toSnowflakeTime(st.StartTime),
		elapsedTime: 0,
		sequence:    0,
		machineID:   hash(st.MachineID, MachineIDBitMask),
	}

	return sf, nil
}

// NextID generates a next unique ID.
// After the Snowflake time overflows, NextID returns an error.
func (sf *Snowflake) NextID() (int64, error) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	// diff timeNow - timeStart
	current := currentElapsedTime(sf.startTime)
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else { // sf.elapsedTime >= current
		sf.sequence = (sf.sequence + 1) & SequenceBitMask
		if sf.sequence == 0 { // if overflow
			sf.elapsedTime++
			overtime := sf.elapsedTime - current
			time.Sleep(sleepTime((overtime)))
		}
	}

	return sf.toID()
}

func validateSettings(st Settings) error {
	if st.StartTime.IsZero() {
		return ErrZeroStartTime
	}
	if st.StartTime.After(time.Now()) {
		return ErrStartTimeAhead
	}

	if st.MachineID == "" {
		return ErrInvalidMachineID
	}

	return nil
}

func (sf *Snowflake) toID() (int64, error) {
	if sf.elapsedTime >= 1<<TimeBitLen {
		return 0, ErrOverTimeLimit
	}

	return sf.elapsedTime<<(MachineIDBitLen+SequenceBitLen) |
		int64(sf.machineID)<<SequenceBitLen |
		int64(sf.sequence), nil
}

func toSnowflakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / SnowflakeTimeUnit
}

func currentElapsedTime(startTime int64) int64 {
	return toSnowflakeTime(time.Now()) - startTime
}

func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime*SnowflakeTimeUnit) -
		time.Duration(time.Now().UTC().UnixNano()%SnowflakeTimeUnit)
}
