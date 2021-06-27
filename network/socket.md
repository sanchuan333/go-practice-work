总结几种 socket 粘包的解包方式: fix length/delimiter based/length field based frame decoder。尝试举例其应用

## fix length
规定每个数据包的长度，用占位符补足长度。
接收方根据规定的长度来判断解包。如果发送数据太小会补足长度，增加了不必要的数据

## delimiter based
规定数据包结尾的分割符号，如换行符。
接收方根据分割符号来判断包结束。

## length field based frame decoder
在数据包中加上一个表示长度的字段，解包根据这个长度进行解析