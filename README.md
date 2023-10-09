# Dummy Data Generator CLI
This CLI tool allows you to efficiently generate a large amount of dummy data in a database. It supports both PostgreSQL and MySQL and provides a flexible configuration file to specify which tables and columns to populate.

## Installation
To install the CLI tool, run the following command:

```bash
go install github.com/ponyo877/dummy_data_generator
```

## Features
- Generate a substantial amount of dummy data in a database.
- Supports both PostgreSQL and MySQL.
- Customize data generation through a configuration file.
- Track progress with a visual progress bar.

# Configuration
| Field       | Description                                                               |
|-------------|---------------------------------------------------------------------------|
| tablename   | Name of the table where the data will be generated.                        |
| recordcount | Total number of records to be generated.                                   |
| buffer      | Buffer size for generating records (useful for optimizing performance).    |
| columns     | List of columns with their respective configurations.                       |
| columns[].name        | Name of the column.                                                        |
| columns[].type        | Data type of the column (e.g., number, varchar, timestamp).                  |
| columns[].rule        | Generation rule for the column.                                             |
| columns[].rule.type        | Dummy rule type (e.g., unique, const, pattern) |
| columns[].rule.format        | [type: unique only] Dummy data format (e.g., UUID(varchar), ULID(varchar), NOW(timestamp)) |
| columns[].rule.value        | [type: const only] Dummy data const value |
| columns[].rule.min        | start of sequential value |
| columns[].rule.max        | [type: pattern only] end of sequential value |
| columns[].patterns[].value        | [type: pattern only] repeated value |
| columns[].patterns[].times        | [type: pattern only] value of how many times to repeat |

#### Example Rules:

- `type: unique`: Generates unique values. sequential number(default), current_timestamp(format: NOW), UUID and ULID is supported
- `type: const`: Assigns a constant value.
- `type: pattern`: Generates values based on specified patterns. If you specify [{value: A, times: 2}, {value: B, times: 1}], it will create repeated values like [A,A,B,A,A,B,...] and so on. And if you specify {Min: 1, Max: 5}, it will create repeated values like [1,2,3,4,5,1,2,3,...] and so on.

##### Example
```yaml
tablename: table1
recordcount: 1000000
buffer: 1000
columns:
- name: id
  type: number
  rule: 
    type: unique
    min: 0
- name: name
  type: varchar
  rule:
    type: const
    value: ponyo877
- name: color
  type: varchar
  rule:
    type: pattern
    patterns:
    - value: male
      times: 3
    - value: female
      times: 2
- name: code
  type: varchar
  rule:
    type: unique
    format: UUID
- name: created_at
  type: timestamp
  rule:
    type: unique
    format: NOW
```
## Sub Command
| Sub Command           | Description                                                                                                  |
|------------------|--------------------------------------------------------------------------------------------------------------|
| dummy_data_generator cnt     | show number of record |
| dummy_data_generator gen     | generate dummy data |

## Option
| Option           | Description                                                                                                  | Default Value |
|------------------|--------------------------------------------------------------------------------------------------------------|---------------|
| -c, --config      | configuration file for dummy data. You can provide multiple configuration files using wildcards <br>(e.g., `-c "cfg_*.yaml"`) or by comma-separating them (e.g., `-c cfg_1.yaml,cfg_2.yaml`). | `config.yaml` |
| -d, --database    | name of the database to use.                                                                   | `mydb`        |
| -u, --dbuser      | database user name.                                                                            | `root`        |
| -e, --engine      | database engine to use. Supports `postgres` and `mysql`.                                     | `postgres`    |
| -h, --host        | database server host or socket directory.                                                      | `127.0.0.1`   |
| -p, --password    | database password to use when connecting to the server.                                       | `password`    |
| -P, --port        | database server port.                                                                         | `5432`        |

## Usage Examples
- Example 1: Check current number of records. (MySQL)
```bash
$ dummy_data_generator cnt -e mysql -h 127.0.0.1 -u root -P 5432 -p password -c sample_1.yaml,sample_2.yaml
+--------+-------+
| TABLE  | COUNT |
+--------+-------+
| table1 |     0 |
| table2 |     0 |
+--------+-------+
```

- Example 2: Generate dummy data to target table designated config file. (PostgreSQL, all default value without config)
```bash
$ dummy_data_generator gen -c "sample_*.yaml"
table1: 534000 / 1000000 in progress  [=====================>-------------------]  53 %
table2: 533000 / 1000000 in progress  [=====================>-------------------]  53 %
table3:    10000 / 10000 done!       [=========================================] 
```

