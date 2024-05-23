package main

type order struct {
	Value int
}

func toOrder(res int) (order, error) {
	return order{
		Value: res,
	}, nil
}

func square(value int) (int, error) {
	return value * value, nil
}

func plusone(value int) (int, error) {
	return value + 1, nil
}

func timesnine(value int) (int, error) {
	return value * 9, nil
}
