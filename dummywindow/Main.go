// Heavily based on the example on https://github.com/go-gl/glfw. 
// That code is under a 3-clause bsd licence. See the file BSD in the root of this project
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Expect one argument: Window title")
	} else {

		if err := glfw.Init(); err != nil {
			panic(err)
		} else {
			defer glfw.Terminate()
			if window, err := glfw.CreateWindow(640, 480, os.Args[1], nil, nil); err != nil {
				panic(err)
			} else {
				window.MakeContextCurrent()
				for !window.ShouldClose() {
					window.SwapBuffers()
					glfw.PollEvents()
				}
			}
		}
	}
}
