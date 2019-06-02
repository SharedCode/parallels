# Parallels
A library for abstraction of software system's interactions.

Hola!
Among the many times we had authored and/or maintained enterprise applications including ingestion systems or data pipelines, there was and still, a big need to have a simple code library for abstracting software pieces' interactions.

A library that is so simple to use, abstracts most/all of the finicky details of software interactions such as messaging, protocols, method of execution such as high concurrency/parallel computing. This library will remove need to be bothered with messaging and puts ownership of management/execution of threads and/or applications (if on a cluster) onto the underlying infrastructure and/or the application. Example, Kubernetes can provide built-in load balancing and instance management as driven by demand. Also, within the application domain, we can implement there high concurrency techniques such as multi-threading, caching, asynchronous-ness, etc...

By using the library, we kind of standardizes the pattern of construction and thus, simplifies the overall system architecture.
There are many benefits in designing along this line such as:
- standards based, i.e. - using the library, resulting application code is consistent per interactions and method of execution.
- simplified *clustered* applications' authorship, i.e. - developer authors (& tests!) the different apps/modules and business functions within a *single* application (i.e. - using the single-app mode of the library), and enable the remote interactions (e.g. - HTTP calls) between these apps during deployment time. This accelerates development time as it removes time spent on (multi) application management in the cluster, enables code level debugging experience (step-in, out, breakpoints, variable inspection, etc...) as code is running in a single app.
- protocol portability, i.e. - apps using this library can switch to use different messaging/protocol if needed. E.g. - from HTTP protocol onto another newer technique, example, perhaps, a new custom protocol that uses Redis and HTTP for reduced latency, etc... The change is localized in the library (or via configuration switch!), and all apps using it, automatically adapts.
- common place for *hooking in* implementation of nifty features *interaction* wise, e.g. - automating application cluster management can be easily implemented & hooked-in in this library, like for example, combined with K8, an application can automatically be given a system *functional role* as it gets auto-instantiated by K8, it can read the *role* information from a source such as config or database, then behave according to that functional role definition. The library's utility for message sourcing, syncing and *function invocation* is a great place for hooking up this application feature.

## Background
Idea is to come up with a "cluster Virtualization" so that we can "program" in the cluster without getting bothered "writing code for the multi-machine hosting"-ness of our code, e.g. - interprocess communication/interaction such as messaging, queueing, polling, etc..., all sounds quite attractive. This allows us to "program" with the feeling that we are in a single "huge" computer. This will simplify and make our code very easily maintainable and portable (between 1 to n "machines" in the cluster and/or between clustered environments, e.g. - K8, AWS, Google Cloud, Azure, etc...).

## Programming Patterns
In order to keep this project simple, below are the identified distributed computing programming patterns that will be implemented. To begin with, there will be at least one pattern that will be implemented, and as time permits, new patterns will be added and thus, increasing the value of this code library.

### Controller/Worker Pattern
The Controller/Worker parallel computing pattern is one of the most typical scenarios in networked or distributed programming models. Where there is a "controller" role that generates jobs and there is a "worker" role that works or does the processing of the jobs created by the "controller".

The library will have this programming functionality internally defined. But using this library from your code, we will attempt to hide it. Thus, keeping the library simple to use and scale up to the cluster if needed and available (multi-threaded "workers" of one(the host) or more apps) or down to the app (multi-threaded "workers" of the host app) if sufficient to accomplish the task or if configured to do so.

Each application instance that uses this library can attain a specific role depending on which functionality is invoked. For example, when calling an execute or invoke function submitting job(s), automatically makes this app that does the call to be the *Controller*, and any instance(s) in the cluster that received (directly or indirectly) the job(s) and work or process them, attains the *Worker* role. This *invoke session* can be done in persisted manner, thus, allows app(s) of a cluster to be able to do app specific logic pertaining to the lifetime of the job(s), such as *retries & give up*, *resume where it left off*, etc...

### Scatter/Gather Pattern
This library will also support *scatter & gather* pattern. Using the API though, will minimize footprint or complexity of doing the distribution of work and gathering of results to/from many processors or workers local or *remote*. The *Scatter* function will distribute the jobs to both local and remote workers as needed. When running in a cluster env't such as K8, jobs processing will be load balanced across different nodes of the cluster.

### Supported Messaging/Protocol Methods
Initially, following will be supported:
- local interactions, i.e. - functions are called locally, no remoting.
- HTTP interactions, i.e. - calls are done via HTTP protocol. Message payload (de)serialization needs to be supplied by the app or use *default* method, e.g. - JSON marshaller offered by the *go* runtime, if using the *go* implementation of the library.
- Persisted interactions, i.e. - message payload persisted to a database, e.g. - Redis + (optional) Cassandra, and interaction calls done via HTTP. Built-in functions can be provided to implement nice to have features such as: retry on failure, (optional, if needed!) non-interrupted operations of source app even on times the target app is offline (message is persisted/queued during offline THEN processed when app resumes operation), etc...
- Others..., i.e. - we can always add new methods of messaging/protocol when/if needed.

