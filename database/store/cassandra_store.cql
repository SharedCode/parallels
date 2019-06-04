create keyspace store with replication = {'class':'SimpleStrategy', 'replication_factor':1};

create table key_value(
    type int,
    key text,
    value Blob,
    is_del boolean,
    updated TIMESTAMP,
    PRIMARY KEY ((type,key))
);

create table key_value_navigable(
    type int,
    key text,
    value Blob,
    is_del boolean,
    updated TIMESTAMP,
    /* partition key on type, cluster key on key.
     allows conditional access on Keys such as where "greater than", etc...
     Use with caution as a Partition can get large. This is typically
     used for storing navigable sets of (small!) metadata, e.g. - time series keyed
     data sets, etc...*/
    PRIMARY KEY (type,key)
);
