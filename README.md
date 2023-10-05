# Summary

This repository contains an opinionated set of types and abstractions for working with parallel tasks communicating via
channels in the [Go](https://www.golang.org) programming language.

## Motivation

Many times throughout the week, I deal with processes that operate on large quantities of data. This data is usually
produced by an HTTP request to an external API, or requires multiple HTTP requests to be made at various stages of
processing. Meanwhile, other parts of the processing pipeline operate on local resources or entirely in memory with
much lower latency than the HTTP calls. Even with the most robust third-party APIs, latency on the requests can come
and go in bursts, with the majority of the calls completing in tens of milliseconds, but occasionally spiking up into
hundreds or -- gods forbid -- thousands of milliseconds.

In these cases, it's helpful to allow quickly-completing calls to move along the pipe asynchronously and not be held
back by a single slow-running request that is ahead of them in the processing queue. Dealing with the locking and 
synchronization of multithreading is a bit of a chore, however, even with Go's well-developed synchronization 
primitives, and race conditions or resource starvation can still bite even the most seasoned developer.

This library aims to help eliminate these problems in most cases by providing a consistent, proven pattern to develop
parallel pipeline processing.

## Workflow

For most of the parallel processing I do, the job starts with a data source that produces sequences of a single data
type. This can be a relational or document database, a file, a message bus, or a network socket/web endpoint. These
items might be requested in bulk, requested by page, be delivered in bulk, or arrive on a message bus or queue.

From there, our application needs to inspect each item and perform some number of operations on them. The operations can
be as simple as validating data sanity, or as complex as performing related queries from other data sources and
transforming the original item or combining it with the resulting data.

Eventually, each item will produce some number of actions, such as writing the item to a data store, invoking another
process with a decorated or annotated version of the data, or similar.

During processing, errors can occur. These errors might be large enough that they require the entire process to abort,
or they might affect only a small subset of the input data. Errors may only need to be logged, or they might require the
process to perform some number of retries, or send the bad data to another flow or external process.

## Concepts

This library deals with three core concepts; `Producers`, which receive data into the pipeline, `Processors`, which
take some action on a datum and then either discard it or pass it along to another step, and `Consumers`, which are the
final recipient of a datum, after which no other processing takes place.