package Models

import "TradeSimulator/Models/Enum"

type PriorityQueue []*Order

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	p1 := pq[i]
	p2 := pq[j]

	samePriceCompare := p1.Price == p2.Price && p1.Timestamp.Before(p2.Timestamp)
	if p1.OrderType == Enum.Buy {
		if p1.IsMarketPrice || p2.IsMarketPrice {
			return true
		} else {
			return p1.Price > p2.Price || samePriceCompare
		}
	}
	if p1.OrderType == Enum.Sell {
		if p1.IsMarketPrice || p2.IsMarketPrice {
			return true
		} else {
			return (p1.Price < p2.Price) || samePriceCompare
		}
	}
	return false
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
