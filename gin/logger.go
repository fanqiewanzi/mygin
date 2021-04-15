package gin

import (
	"log"
	"time"
)

//中间件logger
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.R.RequestURI, time.Since(t))
	}
}
