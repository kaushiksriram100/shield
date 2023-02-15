# Shield
Shield tells if a url is blacklisted or not.

## TL;DR - How it works
- Send Request via browser or POSTMAN (to check if a URL is present in malware DB):
`GET http://127.0.0.1:8080/urlinfo/1/google.com:8080/search-the-internet.html?v=thevalue`

- Response JSON (after looking up malware DB). `is_malware_infected` field captures the result.
`{"url":"google.com:8080/search?v=thevalue","is_malware_infected":false}`. false = Ok to proceed. true = malware infected, stop.
- Add/Remove new blacklisting URL (notice port number 8081 for PUT and DELETE)
    ```
    //Use port number = 8081 for PUT & DELETE (use 8080 for GET)
    PUT http://127.0.0.1:8081/urlinfo/1/google.com:8080/search?v=8
    DELETE http://127.0.0.1:8081/urlinfo/1/google.com:8080/search?v=8
    ```
- Service Ports: 8080 & 8081 (more details below)

## Setup Instructions
### Pre-requisites (one time setup - may need sudo and processor specific installations)
- Macbook+Docker environment
- Ensure Docker engine is setup/running (https://docs.docker.com/desktop/install/mac-install/). verify with `docker ps` command.
- Install `kind` -> `brew install kind` & export to $PATH. May require sudo
- Install `ansible` -> `brew install ansible`
- Install `kubectl` depending on your MAC processor - https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/
- Install `helm` -> `brew install helm`

For Ubuntu Machine (skipping baking this in Ansible playbook due to specific sudo requirement & leaving it to reviewers discretion to install)

```
//ansible
sudo apt install ansible

//Install kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.17.0/kind-linux-amd64
ls -lrt
echo $PATH
cp -Rp kind /usr/local/bin/
which kind

//HELM
wget https://get.helm.sh/helm-v3.9.3-linux-amd64.tar.gz
tar xvf helm-v3.9.3-linux-amd64.tar.gz
cp -Rp /home/srkaushi/linux-amd64/helm /usr/local/bin/

//Docker
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg lsb-release
sudo mkdir -m 0755 -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo   "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

//kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod 755 kubectl 
cp -Rp kubectl /usr/local/bin/
```

Ensure all pre-reqs are in PATH. 

```
SRKAUSHI-M-8A2X:shield srkaushi$ kind --version
kind version 0.11.0
SRKAUSHI-M-8A2X:shield srkaushi$ which kind
/Users/srkaushi/go/bin/kind
SRKAUSHI-M-8A2X:shield srkaushi$ which ansible
/Users/srkaushi/Library/Python/3.9/bin/ansible
SRKAUSHI-M-8A2X:shield srkaushi$ which kubectl
/usr/local/bin/kubectl
SRKAUSHI-M-8A2X:shield srkaushi$ which helm
/usr/local/bin/helm
```
- Ensure all the below are exported to $PATH
`export PATH=$PATH:<paths_above>`

### Setup the Service (using ansible)
```
git clone https://github.com/kaushiksriram100/shield.git
cd shield
docker system prune -a -f (to ensure space is reclaimed, else may get disk issues)
docker volume prune (optional - to ensure vol space is reclaimed)

//Setup the app (this builds and deploys the app and dependencies)
ansible-playbook setup-via-ansible.yml  //may take 3-4 mins. 

PLAY RECAP *************************************************************************************************************************************************************************************************
localhost                  : ok=7    changed=6    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   


//Note: Docker image already published - docker pull registry.hub.docker.com/kaushik100/shield:1.3.0
```
#### Verification and try it out!
1. Ensure all pods are "running" for atleast 2 mins (just to ensure no crashloops). It may take 2-3 mins. 
```
SRKAUSHI-M-8A2X:shield srkaushi$ kubectl get pods
NAME                                   READY   STATUS    RESTARTS   AGE
shield-etcd-0                          1/1     Running   0          8m19s
shield-etcd-1                          1/1     Running   0          8m19s
shield-etcd-2                          1/1     Running   0          8m19s
shieldapp-deployment-5b4ddbd5c-9brhg   1/1     Running   0          8m18s
shieldapp-deployment-5b4ddbd5c-c28kg   1/1     Running   0          8m18s
shieldapp-deployment-5b4ddbd5c-zkt7b   1/1     Running   0          8m18s
```
2. Once pods are running fine, then try this URL from a browser or POSTMAN
http://127.0.0.1:8080/urlinfo/1/google.com:8080/search?v=8

response: {"url":"google.com:8080/search?v=8","is_malware_infected":false}

3. To ensure the service is secure, admin operations like PUT new URL, DELETE URL are performed using another port - 8081. 

###### Security
4. Access to "admin port - 8081" is restricted for internal use-cases (updating/deleting URLs). This can be done via appropriate firewalls, network acls, private subnets in a cloud/on-prem environment. In addition, services are secured using TLS.

##### Admin Operations
5. Add a new URL - PUT http://127.0.0.1:8081/urlinfo/1/google.com:8080/search?v=8  (Notice host-port = 8081) - Response 200
6. Check if URL is blacklisted - GET http://127.0.0.1:8080/urlinfo/1/google.com:8080/search?v=8 (Notice port = 8080) - Response JSON
7. Remove a URL - DELETE http://127.0.0.1:8081/urlinfo/1/google.com:8080/search?v=8  (Notice port = 8081) - Response Code = 200 (If the url exists, it will be deleted)
8. Updating & Deleting URLs can be requested in batches using a ordered queuing system. 

##### Scale Up To Handle more load
8. To scale up the services to handle more load beyond a single host. 
   a). Increase shield service replicas in this file (Line 8) `k8s_deployment_specs/shield.yml`. Set to 5 for example. 
   b). Run ansible with tags - `ansible-playbook setup-via-ansible.yml --tags [upgrade]`
   c) `kubectl get pods` -> Should show 3 pods.
   d) Kubernetes native load balancer will distribute requests to multiple pods (services). 
   e) In real world, Kubernetes workers (in EC2) will be spread across multiple AZs and is resilient for fault. Our kind cluster model is just a emulation of a real production deployment. This service can scale up based on requests and made 100% available.
   f) ETCD provides a highly consistent k,v database that can handle several transactions atomically at scale and compact keys. ETCD in our case has 3 replicas & can be horizontally scaled up in real environments.

