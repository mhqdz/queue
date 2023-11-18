package queue

import "sync"

// 循环队列 用来避开slice的循环问题
type queue[T any] struct {
	mutex sync.Mutex

	// 指向现有的头
	index int

	// 队列总长度
	len int

	// 数据存储的地方
	data []T
}

// New 传入队列最大长度 并返回一个空队列
func New[T any](len int) *queue[T] {
	if len < 0 {
		panic("queue length must be greater than 0")
	}
	return &queue[T]{len: len, index: 0, data: make([]T, 0)}
}

// NewBySlice 从一个数组创建一个queue
func NewBySlice[T any](slice []T) *queue[T] {
	return &queue[T]{len: len(slice), index: 0, data: slice}
}

// Get 获取第i个数据 从0开始 相当于slice Data的 data[i]
func (q *queue[T]) Get(i int) T {
	return q.data[(i+q.index)%q.len]
}

// Set 设置第i个数据的值为v 相当于slice Data的 data[i] = v
func (q *queue[T]) Set(i int, v T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.data[(i+q.index)%q.len] = v
}

// Len 获取queue最大长度 相当于len(data)
func (q *queue[T]) Len() int {
	return q.len
}

// Append 添加一个元素
func (q *queue[T]) Append(v T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.data) < q.len {
		q.data = append(q.data, v)
		return
	}
	q.data[q.index] = v
	q.index = (q.index + 1) % q.len
	// q.Cnt++
}

// GetData 获取Data 注意这个Data是乱序的
func (q *queue[T]) GetData() []T {
	return q.data
}

// GetSlice 返回一个index从0开始的数组
func (q *queue[T]) GetSlice() []T {
	newAyy := make([]T, len(q.data))
	copy(newAyy, append(q.data[q.index:], q.data[:q.index]...))
	return newAyy
}

/*
如果新的len小于旧的 则会舍弃一部分旧数据
如果是扩容则相当于在index前面(如果index是0则是最后面) 插入一段 然后把index后面的拷贝到最后去
插入那段不会置空 也不会在被覆盖之前被使用 但有可能被get取到
*/
// ReSetLen 重新设置最大长度
func (q *queue[T]) ReSetLen(Len int) {
	if Len < 0 {
		panic("queue length must be greater than 0")
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.index == 0 {
		if len(q.data) > q.len {
			q.data = q.data[:q.len]
		}
		q.len = Len
		return
	}

	if q.len > Len {
		arr := make([]T, Len)
		if q.index+1 >= Len {
			copy(arr, q.data[q.index-Len:q.index])
		} else {
			copy(arr, q.data[q.index:])
			copy(arr[q.len-q.index:], q.data[q.len-Len+q.index:])
		}
		q.data = arr
		q.index = 0
		q.len = Len
	} else if q.len < Len {
		q.data = append(q.data, make([]T, Len-q.len)...)
		copy(q.data[q.index+Len-q.len:], q.data[q.index:])
		q.len = Len
	}
}

// Range 遍历
func (q *queue[T]) Range(f func(T)) {
	for _, v := range q.data {
		f(v)
	}
}

// IsEmpty 判空
func (q *queue[T]) IsEmpty() bool {
	return len(q.data) == 0
}
