# Parallels Database

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