## Test Results
Unit Tests:
To Run Tests, create a etcd docker container (Do not use the above kubernetes pods)
```
//start a etcd pod
rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && docker rmi gcr.io/etcd-development/etcd:v3.5.0 || true && docker run -d -p 2379:2379 -p 2380:2380 --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data --name etcd-gcr-v3.5.0 gcr.io/etcd-development/etcd:v3.5.0 /usr/local/bin/etcd --name s1 --data-dir /etcd-data --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380 --initial-advertise-peer-urls http://0.0.0.0:2380 --initial-cluster s1=http://0.0.0.0:2380 --initial-cluster-token tkn --initial-cluster-state new --log-level info --logger zap --log-outputs stderr

SRKAUSHI-M-8A2X:shield srkaushi$ cd server;go test -v
2023/02/12 21:38:13 Info: Connection to ETCD Succeeded
=== RUN   TestLookUpMalwareDB
--- PASS: TestLookUpMalwareDB (0.00s)
=== RUN   TestLookupMalwareEtcD
--- PASS: TestLookupMalwareEtcD (0.04s)
PASS
ok  	command-line-arguments	0.129s
SRKAUSHI-M-8A2X:shield srkaushi$ 

//After adding ETCD backend
SRKAUSHI-M-8A2X:server srkaushi$ go test -v 
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

Added a concurrent GET client to send approx 1000 requests in parallel. Refer to examples/ folder & readme.md in there. 

### Some Known Issues:
Testing concurrent GET/PUT URLs(about ~1000) to the service as-is in a MAC+Docker results in some "connection_reset" issues. Upon debugging, this attributes to specific MAC/Docker internal networking bottlenecks. 

Alternatively, deployed this service in ubuntu & tested around ~1100 concurrent requests and seems to be fine (refer examples/README.md)

