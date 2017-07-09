package main

/*

#include <stdio.h>

void hello() {
	printf("hello world!\r\n")
}


*/

import "C"

func Hello() {
	C.hello()
}

func main() {
	Hello()
}
