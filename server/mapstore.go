package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/KumKeeHyun/ttcp/internal"
	"github.com/cilium/ebpf"
)

type MapStore interface {
	GetAll() ([]string, error)
	Put(ipStr string) error
	Delete(ipStr string) error
}

type mapOp func(m *ebpf.Map) error
type mapOpReq struct {
	mapOp  mapOp
	doneCh chan error
}

func NewMapStore(ctx context.Context, bpfMap *ebpf.Map) MapStore {
	mapStore := &mapStore{
		bpfMap: bpfMap,
		reqCh:  make(chan mapOpReq, 1),
	}
	go mapStore.startHandle(ctx)

	return mapStore
}

type mapStore struct {
	bpfMap *ebpf.Map
	reqCh  chan mapOpReq
}

var _ MapStore = &mapStore{}

func (s *mapStore) startHandle(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("mapStore: received done. exiting...")
			return
		case req := <-s.reqCh:
			s.handleReq(req)
		}
	}
}

func (s *mapStore) handleReq(req mapOpReq) {
	defer close(req.doneCh)
	if err := req.mapOp(s.bpfMap); err != nil {
		req.doneCh <- err
	}
}

func (s *mapStore) sendReq(op mapOp) error {
	done := make(chan error)
	s.reqCh <- mapOpReq{
		mapOp:  op,
		doneCh: done,
	}
	return <-done
}

func (s *mapStore) GetAll() ([]string, error) {
	res := make([]string, 0)

	err := s.sendReq(func(m *ebpf.Map) error {
		var ipInt, notUsed uint32
		iter := m.Iterate()

		for iter.Next(&ipInt, &notUsed) {
			ip := internal.IntToIP(ipInt)
			res = append(res, ip.String())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *mapStore) Put(ipStr string) error {
	ipInt, err := ipStrToipInt(ipStr)
	if err != nil {
		return fmt.Errorf("failed to put %s: %s", ipStr, err)
	}

	err = s.sendReq(func(m *ebpf.Map) error {
		return m.Put(&ipInt, &ipInt)
	})
	if err != nil {
		return fmt.Errorf("failed to put %d: %s", ipInt, err)
	}
	return nil
}

func ipStrToipInt(ipStr string) (uint32, error) {
	ip, err := internal.StringToIPv4(ipStr)
	if err != nil {
		return 0, err
	}
	return internal.IPToInt(ip), nil
}

func (s *mapStore) Delete(ipStr string) error {
	ipInt, err := ipStrToipInt(ipStr)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %s", ipStr, err)
	}

	err = s.sendReq(func(m *ebpf.Map) error {
		return m.Delete(&ipInt)
	})
	if err != nil && !errors.Is(ebpf.ErrKeyNotExist, err) {
		return fmt.Errorf("failed to delete %s: %s", ipStr, err)
	}

	return nil
}
