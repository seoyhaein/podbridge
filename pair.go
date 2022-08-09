package podbridge

// 대문자로 할지 고민해보자.
type Pair struct {
	Key   any
	Value any
}

func MakePair(a any, b any) *Pair {
	return &Pair{
		Key:   a,
		Value: b,
	}
}

// 이렇게 만드는 것 고민해보자. 많이 사용하는 field에 대해서
func Name(name string) *Pair {
	return &Pair{
		Key:   "Name",
		Value: name,
	}
}
