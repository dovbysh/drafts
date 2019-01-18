create database http_log;

CREATE TABLE http_log.tables_delta
(
    tables_delta_id UInt64,
    query_type String,
    table_name String,
    column_data Nested
    (
        column_type String,
        column_name String,
        old_value Nullable(String),
        new_value Nullable(String)
    )
) ENGINE = MergeTree() order by tables_delta_id PRIMARY KEY (tables_delta_id);

insert into http_log.tables_delta values (1,'insert', 'user', ['int','string'],['user_id','name'],[null, null],['1','admin']);
insert into http_log.tables_delta values (2,'insert', 'user', ['int','string'],['user_id','name'],[null, null],['2','notAdmin']);

select * from http_log.tables_delta format JSON;
select * from http_log.tables_delta format JSONCompact;
select * from http_log.tables_delta format JSONEachRow;

