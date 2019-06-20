create keyspace store with replication = {'class':'SimpleStrategy', 'replication_factor':1};

create table key_value(
    group text,
    key text,
    value blob,
    is_del boolean,
    updated TIMESTAMP,
    PRIMARY KEY ((group,key))
);

create table key_value_navigable(
    /* group is a data grouping construct, usually contains time series data set for optimal Partitioning in Cassandra.
     It can be time period like by the hour, by the day, etc... or custom grouping construct */
    group text,
    key text,
    value blob,
    is_del boolean,
    updated TIMESTAMP,
    /* partition key on type, cluster key on key.
     allows conditional access on Keys such as where "greater than", etc...
     Use with caution as a Partition can get large. This is typically
     used for storing navigable sets of (small!) metadata, e.g. - time series keyed
     data sets, etc...*/
    PRIMARY KEY (group,key)
);