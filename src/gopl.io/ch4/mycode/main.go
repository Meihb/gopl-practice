package main

import (
	"crypto/sha256"
	"fmt"
	"sort"
)

func main() {
	//go中数组是一个固定长度的序列,和php中变长的概念不一致
	a1 := [...]int{1, 2, 3}
	fmt.Printf("%T %[1]v \n", a1) //[3]int [1 2 3]

	//这居然是一个[100]int的数组你敢信,所以他是按照key来计算的，array的key严格按照从0到len()-1
	a2 := [...]int{99: -1}
	fmt.Printf("%T %[1]v \n", a2) //[3]int [1 2 3]

	//array比较 以下c1 c2是[32]byte即256bit的数组,只有满足每个key对应value都相等的情况下才==true
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))
	// %x表示用16精致打印, %t表示 boolen型,%T 表示数据类型
	fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)
	// Output:
	// 2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881
	// 4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
	// false
	// [32]uint8

	/*在函数传参处理中,大多数其他语言据我所知都是把array做引用传值对象的,但是go中相对应的概念是slice切片,array是普通数据格式,是copy处理的
	  感觉他们有种老学究式的坚守呀,虽知道且承认变长的slice更加风靡,却硬是照本宣科造出来个注定不会被大量食用的原教旨array出来,哈哈
	*/
	a3 := [4]int{1, 2, 3, 4}
	fmt.Printf("%T\n", &a3)
	a4 := zeroNotCopy(a3)
	a5 := zero(&a3)
	fmt.Printf("%v %v %v\n", a3, a4, a5)

	/*
					   slcie是由三个部分组成,hashtable都一样吧,不对,就是数组,不是hasntable,指针(ptr),长度(len),容量(capacity)
				      ptr指向slice第一个元素对应的底层数组的地址,却不一定是数组的第一个元素

		            slice概念上类似php的数组,但是还记得他们的key固定从0开始对吧,那就是说如果你如下months定义一个切片从key=1开始,key=0处会被零值化
		            所以len=cap=13 记住这一点
	*/

	months := []string{
		1:  "January",
		2:  "February",
		3:  "March",
		4:  "April",
		5:  "May",
		6:  "June",
		7:  "July",
		8:  "August",
		9:  "September",
		10: "October",
		11: "November",
		12: "December",
	}
	fmt.Printf("%#v %d %d\n", months, len(months), cap(months))
	months2 := months[1:]
	fmt.Printf("%#v %d %d\n", months2, len(months2), cap(months2))

	Q2 := months[4:7]
	summer := months[6:9]
	fmt.Println(Q2)     // ["April" "May" "June"]
	fmt.Println(summer) // ["June" "July" "August"]

	// fmt.Println(summer[:20]) // panic: out of range

	endlessSummer := summer[:5] // extend a slice (within capacity)
	/*两个指针是一样的,i.e.底层共享者同一个底层数组,只是切片后生成的新的slice指向自己元素所对应的底层数组地址,这也是上面slice不一定指向底层数组首部的解释*/
	fmt.Printf("%p %p \n", &months[6], &summer)
	fmt.Println(endlessSummer) // "[June July August September October]" 看出来这是指针的本质了吧

	var s []int           // len(s) == 0, s == nil
	s = nil               // len(s) == 0, s == nil
	s = []int(nil)        // len(s) == 0, s == nil
	s = []int{}           // len(s) == 0, s != nil
	s = make([]int, 2, 3) //len(s)=2,s!=nil
	fmt.Println(s)

	/*
	   append
	*/
	var runes []rune
	for _, r := range "Hello, 世界" {
		runes = append(runes, r) //为什么我们用append的返回值给runes呢,原因在于如果发生了hash扩容,之前的指针就不是指向正确的数组了,切记切记
	}
	fmt.Printf("%q %[1]v\n", runes) // "['H' 'e' 'l' 'l' 'o' ',' ' ' '世' '界']"

	/*
		更新slice变量不仅对调用append函数是必要的，实际上对应任何可能导致长度、容量或底层数组变化的操作都是必要的。要正确地使用slice，
		需要记住尽管底层数组的元素是间接访问的，但是slice对应结构体本身的指针、长度和容量部分是直接访问的。要更新这些信息需要像上面例子那样一个显式的赋值操作。
		从这个角度看，slice并不是一个纯粹的引用类型，它实际上是一个类似下面结构体的聚合类型：
		type IntSlice struct {
		    ptr      *int
		    len, cap int
		}
	*/

	var x []int
	x = append(x, 1)
	x = append(x, 2, 3)
	x = append(x, 4, 5, 6)
	x = append(x, x...) // append the slice x 这个使用方法还算熟悉吧,变长参数
	fmt.Println(x)      // "[1 2 3 4 5 6 1 2 3 4 5 6]"

	/*
						map 无法对map的elem进行取址操作,因为map经常需要分配更大的内存空间,从而导致之前的地址无效
				      且map的的迭代顺序是不确定的,不同的哈希函数可能导致不同的遍历顺序,在实践中，遍历的顺序是随机的，每一次遍历的顺序都不相同。
				      这是故意的，每次都使用随机的遍历顺序可以强制要求程序不会依赖具体的哈希函数实现。如果要按顺序遍历key/value对，我们必须显式地对key进行排序，
				      可以使用sort包的Strings函数对字符串slice进行排序

		            按照数组,slice的规则,为什么map不能直接用[keyType][valueType]呢,想象一下,如果遇到多维map  map[int]map[int][string]写成[int][int][string]好像也无不妥吧???
	*/

	ages := make(map[string]int) // mapping from strings to ints
	ages = map[string]int{
		"alice":   31,
		"charlie": 34,
	}
	ages["Bob"] = 2 //不可以用单引号哦,单引号是rune类型,也就是说要么int,要么是unicode码或者值(H什么的)
	ages["Bob"]++
	fmt.Println(ages)
	//排序
	var names []string
	for name := range ages {
		names = append(names, name) //记得一定要赋值
	}
	sort.Strings(names) //好坑！
	for _, name := range names {
		fmt.Printf("%s\t%d\n", name, ages[name])
	}
	//仅仅申明的map是零值化的,也就是说没有指向任何哈希表,这个时候是无法之间加入元素的
	var ages2 map[string]int
	fmt.Println(ages2 == nil)    // "true"
	fmt.Println(len(ages2) == 0) // "true"
	//ages2["Bob"] = 2 panic: assignment to entry in nil map
	ages2 = make(map[string]int)
	ages2["Bob"] = 2
	fmt.Println(ages2)
	_, ok := ages2["Bob1"] //尤其两个ages对比时,一定要用到第二个参数
	if !ok {
		fmt.Println("not a key")
	}

	addEdge("A", "B")
	fmt.Println(hasEdge("A", "C"))

	/*
	   struct
	*/
}

func zero(ptr *[4]int) *[4]int {
	for i := range ptr {
		ptr[i] = 0 //==(*ptr)[i] = 0 相当于语法糖吧
	}
	return ptr
}

func zeroNotCopy(ptr [4]int) [4]int {
	for i := range ptr {
		ptr[i]++
	}
	return ptr
}

//邻接表
var graph = make(map[string]map[string]bool)

func addEdge(from, to string) {
	edges := graph[from]
	if edges == nil {
		edges = make(map[string]bool)
		graph[from] = edges
	}
	edges[to] = true
}

func hasEdge(from, to string) bool {
	return graph[from][to]
}
