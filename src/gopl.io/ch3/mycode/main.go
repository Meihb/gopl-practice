package main

import (
	"fmt"
	"math"
	"unicode/utf8"
)

func main() {
	o := 0666 //0开头表示8进制数
	/*
		通常Printf格式化字符串包含多个%参数时将会包含对应相同数量的额外操作数，但是%之后的[1]副词告诉Printf函数再次使用第一个操作数。第二，%后的#副词告诉Printf在用%o、%x或%X输出时生成0、0x或0X前缀。
	*/
	fmt.Printf("%d %[1]o %#[1]o\n", o) //printf %参数时可以尝试用[1]表示复用第一个操作数
	fmt.Printf("%d %o %#o\n", o, o, o) //上面两式相等

	x := 0xdeadbeef //0x十六进制
	fmt.Printf("%d\n", x)
	fmt.Println(uint32(x))

	//const 常量表达式的值在编译期计算,所以一定不要在表达式中包含runtime运行时才能算出来来的变量啥的,这一点和java一样,且都是基础类型,boolean,string或数字

	//iota
	const (
		_   = 1 << (10 * iota)
		KiB // 1024
		MiB // 1048576
		GiB // 1073741824
		TiB // 1099511627776             (exceeds 1 << 32)
		PiB // 1125899906842624
		EiB // 1152921504606846976
		ZiB // 1180591620717411303424    (exceeds 1 << 64)
		YiB // 1208925819614629174706176
	)
	// fmt.Println(KiB, MiB, GiB, YiB) YIB overflows

	//无类型常量,常量无需指定基础类型,那么其被提供了至少256bit的运算精度,无符号布尔、整形、浮点型、字符、复数、字符串
	//比如上述例子中,YiB远超任何go中整形范围,但是依然合法(这不代表你在runtime中就能输出来 )，但是输出YiB/ZiB(1024)则是可以的

	var x1 float32 = math.Pi
	var y float64 = math.Pi
	var z complex128 = math.Pi
	fmt.Println(x1, y, z) //注意到x1,y,z三个数都是从math.pi做表达式的,但是却不用强制类型转换,因为后者是无类型的

	//常量被赋值给变量时,会进行隐式的类型转换
	var f1 float64 = 3 + 0i // untyped complex -> float64
	f1 = 2                  // untyped integer -> float64
	f1 = 1e123              // untyped floating-point -> float64
	f1 = 'a'                // untyped rune -> float64
	fmt.Println(f1)

	var rune1 rune = 2
	fmt.Println(rune1) // \x 表示十六进制,后面跟两位,表示单字节编码;\u表示unicode码.四个十六进制
	fmt.Printf("%c\n", '\u0061')

	/*
	   unicode
	   0xxxxxxx                             runes 0-127    (ASCII)
	   110xxxxx 10xxxxxx                    128-2047       (values <128 unused)
	   1110xxxx 10xxxxxx 10xxxxxx           2048-65535     (values <2048 unused)
	   11110xxx 10xxxxxx 10xxxxxx 10xxxxxx  65536-0x10ffff (other values unused)
	   前缀各不相同,每个unicode编码都不会是别的unicode编码的子串(!太重要了)
	*/
	s := "Hello, 世界"
	fmt.Println(len(s))                    // "13" 个字节长度  1*7(Hello, )+3*2(世界)
	fmt.Println(utf8.RuneCountInString(s)) // "9" 个unicode码
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%d\t%c\n", i, r)
		i += size
	}

}

/*得益于utf8的编码设计,首先UTF8编码比较紧凑，完全兼容ASCII码，并且可以自动同步：它可以通过向前回朔最多3个字节就能确定当前字符编码的开始字节的位置。
它也是一个前缀编码，所以当从左向右解码时不会有任何歧义也并不需要向前查看（译注：像GBK之类的编码，如果不知道起点位置则可能会出现歧义）。没有任何字符
的编码是其它字符编码的子串，或是其它编码序列的字串，因此搜索一个字符时只要搜索它的字节编码序列即可，不用担心前后的上下文会对搜索结果产生干扰。同时
UTF8编码的顺序和Unicode码点的顺序一致，因此可以直接排序UTF8编码序列。同时因为没有嵌入的NUL(0)字节，可以很好地兼容那些使用NUL作为字符串结尾的编程
语言。
*/
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
