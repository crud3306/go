golang类型转换
===============


// 整形转字符串
// ------------
fmt.Println(strconv.Itoa(100))
//该方法的源码是：
// Itoa is shorthand for FormatInt(i, 10).
func Itoa(i int) string {
	return FormatInt(int64(i), 10)
}
// 可以看出是FormatInt方法的简单实现。



// 字符串转整形
// ------------
i, _ := strconv.Atoi("100")
fmt.Println(i)



// 64位整形转字符串
// ------------
var i int64
i = 0x100
fmt.Println(strconv.FormatInt(i, 10))
//FormatInt第二个参数表示进制，10表示十进制。


// 字节转字符串
// ------------
fmt.Println(string([]byte{97, 98, 99, 100}))


// 字符串转字节
// ------------
fmt.Println([]byte("abcd"))



// 字节转32位整形
// ------------
b := []byte{0x00, 0x00, 0x03, 0xe8}
bytesBuffer := bytes.NewBuffer(b)

var x int32
binary.Read(bytesBuffer, binary.BigEndian, &x)
fmt.Println(x)
其中binary.BigEndian表示字节序，相应的还有little endian。通俗的说法叫大端、小端。



// 32位整形转字节
// ------------
var x int32
x = 106
bytesBuffer := bytes.NewBuffer([]byte{})
binary.Write(bytesBuffer, binary.BigEndian, x)
fmt.Println(bytesBuffer.Bytes())




// 继续补充...


