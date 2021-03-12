package main

import (
	"fmt"
	"math"
)

//测试一下同一个类型,如果分别创建同名的receiver指针型/非指针型是否可以

type bridge struct {
	x, y int
}

//non-pointer receiver
func (bridge bridge) scaleBy(factor int) bridge {
	// bridgeAddr := &bridge
	// bridgeAddr.x *= factor
	// bridgeAddr.y *= factor

	bridge.x *= factor
	bridge.y *= factor
	fmt.Printf("in non-pointer receiver:%#v", bridge) //in non-pointer receiver:main.bridge{x:4, y:6}
	return bridge
}

//看来是不允许的,这很正常
// func (bridge *bridge) scaleBy(factor int)(bridge){

// }

func (bridge *bridge) scaleBy2(factor int) {

	bridge.x *= factor
	bridge.y *= factor
	fmt.Printf("in non-pointer receiver:%#v \n", bridge) //in non-pointer receiver:main.bridge{x:4, y:6}

}

// An IntList is a linked list of integers.
// A nil *IntList represents the empty list.
type IntList struct {
	Value int
	Tail  *IntList
}

// Sum returns the sum of the list elements.
func (list *IntList) Sum() int {
	if list == nil {
		return 0
	}
	return list.Value + list.Tail.Sum()
}

/*values测试slice这类在receiver中的改变*/
type values map[string][]string

func (v values) get(key string) string {
	if vs := v[key]; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (v values) add(key, value string) {
	v[key] = append(v[key], value) //append这个操作是要赋值的,还记得吗
}

func main() {

	defer func() {
		p := recover()
		fmt.Println(p)
	}()
	bridge := bridge{
		x: 2,
		y: 3,
	}
	bridge.scaleBy(2)
	fmt.Println(bridge) //{2 3} copy而非引用

	bridge.scaleBy2(3)
	fmt.Println(bridge) //{6 9} 引用,指向同一块内存地址

	var l1 IntList
	fmt.Println(&l1 == nil)

	m := values{"lang": {"en"}} // direct construction
	m.add("item", "1")
	m.add("item", "2")
	fmt.Println(m)
	fmt.Println(m.get("lang")) // "en"
	fmt.Println(m.get("q"))    // ""
	fmt.Println(m.get("item")) // "1"      (first value)
	fmt.Println(m["item"])     // "[1 2]"  (direct map access)

	m = nil
	fmt.Printf("type is %T,%[1]v", m)
	fmt.Println(m.get("item")) // ""
	// m.add("item", "3")         // panic: assignment to entry in nil map

	/*
	 method 和 method expression 方法表达式 、 method value 方法值
	*/
	p := point{1, 2}
	q := point{4, 6}

	//method value
	distanceFromP := p.distance
	fmt.Printf("type is %T,%v \n", distanceFromP, p) //type is func(main.point) float64
	fmt.Println(distanceFromP(q))

	//method expression 方法表达式是特殊的方法值,其并未绑定具体的类型实例,而是类型本身,因此传参还需要receiver本身
	distanceOfPoint := point.distance
	fmt.Printf("type is %T \n", distanceOfPoint) //type is func(main.point, main.point) float64
	fmt.Println(distanceOfPoint(p, q))

	//因此你甚至可以 申明一个method expression 变量
	var methodExpre func(p, q point) float64
	methodExpre = point.distance
	fmt.Println(methodExpre(p, q))
	methodExpre = distance
	fmt.Println(methodExpre(p, q)) //感觉这样就把tradition function和method 统一起来了,is this what u want?

}

type point struct {
	x, y float64
}

//traditional function
func distance(p point, q point) float64 {
	return math.Hypot(q.x-p.x, q.y-p.y)
}

//method of the point type
func (p point) distance(q point) float64 {
	return math.Hypot(q.x-p.x, q.y-p.y)
}
