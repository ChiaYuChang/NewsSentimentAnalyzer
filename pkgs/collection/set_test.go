package collection_test

import (
	"testing"

	"github.com/ChiaYuChang/NewsSentimentAnalyzer/pkgs/collection"
	"github.com/stretchr/testify/require"
)

func TestSetAddandDel(t *testing.T) {
	s := collection.NewSet(1)
	require.True(t, s.Has(1))

	s.Add(1)
	require.True(t, s.Has(1))

	s.Del(1)
	require.False(t, s.Has(1))

	s.Add(1)
	require.True(t, s.Has(1))

	s.Del(1)
	require.False(t, s.Has(1))
}

func TestSetMerge(t *testing.T) {
	s1 := collection.NewSet(1, 2, 3)
	s2 := collection.NewSet(1, 2, 4)
	s3 := s1.Merge(s2)

	require.True(t, s1.Has(1))
	require.True(t, s1.Has(2))
	require.True(t, s1.Has(3))
	require.False(t, s1.Has(4))

	require.True(t, s2.Has(1))
	require.True(t, s2.Has(2))
	require.False(t, s2.Has(3))
	require.True(t, s2.Has(4))

	require.True(t, s3.Has(1))
	require.True(t, s3.Has(2))
	require.True(t, s3.Has(3))
	require.True(t, s3.Has(4))

	s3.Add(10)
	require.False(t, s1.Has(10))
	require.False(t, s2.Has(10))
	require.True(t, s3.Has(10))
	require.ElementsMatch(t, []int{1, 2, 3, 4, 10}, s3.Key())
}
