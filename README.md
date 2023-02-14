# Shield
Shield tells if a url is blacklisted or not.

## TL;DR - How it works
- Send Request via browser or POSTMAN (to check if a URL is present in malware DB):
`GET http://127.0.0.1:8080/urlinfo/1/google.com:8080/search-the-internet.html?v=thevalue`

- Response JSON (after looking up malware DB). `is_malware_infected` field captures the result.
`{"url":"google.com:8080/search?v=thevalue","is_malware_infected":false}`. false = Ok to proceed. true = STOP.
- Add/Remove new blacklisting URL
    ```
    //Notice port number = 8081 for PUT & DELETE
    PUT http://127.0.0.1:8081/urlinfo/1/google.com:8080/search?v=8
    DELETE http://127.0.0.1:8081/urlinfo/1/google.com:8080/search?v=8
    ```
- Service Ports: 8080 & 8081 (more details below)

## Setup Instructions
### Pre-requisites
- MACBOOK+Docker environment. 
- Docker engine is setup/running (https://docs.docker.com/desktop/install/mac-install/). verify with `docker ps` command.
- Install `kind` -> `brew install kind` & export to $PATH. May require sudo
- Install `ansible` -> `brew install ansible`
- Install `kubectl` depending on your MAC processor - https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/
- Install `helm` -> `brew install helm`

```
SRKAUSHI-M-8A2X:shield srkaushi$ kind --version
kind version 0.11.0
SRKAUSHI-M-8A2X:shield srkaushi$ which kind
/Users/srkaushi/go/bin/kind
SRKAUSHI-M-8A2X:shield srkaushi$
SRKAUSHI-M-8A2X:shield srkaushi$ which ansible
/Users/srkaushi/Library/Python/3.9/bin/ansible
SRKAUSHI-M-8A2X:shield srkaushi$ 
SRKAUSHI-M-8A2X:shield srkaushi$ which kubectl
/usr/local/bin/kubectl
SRKAUSHI-M-8A2X:shield srkaushi$ 
SRKAUSHI-M-8A2X:shield srkaushi$ which helm
/usr/local/bin/helm
SRKAUSHI-M-8A2X:shield srkaushi$ 
```
- Ensure all the below are exported to $PATH (export $PATH=$PATH:<new_paths>)

### Setup
```
git clone https://github.com/kaushiksriram100/shield.git
cd shield
docker system prune -a -f (optional - to ensure space is reclaimed)
docker volume prune (optional - to ensure vol space is reclaimed)
//setup the app
ansible-playbook setup-via-ansible.yml  //may take 3-4 mins. 
```
#### Verification and try it out!
1. Ensure pods are running for atleast 2 mins (just to ensure no crashloops). 
```
SRKAUSHI-M-8A2X:shield srkaushi$ kubectl get pods
NAME                                   READY   STATUS    RESTARTS   AGE
shield-etcd-0                          1/1     Running   0          2m44s
shieldapp-deployment-5b4ddbd5c-8s52h   1/1     Running   0          2m43s
shieldapp-deployment-5b4ddbd5c-xjtzj   1/1     Running   0          2m43s
SRKAUSHI-M-8A2X:shield srkaushi$
```
2. Wait for all pods to start (may take 1-2 mins). Once pods are running fine, then try this URL from a browser or POSTMAN
http://127.0.0.1:8080/urlinfo/1/google.com:8080/search?v=8

response: {"url":"google.com:8080/search?v=8","is_malware_infected":false}

3. To ensure the service is secure, admin operations like PUT new URL, DELETE URL are performed using another port (8081). 

###### Security
4. Generally in additiona to secure services using SSL, admin port has to be also secured via approriate firewalls within the hosting platform. Like Network ACLs, private subnets in AWS. Access to admin port is restricted for internal access.

##### Admin Operations
5. Add a new URL - PUT http://127.0.0.1:8081/urlinfo/1/google.com:8080/search?v=8  (Notice host-port = 8081) - Response 200
6. Check if URL is blacklisted - GET http://127.0.0.1:8080/urlinfo/1/google.com:8080/search?v=8 (Notice port = 8080) - Response JSON
7. Remove a URL - DELETE http://127.0.0.1:8081/urlinfo/1/google.com:8080/search?v=8  (Notice port = 8081) - Response Code = 200 (If the url exists, it will be deleted)
8. Updating & Deleting URLs can be requested in batches using a ordered queuing system. 

##### Scale Up
8. To scale up the services to handle more load beyond a single host. 
   a). Increase replicas in this file (Line 8) `k8s_deployment_specs/shield.yml`
   b). Run with tags - `ansible-playbook setup-via-ansible.yml --tags [upgrade]`
   c) `kubectl get pods` -> Should show 3 pods.
   d) These services are being a kubernetes native load balancer that will distribute requests to multiple pods (services). 
   e) In real world, Kubernetes workers (in EC2) will be spread across multiple AZs and is resilient for fault. Our kind cluster model is just a emulation of a real production deployment. This service can scale up based on requests. 

