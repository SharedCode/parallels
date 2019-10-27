create keyspace store with replication = {'class':'SimpleStrategy', 'replication_factor':1};

create table key_value(
    /* group is domain specific data grouping construct. sample of groups are: string value that represents entity types like
    Album, TVProgram, etc...*/
    group text,
    key text,
    value blob,
    is_del boolean,
    updated TIMESTAMP,
    PRIMARY KEY ((group,key))
);

create table key_value_navigable(
    /* group in this case, is a data grouping construct, usually contains time series data set for optimal Partitioning in Cassandra.
        It can be time period like by the hour, by the day, etc... or custom grouping construct */
    group text,
    key text,
    value blob,
    is_del boolean,
    updated TIMESTAMP,
    /* partition key on group, cluster key on key. Allows conditional access on key such as using CQL "where clause"*/
    PRIMARY KEY (group,key)
);
