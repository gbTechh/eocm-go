package suma

func Add(a,b int) int  {
	return a + b
}

func AddMultiple(numbers ...int) int {
	resp := 0
	for _, v := range numbers {
		resp += v
	}
	return resp
}