## Test Results
Unit Tests:
To Run Tests, create a etcd docker container (Do not use the above kubernetes pods)
```
//start a etcd pod
rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && docker rmi gcr.io/etcd-development/etcd:v3.5.0 || true && docker run -d -p 2379:2379 -p 2380:2380 --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data --name etcd-gcr-v3.5.0 gcr.io/etcd-development/etcd:v3.5.0 /usr/local/bin/etcd --name s1 --data-dir /etcd-data --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380 --initial-advertise-peer-urls http://0.0.0.0:2380 --initial-cluster s1=http://0.0.0.0:2380 --initial-cluster-token tkn --initial-cluster-state new --log-level info --logger zap --log-outputs stderr

SRKAUSHI-M-8A2X:shield srkaushi$ go test -v server/*
2023/02/12 21:38:13 Info: Connection to ETCD Succeeded
=== RUN   TestLookUpMalwareDB
--- PASS: TestLookUpMalwareDB (0.00s)
=== RUN   TestLookupMalwareEtcD
--- PASS: TestLookupMalwareEtcD (0.04s)
PASS
ok  	command-line-arguments	0.129s
SRKAUSHI-M-8A2X:shield srkaushi$ 

//After adding ETCD backend
SRKAUSHI-M-8A2X:server srkaushi$ go test -v *
2023/02/13 13:18:37 Info: Connection to ETCD Succeeded
=== RUN   TestLookUpMalwareDB
--- PASS: TestLookUpMalwareDB (0.00s)
=== RUN   TestLookupMalwareEtcD
--- PASS: TestLookupMalwareEtcD (0.02s)
=== RUN   TestPutMalwareUrlToEtcD
&{200 map[]  false <nil> map[] false}Successfully Added Key--- PASS: TestPutMalwareUrlToEtcD (0.01s)
=== RUN   TestDeleteMalwareUrlToEtcD
&{200 map[]  false <nil> map[] false}Successfully Deleted Key--- PASS: TestDeleteMalwareUrlToEtcD (0.01s)
PASS
ok  	command-line-arguments	0.106s
SRKAUSHI-M-8A2X:server srkaushi$ 
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

//After adding etcd Backend
SRKAUSHI-M-8A2X:server srkaushi$ go test -bench=.
2023/02/13 13:17:14 Info: Connection to ETCD Succeeded
&{200 map[]  false <nil> map[] false}Successfully Added Key&{200 map[]  false <nil> map[] false}Successfully Deleted Keygoos: darwin
goarch: amd64
pkg: shield/server
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
BenchmarkShieldServer-8   	     301	   4473146 ns/op
PASS
ok  	shield/server	1.787s
SRKAUSHI-M-8A2X:server srkaushi$ 
```

### Examples:
![GET](https://github.com/kaushiksriram100/shield/tree/main/examples/Postman_example.png?raw=true)

![PUT](https://github.com/kaushiksriram100/shield/tree/main/examples/Postman_PUT_example.png?raw=true)

![DELETE](https://github.com/kaushiksriram100/shield/tree/main/examples/Postman_Delete_example.png?raw=true)
