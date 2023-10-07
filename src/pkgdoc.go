/*
Package pipeline provides an abstraction around the sequential processing of elements where any step in the processing can be parallelized.

Pipelines are organized into stages, where each Stage accepts input messages from a channel and provides output to a channel,
potentially of a different type. A Stage can buffer its output for instances where the next stage has a variable execution
time that can sometimes fall behind and sometimes outstrip the current stage's execution time. A Stage can also be run
across multiple goroutines, allowing an I/O bound task to continue processing subsequent incoming messages while waiting
for responses from the I/O operation(s).

A Pipeline can send its output to a channel which is closed upon completion, or to a function that is invoked repeatedly
until the data source is exhausted. Pipelines can be canceled, and each Stage can abort the entire pipeline if an unrecoverable
error is encountered.
*/
package pipeline
