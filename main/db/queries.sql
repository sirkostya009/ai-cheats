/* This document contains multiple SQL queries for
   analyzing, manipulating and monitoring the database. */

select *, active_till > now() as active
from customers;

select * from requests
order by created_at desc;

-- add new customer
insert into customers (telegram, active_till)
values ('@sirkostya009', now() + interval '1 month');

-- reset customer hashes
update customers
set hashes = '{}'
where id = 1;

-- move customer to a newer model
update customers
set model = 'gpt-4-1106-preview'
where id = 1;

-- request analytics per model, uncomment a where clause to filter only recent requests
-- input3_5_price: 0.001
-- input4_price: 0.01
-- output3_5_price: 0.002
-- output4_price: 0.03
select *,
       prompt_tokens / total_requests                                                        as prompt_tokens_per_request,
       completion_tokens / total_requests                                                    as completion_tokens_per_request,
       prompt_tokens * (input_price)                                                         as total_input_price,
       completion_tokens * (output_price)                                                    as total_output_price,
       prompt_tokens * (input_price) + completion_tokens * (output_price)                    as total_price,
       (prompt_tokens * (input_price) + completion_tokens * (output_price)) / total_requests as price_per_request
from (select model,
             avg(finished_at - created_at) as average_completion_time,
             sum(prompt_tokens)            as prompt_tokens,
             sum(completion_tokens)        as completion_tokens,
             count(*)                      as total_requests,
             case
                 when model like '%3.5%' then :input3_5_price
                 when model like '%4%' then :input4_price
                 end / 1000                as input_price,
             case
                 when model like '%3.5%' then :output3_5_price
                 when model like '%4%' then :output4_price
                 end / 1000                as output_price
      from requests
      where status = 200
--         and created_at >= now() - interval '1 day'
      group by model) statistics;

-- average request interval time
select avg(time_diff), count(time_diff)
from (select extract(epoch from created_at - lag(created_at) over (order by created_at)) as time_diff
      from requests
      where customer_id = 1
        and status = 200
        and created_at > now() - interval '1 hour') time_diffs;

-- request analytics per model and customer, comment out the where clause to include all requests
select *,
       prompt_tokens / total_requests                                                        as prompt_tokens_per_request,
       completion_tokens / total_requests                                                    as completion_tokens_per_request,
       prompt_tokens * (input_price)                                                         as total_input_price,
       completion_tokens * (output_price)                                                    as total_output_price,
       prompt_tokens * (input_price) + completion_tokens * (output_price)                    as total_price,
       (prompt_tokens * (input_price) + completion_tokens * (output_price)) / total_requests as price_per_request
from (select customer_id,
             model,
             avg(finished_at - created_at) as average_completion_time,
             sum(prompt_tokens)            as prompt_tokens,
             sum(completion_tokens)        as completion_tokens,
             count(*)                      as total_requests,
             case
                 when model like '%3.5%' then :input3_5_price
                 when model like '%4%' then :input4_price
                 end / 1000                as input_price,
             case
                 when model like '%3.5%' then :output3_5_price
                 when model like '%4%' then :output4_price
                 end / 1000                as output_price
      from requests
      where status = 200
        and created_at >= now() - interval '1 day'
      group by customer_id, model) statistics;
