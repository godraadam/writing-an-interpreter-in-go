let forEach = fn([x, ...xs], cb) {
    if (x != nil) {
        cb(x)
        forEach(xs, cb) 
    }
} 

forEach([1, 2, 3], fn(x) { print x })