# Lutral – Experminetal Na'vi dictionary engine.

I had this itch to implement a Na'vi dictionary parser to see if my idea is more 
performant than existing implementations.

The general idea is to build a massive parser tree and return a match for each time
it reaches a Result-node at the end of a boundary. A lot of the processing is done
building the tree to keep lookups fast.

As for the results, the performance is on par with Fwew for most cases, but Lutral
wins by a large margin on ridiculously affixed words.

The project in and of itself is not a dictionary, nor will it be.
