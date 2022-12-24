package election

import (
	"context"
	"fmt"
	"github.com/hhzhhzhhz/gopkg/log"
	"github.com/hhzhhzhhz/gopkg/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sync"
)

const (
	ttl = 5
)

type OptionEtcdEl struct {
	context.Context
	Namespace string
	LeaderId  string
	*clientv3.Client
}

func (o *OptionEtcdEl) valify() {

}

func NewElection(opt *OptionEtcdEl) *electionEtcd {
	if opt.LeaderId == "" {
		opt.LeaderId = utils.UUID()
	}
	return &electionEtcd{
		OptionEtcdEl: opt,
		unmap:        make(map[FunType][]func(), 0),
	}
}

type electionEtcd struct {
	*OptionEtcdEl
	mux            sync.Mutex
	isLeader       bool
	leaderIdentity string
	// Have you ever been a leader
	once  bool
	unmap map[FunType][]func()
	se    *concurrency.Session
	watch <-chan clientv3.GetResponse
}

func (e *electionEtcd) Register(ft FunType, f func()) {
	fs, ok := e.unmap[ft]
	if ok {
		fs = append(fs, f)
		e.unmap[ft] = fs
		return
	}
	e.unmap[ft] = []func(){f}
}

func (e *electionEtcd) Run() error {
	log.Logger().Info("Election.Start starting")
	se, err := concurrency.NewSession(
		e.Client,
		concurrency.WithTTL(ttl),
		concurrency.WithContext(e.Context),
	)
	if err != nil {
		return err
	}
	el := concurrency.NewElection(se, e.Namespace)
	e.watch = el.Observe(e.Context)
	utils.Wrap(func() {
		for {
			select {
			case w, ok := <-e.watch:
				if !ok {
					//log.Logger().Error("Election.Start Stop chan closed")
					return
				}
				if len(w.Kvs) == 0 {
					continue
				}
				e.onNewLeader(string(w.Kvs[0].Value))
			}
		}
	})
	if err := el.Campaign(e.Context, e.LeaderId); err != nil {
		return err
	}
	log.Logger().Info("Election.Start started")
	return nil
}

func (e *electionEtcd) onNewLeader(leader string) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.leaderIdentity = leader
	// we're notified when new leader elected
	log.Logger().Info(fmt.Sprintf("Election.OnNewLeader isleader=%t leader=%s current=%s", leader == e.LeaderId, leader, e.LeaderId))
	if leader == e.LeaderId {
		e.once = true
		e.isLeader = true
		fs, ok := e.unmap[Leader]
		if ok {
			for _, fun := range fs {
				fun()
			}
		}
		return
	}
	e.isLeader = false
	fs, ok := e.unmap[Follower]
	if ok && e.once {
		for _, fun := range fs {
			fun()
		}
	}
}

func (e *electionEtcd) IsLeader() (bool, string) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.isLeader {
		return true, e.leaderIdentity
	}
	return false, e.leaderIdentity
}

func (e *electionEtcd) Close() error {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.se.Close()
}
