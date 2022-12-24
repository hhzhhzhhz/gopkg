package runtime

import (
	"fmt"
	"github.com/hhzhhzhhz/gopkg/log"
	"net/http"
	_ "net/http/pprof"
)

const Pprof = ":6060"

func StartPprof(addr string) error {
	log.Logger().Info(fmt.Sprintf("pprof is listening and serving on %s", addr))
	go func() {
		fmt.Errorf("%s", http.ListenAndServe(addr, nil))
	}()
	return nil
}
