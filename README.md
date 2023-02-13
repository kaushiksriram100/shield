# shield

V1:
shield returns if a url is blacklisted or not upon inspected a file. This is a in-memory approach. 

## Test Results

Unit Tests:
```
SRKAUSHI-M-8A2X:shield srkaushi$ go test -v server/*
2023/02/12 21:38:13 Info: Connection to ETCD Succeeded
=== RUN   TestLookUpMalwareDB
--- PASS: TestLookUpMalwareDB (0.00s)
=== RUN   TestLookupMalwareEtcD
--- PASS: TestLookupMalwareEtcD (0.04s)
PASS
ok  	command-line-arguments	0.129s
SRKAUSHI-M-8A2X:shield srkaushi$ 
```

Benchmark Tests:
```
SRKAUSHI-M-8A2X:server srkaushi$ go test -bench=. -count=4
2023/02/12 21:31:32 Info: Connection to ETCD Succeeded
goos: darwin
goarch: amd64
pkg: shield/server
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
BenchmarkShieldServer-8   	     147	   7212364 ns/op
BenchmarkShieldServer-8   	     170	   6924004 ns/op
BenchmarkShieldServer-8   	     195	   6557280 ns/op
BenchmarkShieldServer-8   	     182	   6033491 ns/op
PASS
ok  	shield/server	7.467s
SRKAUSHI-M-8A2X:server srkaushi$ 
```

## Packaging

Kind Cluster:

Pre-reqs:
1. Docker running in local MAC
2. `brew install kind`
3. helm installed
4. kubectl installed

To deploy the full stack:
`ansible-playbook setup-via-ansible.yml`

To Upgrade only the shield App:
`ansible-playbook setup-via-ansible.yml --tags [upgrade]`
