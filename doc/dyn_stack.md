I am designing statically typed, systems programming language
and compiler for it. I would like to have manual memory managment
in my language, maybe with some automatic checks and cleanup, but
nothing heavy-handed like in Rust for sure. I like the idea of using
different allocators for different situations, to have "right
tool for the job" so to say.

Right now I am contemplating the idea of "dynamic side stack"
allocation (dont know if that idea was already named, you can
point me to actual name, but it is not that important). The idea
rests on several concepts:

1. Compiler passes implicit context variable as first hidden argument
to each function call. Context carries common things like logger,
various allocators, and maybe some other things. For efficiency
compiler can inspect if function actually needs the context and
do not pass it if there is no need. It cretes some work in compiler
to "color" different functions, but nothing unmanageable in my book.

2. "Dynamic side stack" works on stack principle and memory allocated
through it lives only for the duration of the call. It kinda works like
classic frame-based CPU managed stack, but with more flexibility.
It can be used to allocate things with dynamic size and its capacity
can grow dynamically as needed. Implementation places it on the heap
so it can have much bigger capacity than regular stack.

3. To work properly compiler needs to place hidden cleanup calls
(similar to defer, but somewhat automatic). Consider example below:

```
// in this snippet $ means context variable
//
// dynamic side stack works similar to arena allocator,
// but with automatic checkpoints managed but compiler
fun dynstack_example() {
    // compiler knows that we use dynamic stack allocation
    // in this function body, so it places hidden
    // code to save dynamic stack state in a checkpoint
    // GENERATED:
    dynstack_checkpoint := $.dynstack.save();

    size := 0;
    // ... calculate needed size

    // memory allocated is only valid until we
    // leave the function
    data := $.dynstack.alloc(size);

    // allocated data is used by another function
    save_buffer(data, size);

    // when control flow leaves the function compiler
    // places code to restore dynamic stack checkpoint
    // GENERATED:
    $.dynstack.restore(dynstack_checkpoint);
}
```

Unlike traditional arena allocator we can grow dynamic stack
capacity as we allocate more and more by managing (for example)
a linked list of already allocated backing memory regions.

What do think of this design? Provide critique and good points
about it. Suggest impovements and alternatives. Point to similar
approaches and how they work.
