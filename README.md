# Tilores CLI

## What is Tilores?

Tilores is a highly-scalable entity-resolution technology that was
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

You need to have a valid license or test agreement in order to install Tilores
into your own account. For a quick test, you can visit the
[public SaaS version of Tilores](https://app.tilores.io).

1. Install the CLI

```
go install github.com/tilotech/tilores-cli@latest
tilores-cli version
```

2. Initialize the project

```
mkdir foocustomer
cd foocustomer
tilores-cli init
```

1. Modify the schema files in the newly created `schema` folder.

2. If you want to test the API in your own AWS account, you can do so via

```
tilores-cli deploy --region <your-aws-region>
```

Please note, that this requires at least
[Terraform in version 1.x.x](https://www.terraform.io) to be installed.

6. Removing the API again from your AWS can be done by

```
tilores-cli destroy --region <your-aws-region>
```

More help to get started:

* [Tilores Documentation](https://docs.tilotech.io)
* [Tilores Website](https://tilores.io) for general information
* [Tilores Public SaaS](https://app.tilores.io)
