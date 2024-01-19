# G.E.C.K
## Golang Entity Component Kit
![](./assets/geck.png)

Taking the relational archetype ideas from [flecs](https://www.flecs.dev) and port them to Go.

One of the key implementation details is around using [roaring bitmaps](https://roaringbitmap.org/about/) to represent the set of entities that match a given archetype, this allows for fast set operations and iterations in a 64 bit space.


This is a work in progress, and the API is subject to change.