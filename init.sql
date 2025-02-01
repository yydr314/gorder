create database if not exists gorder;
use gorder;

drop table if exists `o_stock`;

create table `o_stock` (
    id int unsigned auto_increment primary key,
    product_id varchar(255) not null,
    quantity int unsigned not null default 0,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci;