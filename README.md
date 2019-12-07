# Dune Analytics assignment 
## Basic Query Executor
In few words, Query Executor runs SELECT queries on top of the PostgreSQL database, specifically BLOCKS table.

The structure of the BLOCKS table is the following:

```
TABLE blocks (
     "time" timestamp with time zone NOT NULL,
     number numeric NOT NULL,
     hash bytea NOT NULL,
     parent_hash bytea NOT NULL,
     gas_limit numeric NOT NULL,
     gas_used numeric NOT NULL,
     miner bytea NOT NULL,
     difficulty numeric NOT NULL,
     total_difficulty numeric NOT NULL,
     nonce bytea NOT NULL,
     size numeric NOT NULL
)
```
## Project frontend

#### Simple query
<img src="./screenshots/simple_query.png" width="400"/>

#### Aggregate query
<img src="./screenshots/aggregate_query.png" width="400"/>

#### How To Run This Project


###### Starting the server

Server configuration is located in the `config.json` file in the root directory. In order to create the table BLOCKS and fill with initial data, change the `database.initialize` to `true`. Change the parameters and run following commands:

```make build```

```./bin/dune```

###### Starting the frontend

Frontend configuration is located in the `frontend/config.json`. Change the parameters and open the `frontend/index.html` in your browser.