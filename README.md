# go-control-plane
This project builds a reliable control-plane service that safely manages shared resources.
The system guarantees correctness even when requests are retried, duplicated, or interrupted by failures.
It uses a desired-vs-current state model and reconciliation to converge to the correct final state.

# Why this exists
Real distributed systems face crashes, retries, and concurrent updates. This project demonstrates how to design for correctness rather than assuming everything works.
