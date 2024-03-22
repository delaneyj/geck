# G.E.C.K

![](./assets/geck.png)

## Golang Entity Component Kit

Taking the relational archetype ideas from [flecs](https://www.flecs.dev) and port them to Go.

What's different from alternatives like [Arche](https://mlange-42.github.io/arche/)

1.  Pros
    1. Generates custom API based on config file, your API is bespoke to your application and making more like a DSL.
    1. Uses sparse sets with groups allowing perfect SOA data layout.
    1. Initial benchmarks show its faster to add and remove components and event query.
    1. Handles relationships with any number of components.
    1. Has seperate API generated for tags vs components.
    1. Has custom API for components with only a single field to optimize out function calls.
1.  Cons
    1. Code generation is required. But as soon as you generics I think you are tied to recompiling your code anyway.
    2. Have to write a config file. <sup>[1](#config)</sup>

It uses Go's generics however there are no good way to make clean API with the current interface constraints. For example having a slice of different `SparseSet[T]` has to be done with `any`, whereas with code gen we can be more explicit.

This is a work in progress, and the API is subject to change.

<a name="config">1</a> Actively working on a Web based real-time config generator built into the generator.
