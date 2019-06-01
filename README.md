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

A new database virtualizer software module that combines the best of both worlds, of an in- memory and NoSQL DB to achieve yet another level of high Performance, Reliability & Availability (PRA).
Reference implementation uses Redis for in-memory data management and Cassandra for the backend database persistent store, and will easily plug-in with any other in-memory caching and NoSQL database solutions if required.
## Background
With the advent of in-memory database systems such as Redis, engineering teams get armed with an ultra fast data storage and management system. But since it is an in-memory solution, i.e. - data is not adapted for persistence to a backend storage system, data can get lost.
On the other hand, with the advent of NoSQL database systems such as Cassandra, MongoDB, etc..., engineering teams get armed with a tunable modern database engine that Scales in performance, reliability and availability. However, NoSQL DB still canʼt compete with the speed offered by in-memory solution, due to their ultra fast storage and management form factor, i.e. - the computerʼs (or a set of computers in a Cluster!) RAM.
Then, here is where *parallels DB* comes in. It combines the in-memory storage and operational characteristics with that of the NoSQL DBʼs, thus, achieving a new level of PRA!
## Sample Code
Here is a sample code equivalent of Hello World using this *parallels DB* code library.
```
    // this is from one of *parallels DB*'s unit tests code files.
    import "testing"
    import "parallels/database/common"
    import "os"
    func TestUpsert(t *testing.T) {
        dir,_ := os.Getwd()
        // 'config.json' file containing Redis & Cassandra cluster settings
        var config, _ = LoadConfiguration(dir + "/config.json")
        repo,_ := NewRepository(config)
        const Blob = 1
        rs := repo.Upsert([]common.KeyValue{*common.NewKeyValue(Blob, "123", []byte("Hello World!"))})
        if !rs.IsSuccessful() {
            t.Error(rs.Error)
        } 
    }
```
## API
The library provides two main use-cases:
* To provide a Repository interface that can be used to access/manage data
* To provide a Queryable or Navigable Repository<br/>
On top of the regular use-case of a Repository, this can be used to store/access data useful (for example), for storing metadata or summary data describing or having pointers/references to the data stored on the regular Repository. A typical use-case this addresses is, for example, storage of checkpoints & generated *data keys* of these checkpoints. Checkpoints typically captures a given point in time. And data *keys* generated in this point in time. These two can be stored/access from the Navigable Store, which, can then be used for cases, where the App needs to find and process a given set of data of a given point in time.

Following are the main API/methods available to use the library:
* Repository interface has *Upsert, Delete & Get* methods usable for data access and management
* NewRepository method instantiates a Repository that can be used in code to access/manage data
* NewRepositorySet<br/>
This method instantiates a RepositorySet that has two members: Repository and Navigable Repository. The latter, provides a Repository that has method to allow data query, or filter expression.

## Solution
The software solution is very easy, it virtualizes the two mentioned database storage systems. More like, how the Operating System does data paging, in this case, the data rows in in- memory of the Cluster (e.g. - Redis) gets swapped in/out and to/from the backend permanent data store (e.g. Cassandra). Most Recently Used (MRU) data sets get to be cached (persisted!) and the least used, stored in the backend DB, ready for getting cached when the application needs it. And when all of the data can fit in the in-memory cache, then they can all be cached and available to the Application. Reloading them back to in-memory cache is a fast and (easy/implicit) manner. Just use the software library and youʼre good to go.
* Redis In-memory Cache<br/>
Redis is the number one in-memory caching solution on the planet. Thus, it is industry proven when it comes to its performance, reliability, support, etc... including backing of the open source community.
* Cassandra NoSQL<br/>
Cassandra is the number one NoSQL engine that offers highly tunable Consistency, Availability and Performance (CAP) characteristics of the database engine. This is the reason it was chosen for initial implementation as *parallels DB* backend Storage engine. Via configuration, we can fine tune the CAP attributes to satisfy the Application requirements. Also, the DB schema was fine tuned for optimal Cassandra performance, e.g. - soft deletes, data partitioning & data batching models that will not generate hot spot (or bottleneck in the Cluster) during production use.

# Software Construction Patterns
The software library adapts the following software construction patterns to achieve the adaptability that it sports:
* Facade Pattern, e.g. Repository<br/>
The Repository interface defines the methods needing implementation for plugin to a database system such as Redis and Cassandra. In the future, supporting another caching solution, other than Redis, and another backend Storage, other than Cassandra, will be a jiffy, via Repository interface implementation.
* Test Driven Devʼt (TDD)<br/>
Unit tests were authored in all phases of the software life cycle, in order to achieve a good quality, from designs to implementations & testing.
* Module<br/>
The code library is written as a *go* module and thus, can be easily & safely reused in any application devʼt. Furthermore, adding an HTTP service on top of the module, will open it up for reuse to other programming languages, other than golang, via HTTP protocol interactions.

# Build & Deploy Info
* golang compiler and runtime version 1.12
* other dependencies such as golang client for Redis, Cassandra,... will automatically get downloaded during build process
* specify Redis & Cassandra connectivity and cluster info in the project's config.json file and you're good to go!
