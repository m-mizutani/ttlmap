package ttlmap_test

import (
	"testing"

	"github.com/m-mizutani/ttlmap"
	"github.com/stretchr/testify/assert"
)

type testData string

func TestBasicUsage(t *testing.T) {
	ttlMap := ttlmap.New(10)
	v1 := testData("hoge")
	k1 := []byte("k1")

	assert.Nil(t, ttlMap.Set(k1, &v1, 5))
	p2 := ttlMap.Get(k1)
	assert.NotNil(t, p2)
	v2, ok := p2.(*testData)
	assert.True(t, ok)
	assert.Equal(t, "hoge", string(*v2))

	assert.Equal(t, 0, len(ttlMap.Prune(5)))
	assert.NotNil(t, ttlMap.Get(k1))
	assert.Equal(t, 1, len(ttlMap.Prune(1)))
	assert.Nil(t, ttlMap.Get(k1))
}
