package queue_test

import (
	"testing"

	"github.com/mhqdz/queue"
)

func TestMain(t *testing.T) {
	q1 := queue.New[string](5)
	q2 := queue.NewBySlice[int]([]int{0, 1, 2, 3, 4})
	t.Log("queue1, 原数据:", q1.GetData())
	t.Log("isEmpty:", q1.IsEmpty())

	for i := 0; i < 10; i++ {
		q1.Append(string([]byte{byte('a' + i)}))
		arr := []string{}
		q1.Range(func(s string) {
			arr = append(arr, s)
		})
		t.Log(arr)
	}

	t.Log("queue2, 原数据:", q2.GetData())
	t.Log("isEmpty:", q2.IsEmpty())
	for i := 0; i < 10; i++ {
		q2.Append(5 + i)
		t.Logf("slice:%v,data:%v\n", q2.GetSlice(), q2.GetData())
	}
}

func TestReLen(t *testing.T) {
	q := queue.NewBySlice[any]([]any{0, 1, 2, struct{ int }{int: 3}})
	q.Append(4)
	t.Log("原数据: ", q.GetData())

	q.ReSetLen(10)
	t.Log("扩容到10之后", q.GetData())

	q.Append(5)
	t.Log("push尝试", q.GetData())

	q.ReSetLen(2)
	t.Log("缩容到2之后", q.GetData())

	q.Append(6)
	t.Log("push尝试", q.GetData())
}
