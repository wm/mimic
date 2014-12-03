package utils

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"os"
)

//PrintError is wraps fmt.Fprintf(os.Stderr, message, a...) and changes the
//terminal output color to Red
func PrintError(message string, a ...interface{}) {
	ct.ChangeColor(ct.Red, true, ct.None, false)
	fmt.Fprintf(os.Stderr, message, a...)
	ct.ResetColor()
}

//PrintError is wraps fmt.Fprintf(os.Stderr, message, a...) and changes the
//terminal output color to Green
func PrintOk(message string, a ...interface{}) {
	ct.ChangeColor(ct.Green, true, ct.None, false)
	fmt.Printf(message, a...)
	ct.ResetColor()
}

//PrintError is wraps fmt.Fprintf(os.Stderr, message, a...) and changes the
//terminal output color to Yellow
func PrintNotice(message string, a ...interface{}) {
	ct.ChangeColor(ct.Yellow, true, ct.None, false)
	fmt.Printf(message, a...)
	ct.ResetColor()
}
