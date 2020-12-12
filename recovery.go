package pbw

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

var internalServerErr = "Internal Server Error"

// Recovery middleware used error recover.
func Recovery() HandlerFunc {
	return func(c Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n", trace(message))
				c.Data(http.StatusInternalServerError, []byte(internalServerErr))
			}
		}()

		c.Next()
	}
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTrace:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
