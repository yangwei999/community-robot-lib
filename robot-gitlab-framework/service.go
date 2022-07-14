package framework

import (
	"fmt"
	"github.com/opensourceways/community-robot-lib/interrupts"
	"net/http"
	"strconv"
	"time"
)

func Run(robot interface{}, port int, timeout time.Duration) error {
	h := handlers{}
	h.registerHandler(robot)

	hs := h.getHandler()
	if len(hs) == 0 {
		return fmt.Errorf("it is not a robot")
	}

	d := dispatcher{h: hs}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.OnInterrupt(func() {
		d.Wait()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	http.Handle("/gitlab-hook", &d)

	httpServer := &http.Server{Addr: ":" + strconv.Itoa(port)}

	interrupts.ListenAndServe(httpServer, timeout)

	return nil
}
