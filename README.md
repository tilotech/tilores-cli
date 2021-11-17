# TiloRes CLI

## What is TiloRes?

TiloRes is a highly-scalable, “entity-resolution” technology that was
originally developed to connect internal data together. The technology was
developed because we found that no other technology on the market could
deliver the speed, scalability or cost performance we demanded.

### What is Entity Resolution?

Entity resolution (ER) is the connecting of non-identical, related data from
disparate sources to “entities”. Entities can be anything from people, to
companies, to financial transactions.

### Why is this important?

Companies today are collecting more and more data, of varying quality and
from different/disparate sources, but they are only able to productively use
a fraction of this data. Why? Because matching this data together so that one
has the full data picture, is technically very difficult, especially at scale
and when data must be accessed in real-time. 

In order to fully utilise and get value from data resources - both
internal and externally-sourced - the data needs to be matched together -
entity resolution - in a manner which can be searched quickly.

## Quick Start

1. Install the CLI

```
go install gitlab.com/tilotech/tilores-cli@latest
tilores-cli version
```

2. Initialize the project

```
mkdir foocustomer
cd foocustomer
tilores-cli init
```

3. Start the local API webserver with a fake implementation

```
tilores-cli run
```

4. Modify the schema files in the newly created `schema` folder

More help to get started:

* [FooCustomer-Example](https://gitlab.com/tilotech/tilores-foocustomer) - an example on how a possible customer would start
* [Tilo Tech Website](https://tilotech.io) for general information
