# Source Code and Sample Deployment

## The Flameblock Repository

Application code is [available on GitHub](https://github.com/gargath/flameblock)

### Repository Structure

```
┃
┣╸assets                    // Static assets served by the visualizer
┃
┣╸cmd                       // Go code for the executables
┃  ┣╸collector
┃  ┗╸visualizer
┃
┣╸deploy                    // Dockerfiles for the three container images
┃  ┣╸docker
┃  ┃   ┣╸collector
┃  ┃   ┣╸docs
┃  ┃   ┗╸visualizer
┃  ┃
┃  ┗╸kubernetes             // Deployment files used for the sample deployment
┃
┣╸docs                      // The documentation source
┃
┣╸pkg                       // Go code for the application logic
┃  ┣╸api                      // Types for JSON marshalling
┃  ┣╸config                   // Shared configuration handling
┃  ┣╸collector
┃  ┗╸visualizer
┃
┗╸test                      // Sample webhook payload for testing
```

## Sample Deployment

Flameblock is deployed on GKE alongside this documentation:

```
$ kubectl get pods -n flameblock
NAME                                     READY   STATUS    RESTARTS   AGE
blocktown-56d685ffb6-8lsm7               1/1     Running   1          175m
flameblock-collector-b764b654d-t72cq     1/1     Running   0          3h26m
flameblock-docs-d5bb58f8d-85zv9          1/1     Running   0          30m
flameblock-visualizer-85bb67c947-htsjn   1/1     Running   0          3h26m
nsolid-console-7bc6985bbc-v9nts          1/1     Running   0          175m
redisoperator-67c6c79d46-77vqw           1/1     Running   0          31h
rfr-flameblock-redis-0                   1/1     Running   0          31h
rfr-flameblock-redis-1                   1/1     Running   0          31h
rfs-flameblock-redis-6b4c99bbc7-k2qpd    1/1     Running   0          23h
rfs-flameblock-redis-6b4c99bbc7-st4mz    1/1     Running   0          23h
rfs-flameblock-redis-6b4c99bbc7-zffzd    1/1     Running   0          31h
```

### Nsolid Console

The console is exposed via Kubernetes Ingress at [http://nsolid.g.lightweaver.info/](http://nsolid.g.lightweaver.info/)

It is protected with Basic Auth to protect the stored license key


### Blocktown

Blocktown is exposed via Kubernetes Ingress at [http://blocktown.g.lightweaver.info](http://blocktown.g.lightweaver.info)

Using the script from [https://github.com/snyk/sre-exercise-sample-app/blob/master/bin/blast](https://github.com/snyk/sre-exercise-sample-app/blob/master/bin/blast)
it can be load tested like this:
```
$ ./blast blocktown.g.lightweaver.info
```

### Flameblock Collector

The collector is not exposed externally. It is only available to the console inside the cluster. The webhook is configured as `http://flameblock-collector:8000/hook` in the nsolid console


### Flameblock Visualizer

The visualizer exposes two endpoints via Kubernetes Ingress:

* [`/index.html`](http://flameblock.g.lightweaver.info/index.html) for the static web page
and
* [`/flamedata`](http://flameblock.g.lightweaver.info/flamedata) for the dynamically generated data used to render the graph

The data can be retrieved directly using cURL:
```
$ curl -s flameblock.g.lightweaver.info/flamedata | jq .
{
  "name": "root",
  "value": 56100,
  "children": [ ... ]
}
```
