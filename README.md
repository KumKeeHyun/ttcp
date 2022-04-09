# ttcp

[cilium/ebpf/examples/tcprtt](https://github.com/cilium/ebpf/tree/master/examples/tcprtt)에다가 `BPF_MAP_TYPE_HASH`를 추가해서 srcAddr을 필터링하도록 수정함. 어플리케이션에는 Http API를 추가해서 필터링할 IP를 추가, 삭제, 조회할 수 있도록 함.

## Getting Startd

도커만 준비하면 됨. Clone `github.com/KumKeeHyun/ttcp`

```
# 1. ebpf-builder를 통해 bpf 프로그램 컴파일. bpf_bpfel.go, bpf_bpfel.o 생성.
# 2. 도커 이미지 빌드. kbzjung359/ttcp 이미지 생성
make build
```

## Running ttcp

`Docker hub`에 `kbzjugn359/ttcp:v0.0.0` 이미지를 올려둠. 따로 빌드하지 않아도 됨.

```
sudo docker run --rm --privileged \
	-p 8090:8090 \
	kbzjugn359/ttcp:v0.0.0`
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

## example

호스트의 network namespace와 도커 컨테이너의 network namespace에서 각각 http 요청을 보내보는 시나리오임.

`ttcp` 컨테이너는 `172.17.0.2`에서 실행되고 있음. 원래대로라면 다른 네임스페이스의 TCP 요청을 보지 못하지만 BPF를 통해 커널에서 TCP 요청을 감시하기 때문에 호스트의 `10.178.0.5` 요청도 감시할 수 있음. 물론 다른 컨테이너인 `172.17.0.3` 요청도 감시할 수 있음. 

1. ttcp 실행
2. Src Addr IP 필터링 테이블에 호스트 IP인 10.178.0.5 추가
3. 호스트에서 curl 요청 -> ttcp에 요청이 찍힘
4. 테이블에 호스트 IP 제거
5. 호스트에서 curl 요청 -> ttcp에서 더이상 호스트의 IP의 요청을 추적하지 않음
6. 컨테이너에서 curl 요청 -> ttcp에서 컨테이너의 IP의 요청을 추적하지 않음
7. 테이블에 컨테이너 IP 추가
8. 컨테이너에서 curl 요청 -> tcp에 요청이 찍힘

```
                                                            |   2022/04/09 06:05:35 Src addr        Port   -> Dest addr       Port   RTT   
$ curl -L http://localhost:8090 -XPUT -d 10.178.0.5         |                  
$ curl google.com                                           |   2022/04/09 06:06:02 10.178.0.5      47090  -> 172.217.161.206 80     0     
$ curl -L http://localhost:8090 -XDELETE -d 10.178.0.5      |
$ curl google.com                                           |   nothing
$ sudo docker run --rm alpine/curl:3.14 google.com          |   nothing
$ curl -L http://localhost:8090 -XPUT -d 172.17.0.3         |
$ sudo docker run --rm alpine/curl:3.14 google.com          |   2022/04/09 06:10:24 172.17.0.3      40878  -> 172.217.25.174  80     0
```