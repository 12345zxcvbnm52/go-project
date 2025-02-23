package queue

// 队列的默认最小长度
// 必须是 2 的幂,以便进行位运算取模:x % n == x & (n - 1)
const minQueueLen = 16

// queue提供了一个基于快速环形缓冲区的队列
// 注意该Queue不是线程安全的
type Queue[V any] struct {
	buf               []*V
	head, tail, count int
}

func NewQueue[V any]() *Queue[V] {
	return &Queue[V]{
		buf: make([]*V, minQueueLen),
	}
}

func (q *Queue[V]) Length() int {
	return q.count
}

// resize将队列调整为恰好两倍于当前内容的大小
// 如果队列未达到一半容量,这可能会导致队列缩小
func (q *Queue[V]) resize() {
	newBuf := make([]*V, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// 将元素添加到队列的末尾
func (q *Queue[V]) Push(elem V) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = &elem
	// 位运算取模
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

// Peek返回队列头部的元素,如果队列为空,调用此方法会引发panic
func (q *Queue[V]) Head() V {
	if q.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return *(q.buf[q.head])
}

// Get返回队列中索引为i的元素,如果索引无效,调用此方法会引发panic
// 此方法接受正负索引值(类似于redis),索引0表示第一个元素,索引-1表示最后一个元素
func (q *Queue[V]) Get(i int) V {
	// 如果索引为负数,则转换为正索引
	if i < 0 {
		i += q.count
	}
	if i < 0 || i >= q.count {
		panic("queue: Get() called with index out of range")
	}
	// 位运算取模
	return *(q.buf[(q.head+i)&(len(q.buf)-1)])
}

// Remove 移除并返回队列头部的元素,如果队列为空,调用此方法会引发panic
func (q *Queue[V]) Pop() V {
	if q.count <= 0 {
		panic("queue: Remove() called on empty queue")
	}
	ret := q.buf[q.head]
	q.buf[q.head] = nil
	// 位运算取模
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// 如果缓冲区容量为 1/4,则缩小容量
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return *ret
}
