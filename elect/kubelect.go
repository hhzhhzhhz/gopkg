package election

import (
	"context"
	"github.com/google/uuid"
	"github.com/hhzhhzhhz/gopkg/log"
	"github.com/hhzhhzhhz/gopkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"sync"
	"time"
)

type OptionKubeEl struct {
	context.Context
	Namespace string
	LeaderId  string
	kubernetes.Interface
}

func NewK8sElection(opt *OptionKubeEl) *electionKube {
	if opt.LeaderId == "" {
		opt.LeaderId = uuid.New().String()
	}
	return &electionKube{
		OptionKubeEl: opt,
		unmap:        make(map[FunType][]func(), 0),
	}
}

type electionKube struct {
	*OptionKubeEl
	isLeader bool
	// Have you ever been a leader
	once           bool
	unmap          map[FunType][]func()
	mux            sync.Mutex
	leaderIdentity string
	kubeCli        kubernetes.Interface
}

func (e *electionKube) Register(ft FunType, f func()) {
	fs, ok := e.unmap[ft]
	if ok {
		fs = append(fs, f)
		e.unmap[ft] = fs
		return
	}
	e.unmap[ft] = []func(){f}
}

func (e *electionKube) Run() error {
	utils.Wrap(func() {
		leaseLock := &resourcelock.LeaseLock{
			LeaseMeta: metav1.ObjectMeta{
				Name:      "job",
				Namespace: e.Namespace,
			},
			Client: e.kubeCli.CoordinationV1(),
			LockConfig: resourcelock.ResourceLockConfig{
				Identity: e.LeaderId,
			},
		}
		leaderelection.RunOrDie(e.Context, leaderelection.LeaderElectionConfig{
			Lock: leaseLock,
			// IMPORTANT: you MUST ensure that any code you have that
			// is protected by the lease must terminate **before**
			// you call cancel. Otherwise, you could have a background
			// loop still running and another process could
			// get elected before your background loop finished, violating
			// the stated goal of the lease.
			ReleaseOnCancel: true,
			LeaseDuration:   10 * time.Second,
			RenewDeadline:   8 * time.Second,
			RetryPeriod:     5 * time.Second,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {
					//log.Logger().Info("Election.OnStartedLeading identity=%s", e.identity)
					// we're notified when we start - this is where you would
					// usually put your code
				},
				OnStoppedLeading: func() {
					//log.Logger().Info("Election.OnStoppedLeading identity=%s", e.identity)
					// we can do cleanup here
					e.mux.Lock()
					e.isLeader = false
					e.mux.Unlock()
				},
				OnNewLeader: func(leader string) {
					e.mux.Lock()
					defer e.mux.Unlock()
					e.leaderIdentity = leader
					// we're notified when new leader elected
					log.Logger().Info("Election.OnNewLeader isleader=%t leader=%s current=%s ", leader == e.LeaderId, leader, e.LeaderId)
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
				},
			},
		})
	})
	return nil
}

func (e *electionKube) IsLeader() (bool, string) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.isLeader {
		return true, e.leaderIdentity
	}
	return false, e.leaderIdentity
}

func (e *electionKube) Close() error {
	e.mux.Lock()
	defer e.mux.Unlock()
	return nil
}
