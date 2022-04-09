# ttcp

[cilium/ebpf/examples/tcprtt](https://github.com/cilium/ebpf/tree/master/examples/tcprtt)에다가 `BPF_MAP_TYPE_HASH`를 추가해서 srcAddr을 필터링하도록 수정함. 어플리케이션에는 Http API를 추가해서 필터링할 IP를 추가, 삭제, 조회할 수 있도록 함.

## Getting Startd

Clone `ttcp` to `$GOPATH/src/github.com/KumKeeHyun/ttcp` 

```
# 1. ebpf-builder를 통해 bpf 프로그램 컴파일. bpf_bpfel.go, bpf_bpfel.o 생성.
# 2. 도커 이미지 빌드. kbzjung359/ttcp 이미지 생성
make build
```

## Running ttcp

`Docker hub`에 `kbzjugn359/ttcp:v0.0.0` 이미지를 올려둠.

```
# kbzjugn359/ttcp:v0.0.0 이미지 실행
# 8090 포트 연결
make run
```

Http API를 통해 필터링할 IP 테이블(BPF_MAP) 조작.

```
# 테이블에 등록된 IP 조회
curl -L http://localhost:8090 -XGET

# 테이블에 필터링할 새로운 IP 추가
curl -L http://localhost:8090 -XPUT -d 192.168.0.1

# 테이블에 필터링하지 않을 IP 제거
curl -L http://localhost:8090 -XDELETE -d 192.168.0.1
```

