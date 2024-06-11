package Models

import "TradeSimulator/Models/Enum"

type PriorityQueue []*Order

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	samePriceCompare := pq[i].Price == pq[j].Price && pq[i].Timestamp.Before(pq[j].Timestamp)
	if pq[i].OrderType == Enum.Buy {
		return (pq[i].Price > pq[j].Price) || samePriceCompare
	}
	return (pq[i].Price < pq[j].Price) || samePriceCompare
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	order := x.(*Order)
	*pq = append(*pq, order)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	order := old[n-1]
	*pq = old[0 : n-1]
	return order
}
