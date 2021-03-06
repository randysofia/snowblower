# Snowblower

A lightweight high-performance Golang Snowplow collector and enricher. Besides the language choice, Snowblower differs from the official Snowplow implementations the following ways:

- Snowblower supports SNS/SQS as the intermediate data store between stages
- Snowblower uses a JSON serialization for CollectorPayloads instead of Thrift.

It’d be rather trivial to add both a Kinesis stream as a destination for the collector as well as to support Thrift, at which point it would be a complete drop-in replacement for the Snowplow Scala Kinesis Collector. However, for our needs, SQS provides a pretty compelling solution.

## Performance and Cost

In initial testing, the collector service requires between 10 and 20 times fewer front-end compute resources than the Scala-based Snowplow Kinesis collector, based on the observation that we scaled down from 24 c3.xlarge machines to 2 on our initial deployment. There are likely many reasons other than the langauge choice including:

- Snowblower only ships collected payloads that have data. It ignores the large number of empty data requests generated by Snowplow trackers.
- The Scala-based Kinesis collector is clearly marked as beta and likely not optimized.

On the other hand, the two c3.xlarge instances that replaced the Scala cluster handle a peak of over 350,000 requests per minute with an average latency at our load balancer of ~15ms and a CPU load of around 20%. We could scale back to one server, but we’ll likely experiment with smaller instances first.

## On using SQS instead of Kinesis

One advantage to using SNS/SQS instead of Kinesis is that SQS scales transparently without explicit provisioning instruction.

## History

Snowblower originated at [wunderlist/snowblower](https://github.com/wunderlist/snowblower). The commit history seems to indicate that it underwent a one month period of development before the project was entirely abandoned. Its abandonment seems to correlate with the period of time Wunderlist had been acquired by Microsoft.

Spark451 picked it up two years later as a replacement to an existing analytics stack. Much of what was promised in the Readme was never implemented by Wunderlist. It's now in working condition, but is still a work in progress.

## Future

- We'd like to see elements of the enrichment process moved to a javascript interpreter so that business logic could be incorporated without editing core.

- We'd like to make the storage configurable so that a user could choose one or more of many different options as destinations for events.

- Multiple instances of precipitate can not be run in parallel due to a race condition on the log file being processed. This should be moved to a queue based system like the enricher.

## Running

Snowblower has three commands:

- `collect` Runs the collector, sending events to SNS or SQS, acting as the second stage in a Snowplow pipeline.
- `etl` Pulls events from SQS, enriches them, and sends them into storage into MongoDB, acting as the third stage in a Snowplow pipeline.
- `precipitate` Pulls events from Cloudfront logs recorded on S3 and sends them to SNS for future enrichment see: [Setting up the Cloudfront collector](https://github.com/snowplow/snowplow/wiki/Setting-up-the-Cloudfront-collector).


## Configuration

The following environment variables configure the operation of Snowblower when running the collector:

- `SNS_TOPIC` Must contain the ARN of the SNS topic to send events to. **REQUIRED FOR `collect`**
- `SQS_URL` Must contain the URL of the SQS endpoint. **REQUIRED FOR `etl`**
- `S3_PATH` Must contain the URI of the S3 path (ie. s3://bucket/path). **REQUIRED FOR `precipitate`**
- `PORT` Optionally sets the port that the server listens to. Defaults to 8080.
- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_DEFAULT_REGION` Amazon Web Services credentials / config. If not set, Snowblower will attempt to use IAM Roles.
- `MONGO_URI` The mongo connection string for the DB.
- `MONGO_DB` The mongo DB to use.
- `MONGO_COLLECTION` The mongo collection to save events to.
- `COOKIE_DOMAIN` if not set, a domain won't be set on the session cookie
- `UA_REGEX` Location to Useragent regex
- `GEO_DB` Location to binary geo DB, see: http://dev.maxmind.com/geoip/geoip2/geolite2/

These settings can also be placed in a .env file. See the envfile as an example.

## Flags

- `--check` or `-c` Operates in checkmode (or dryrun mode) making no changes to S3, SNS, DB, SQS, etc..

## Installation

Quick install reference:

- Install godep see: [github.com/tools/godep](https://github.com/tools/godep)
- `godep restore` installs the package versions specified in Godeps/Godeps.json to your $GOPATH.
- `godep go install` compiles and places snowblower binary in bin dir.
