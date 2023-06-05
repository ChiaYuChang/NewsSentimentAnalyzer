package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Logger struct {
	Verbose     bool
	RandomDelay bool
}

func (p Logger) Println(a ...any) {
	if p.Verbose {
		fmt.Println(a...)
	}
}

func (p Logger) Printf(format string, a ...any) {
	if p.Verbose {
		fmt.Printf(format, a...)
	}
}

func (p Logger) Print(a ...any) {
	if p.Verbose {
		fmt.Print(a...)
	}
	if p.RandomDelay {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	}
}

func (p Logger) Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}
