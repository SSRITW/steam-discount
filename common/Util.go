package common

import (
	"strconv"
	"fmt"
)


func StrToInt(str string) (result int){
	var err error
	result ,err = strconv.Atoi(str)
	if err!=nil {
		fmt.Println(err.Error())
	}
	return result
}

func StrToFloat64(str string) (result float64)  {
	var err error
	result,err =strconv.ParseFloat(str,32)
	if err!=nil {
		fmt.Println(err.Error())
	}
	return result
}

