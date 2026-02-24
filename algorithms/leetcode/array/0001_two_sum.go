package array

// TwoSum 1. 两数之和
//
// 题目：https://leetcode.cn/problems/two-sum/
// 难度：简单
//
// 思路：
//   1. 暴力解法：双重循环，时间复杂度 O(n²)
//   2. 哈希表：一遍遍历，用 map 存储已遍历的数字及其索引，时间复杂度 O(n)
//
// 复杂度分析：
//   - 时间复杂度：O(n)，只需要遍历一次数组
//   - 空间复杂度：O(n)，哈希表最多存储 n 个元素
//
// 参数：
//   nums: 整数数组
//   target: 目标和
//
// 返回：
//   两个数的索引，不存在则返回空切片
func TwoSum(nums []int, target int) []int {
	// 创建哈希表，存储数字到索引的映射
	numToIndex := make(map[int]int)

	// 遍历数组
	for i, num := range nums {
		// 计算需要的补数
		complement := target - num

		// 检查补数是否已在哈希表中
		if j, exists := numToIndex[complement]; exists {
			// 找到答案，返回两个索引
			return []int{j, i}
		}

		// 将当前数字及其索引存入哈希表
		numToIndex[num] = i
	}

	// 题目保证有解，这里返回 nil 保持完整性
	return nil
}

// TwoSum_BruteForce 暴力解法（用于对比）
//
// 复杂度分析：
//   - 时间复杂度：O(n²)，双重循环
//   - 空间复杂度：O(1)，只使用常量额外空间
func TwoSum_BruteForce(nums []int, target int) []int {
	n := len(nums)

	// 外层循环：遍历每个元素
	for i := 0; i < n-1; i++ {
		// 内层循环：遍历 i 后面的所有元素
		for j := i + 1; j < n; j++ {
			// 检查两数之和是否等于 target
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}

	return nil
}

// TwoSum_TwoPass 两次遍历哈希表解法
//
// 思路：
//   1. 第一次遍历：构建所有数字到索引的哈希表
//   2. 第二次遍历：检查每个数的补数是否存在
//
// 复杂度分析：
//   - 时间复杂度：O(n)，两次遍历
//   - 空间复杂度：O(n)，哈希表存储 n 个元素
//
// 注意：需要处理重复元素的情况（如 [3,3], target=6）
func TwoSum_TwoPass(nums []int, target int) []int {
	// 第一次遍历：构建哈希表
	numToIndex := make(map[int]int)
	for i, num := range nums {
		numToIndex[num] = i
	}

	// 第二次遍历：查找补数
	for i, num := range nums {
		complement := target - num
		if j, exists := numToIndex[complement]; exists && j != i {
			return []int{i, j}
		}
	}

	return nil
}
