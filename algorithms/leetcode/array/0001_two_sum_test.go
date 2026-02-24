package array

import (
	"reflect"
	"testing"
)

// TestTwoSum 测试哈希表解法
func TestTwoSum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		target   int
		want     []int
		wantErr  bool
	}{
		{
			name:   "示例1-正常情况",
			nums:   []int{2, 7, 11, 15},
			target: 9,
			want:   []int{0, 1},
		},
		{
			name:   "示例2-负数",
			nums:   []int{3, 2, 4},
			target: 6,
			want:   []int{1, 2},
		},
		{
			name:   "示例3-重复元素",
			nums:   []int{3, 3},
			target: 6,
			want:   []int{0, 1},
		},
		{
			name:   "边界情况-两个元素",
			nums:   []int{1, 2},
			target: 3,
			want:   []int{0, 1},
		},
		{
			name:   "包含负数",
			nums:   []int{-1, -2, -3, -4, -5},
			target: -8,
			want:   []int{2, 4},
		},
		{
			name:   "混合正负数",
			nums:   []int{0, 4, 3, 0},
			target: 0,
			want:   []int{0, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TwoSum(tt.nums, tt.target)

			// 验证结果
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TwoSum() = %v, want %v", got, tt.want)
				return
			}

			// 验证结果的正确性：两数之和是否等于 target
			if len(got) == 2 {
				sum := tt.nums[got[0]] + tt.nums[got[1]]
				if sum != tt.target {
					t.Errorf("结果验证失败: nums[%d]=%d + nums[%d]=%d = %d, want %d",
						got[0], tt.nums[got[0]], got[1], tt.nums[got[1]], sum, tt.target)
				}
			}
		})
	}
}

// TestTwoSum_BruteForce 测试暴力解法
func TestTwoSum_BruteForce(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{
			name:   "示例1",
			nums:   []int{2, 7, 11, 15},
			target: 9,
			want:   []int{0, 1},
		},
		{
			name:   "示例2",
			nums:   []int{3, 2, 4},
			target: 6,
			want:   []int{1, 2},
		},
		{
			name:   "示例3",
			nums:   []int{3, 3},
			target: 6,
			want:   []int{0, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TwoSum_BruteForce(tt.nums, tt.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TwoSum_BruteForce() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTwoSum_TwoPass 测试两次遍历解法
func TestTwoSum_TwoPass(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{
			name:   "示例1",
			nums:   []int{2, 7, 11, 15},
			target: 9,
			want:   []int{0, 1},
		},
		{
			name:   "示例2",
			nums:   []int{3, 2, 4},
			target: 6,
			want:   []int{1, 2},
		},
		{
			name:   "示例3",
			nums:   []int{3, 3},
			target: 6,
			want:   []int{0, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TwoSum_TwoPass(tt.nums, tt.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TwoSum_TwoPass() = %v, want %v", got, tt.want)
			}
		})
	}
}

// BenchmarkTwoSum 性能测试 - 哈希表解法
func BenchmarkTwoSum(b *testing.B) {
	nums := []int{2, 7, 11, 15, 3, 6, 8, 10, 1, 5}
	target := 9

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TwoSum(nums, target)
	}
}

// BenchmarkTwoSum_BruteForce 性能测试 - 暴力解法
func BenchmarkTwoSum_BruteForce(b *testing.B) {
	nums := []int{2, 7, 11, 15, 3, 6, 8, 10, 1, 5}
	target := 9

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TwoSum_BruteForce(nums, target)
	}
}

// BenchmarkTwoSum_TwoPass 性能测试 - 两次遍历解法
func BenchmarkTwoSum_TwoPass(b *testing.B) {
	nums := []int{2, 7, 11, 15, 3, 6, 8, 10, 1, 5}
	target := 9

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TwoSum_TwoPass(nums, target)
	}
}
