package snowflake

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Snowflake(t *testing.T) {
	t.Parallel()

	podsNum := 10
	prefix := "origin-vod-"
	uniquePodNames := make(map[string]struct{})
	for len(uniquePodNames) < podsNum {
		str := prefix + generateRandomString(10) + "-" + generateRandomString(5)
		uniquePodNames[str] = struct{}{}
	}

	instances := make([]*Snowflake, 0, podsNum)
	for name := range uniquePodNames {
		snowflake, err := New(Settings{
			StartTime: time.Now(),
			MachineID: name,
		})
		assert.NoError(t, err)

		instances = append(instances, snowflake)
	}

	consumer := make(chan int64)
	IDnum := 10000
	generate := func(sf *Snowflake) {
		for i := 0; i < IDnum; i++ {
			id, err := sf.NextID()
			assert.NoError(t, err)
			consumer <- id
		}
	}

	for _, i := range instances {
		go generate(i)
	}

	set := make(map[int64]struct{})
	for i := 0; i < IDnum*podsNum; i++ {
		id := <-consumer
		assert.GreaterOrEqual(t, id, int64(0))
		set[id] = struct{}{}
	}

	// collisions num
	assert.Equal(t, 0, IDnum*podsNum-len(set))
}
