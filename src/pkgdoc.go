/*
Package pipeline provides abstractions for working with interrelated tasks across multiple concurrent processes without
the need to think about channels, goroutines, or contexts.

A pipeline has three key concepts: Sources, Processors, and Consumers. A Source produces sequences of items which are
passed into the pipeline. A Processor accepts items of one type, performs an operation on them, and returns an item,
potentially of a different type. A Consumer receives items and performs some operation on them, but does not produce
further items for the pipeline.

Concurrency is achieved internally by wiring the functions together through buffered channels and running each on its
own goroutine. This means that once a Stage has produced an item, it can init handling the next item while the next
stage of the pipeline works. Additionally, options are available to run any stage of the pipeline in parallel across
some arbitrary number of goroutines. This is useful if, for example, a stage of your pipeline is I/O bound and does
not process items as fast as the stage feeding it, which is especially common when working with web-based APIs or
other network resources.

Another common pitfall of working with concurrency is "zombies", or processes that keep running even though there is no
further work for them to do because another process they depend on has died. Any stage within a Pipeline can notify its
Pipeline that an unrecoverable failure has occurred and all active processing will be terminated. Likewise, every
pipeline has a Cancel method which will terminate the pipeline immediately.

Pipeline also have Source and Consumer options to help integration with existing channel-based code. A ChannelSource
is a Source which wraps an existing channel, feeding the items from the channel directly into the pipeline and signaling
completion when the channel is closed. A Sink is a Consumer which wraps an externally-created channel and passes
pipeline items to it. ChannelSink can be configured to close the channel when it is done, or to simply exit and mark
the pipeline as complete without closing the channel.

Pipelines may also have zero or more functions of type Filter applied between their stages. The job of a Filter is
to determine whether a given item should be processed or silently dropped from the pipeline.

Last but not least, pipelines may be split or joined to create more complex workflows.
TODO: Document those types and point to them from here.
*/
package pipeline