## Implementation Plan
In order to keep this project simple and get focused on defining the (simple) code library API, using the following platform and/or software pieces sound reasonable. NOTE: code will be written in a way that it will be easy to "plugin" another "cluster" solution, e.g. - AWS Cloud Formation "direct", etc... Being able to (plugin!/) run across different cluster solution is a "nice to have" goal, design wise.

### * Kubernetes (K8)
The code library can be used in any application clustering environment, even in standalone app mode, as described above. K8 can be used though, to provide automatic instance scaling and (call) load balancing to these instances' endpoints.

### * Redis & Persistence Store
Redis will be used where appropriate in order to provide built-in *out of process* in-memory data caching that scales. If time permits, an option to hook in Cassandra for *persistence* will be supported. Example, Redis cluster can be restarted and persisted data read from Cassandra to restore previous Redis cluster *state*.
Data persistence to Cassandra will be done non-disruptedly not to impose unnecessary delays due to latencies in Cassandra I/O, which is minimal, though. If needed, *buffering & asynchronous* I/O will be implemented to minimize wait-times or latency when doing Cassandra backend DB I/O, when data is not in Redis cluster (reading) and when saving/syncing data to Cassandra (writing) during upserts to the Redis cache.

### * GO
Initial(i.e. - reference) implementation will be written in the "go" programming language. Go is the initial language of choice due to a number of reasons, including code written in "go" tend to be easily maintainable. Plus, K8 and other open source software bits such as docker, were written in "go" language as well. This gives us a "feeling" of using a "first" class citizen in the software stack.

## Code Library API
The "parallels" library will have the following essential parts:
### * Configuration
Configuration information is needed to allow the application to define how it wants to use the "parallels" library. Here is where the developer defines information such as:
- whether "parallels" will run "workers" as threads of the Host app, or they can run on other App(s) in the cluster.
- cluster run-time specification. Here is where developer defines things such as number of "worker" limits (or no-limits), etc...

### * Invoke Delegator
Invoke Delegator will take care of actually sending the message(s) to the target endpoint, via selected protocol, passing necessary parameter(s) and other message payload as needed.

### * Lifetime Mgmt
(Optional) Methods for controlling (e.g. - aborting or signalling aborts) lifetime of running or in-standby "workers" in the cluster. This is optional to keep usage of "parallels" easy. Default lifetime mechanism that works in "most/all" cases will be automatically provided, if not yet available.

## Use-Cases
Following are the different programming use-cases that will be supported by the library, or (minimally!) provide utility functions to aid implementation in the app. Ideally, these cases are best implemented in the application as it knows best how to deal with it, or the infrastructure/env't hosting the app. During implementation, we will see what utility functions are needed in the library. This feature works with *message persistence* option, as a feature provided by the library.

### * Process a list of Items
This is a case where we have "sourced" a list of Items (a.k.a. jobs) and which, we want to process them in parallel or concurrently, using potentially "multiple workers" to accomplish the task quickly.

### * Streaming: Periodic Batch Processing
This is a case in which, being able to (periodically) "source" job(s) from data source(s) then "sink" these job(s) for processing to "target" worker(s).

### * Stateful Processors
This is a case where the "processing" function can store state information(to the "cluster" memory, e.g. - redis, or to a disk, e.g. - persisted volume in K8) as it performs different processing stages, then during a crash and restart, the same (can be different!) process can detect "last" persisted state then resume processing where it left off.

- A. stateful from Controller Context<br/>
This means that the instance achieves an "identifiable" state as used from the controller context. For example, an instance is used to host a "MongoDB", then it is given a name that can be referenced from any of the other instance(s) in the cluster, essentially allowing them to interact with this "MongoDB" instance via its name.
- B. stateful from Instance Context<br/>
This means that the instance tracks some level of stateful information, that during restart, it can detect last persisted state and then, be able to resume operations where it left off.

An instance in the cluster can claim either or both of the above listed stateful-ness, depends on the application use-case.
See K8's "Stateful Set" feature for supporting letter A classified stateful-ness. Letter B statefulness can be achieved using Redis by storing instance state on Redis, and then reading such state during startup. And/Or via "persisted volume" in K8.

## Parallels Database

Parallels DB combines in-memory (e.g. - Redis) & NoSQL (e.g. - Cassandra) DBs to create a database that can *accomodate* the ULTRA high performance/throughput required by *huge multi-processing bandwidth* of parallels' workers.

Please feel free to read more details of this database here: https://github.com/SharedCode/parallels/blob/master/README.DB.md
