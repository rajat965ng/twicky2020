# Chapter 7: Parallel data processing and performance

For instance, before Java 7, processing a collection of data in parallel was extremely cumbersome. 
First, you needed to explicitly split the data structure containing your data into subparts. 
Second, you needed to assign each of these subparts to a different thread. 
Third, you needed to synchronize them opportunely to avoid unwanted race conditions, wait for the completion of all threads, and 
finally, combine the partial results.
    
    Java 7 introduced a framework called fork/join to perform these operations more consistently and in a less error-prone way.

## Parallel streams
    
A parallel stream is a stream that splits its elements into multiple chunks, processing each chunk with a different thread. 
Thus, you can automatically partition the workload of a given operation on all the cores of your multicore processor and keep all of them equally busy. 
    
    
Let’s suppose you need to write a method accepting a number n as argument and returning the sum of the numbers from one to n. 
    
A straightforward (perhaps naïve) approach is to generate an infinite stream of numbers, limiting it to the passed numbers, 
and then reduce the resulting stream with a BinaryOperator that sums two numbers, as follows:
    
    public long sequentialSum(long n) {
        return Stream.iterate(1L, i -> i + 1)
                     .limit(n)
                     .reduce(0L, Long::sum);
    }

This operation seems to be a good candidate to use parallelization, especially for large values of n. 
But where do you start? 
Do you synchronize on the result variable? 
How many threads do you use? Who does the generation of numbers? 
Who adds them up?
    
### Turning a sequential stream into a parallel one

Call the method parallel on the sequential stream:    
    
    public long parallelSum(long n) {
        return Stream.iterate(1L, i -> i + 1)
                     .limit(n)
                     .parallel()
                     .reduce(0L, Long::sum);
    }


Parallel streams internally use the default ForkJoinPool, which by default has as many threads as you have processors, 
as returned by Runtime.getRuntime().availableProcessors().
    
You can expect a significant performance improvement in its parallel version when running it on a multicore processor. 
    
    
### Measuring stream performance

Java Microbenchmark Harness (JMH): This is a tool- kit that helps to create, in a simple, annotation-based way, reliable microbenchmarks for Java programs 
and for any other language targeting the Java Virtual Machine (JVM). 
    
    
Add a couple of dependencies to your pom.xml file
    
    <dependency>
      <groupId>org.openjdk.jmh</groupId>
      <artifactId>jmh-core</artifactId>
      <version>1.17.4</version>
    </dependency>
    <dependency>
      <groupId>org.openjdk.jmh</groupId>
      <artifactId>jmh-generator-annprocess</artifactId>
      <version>1.17.4</version>
    </dependency>
    
    <build>
        <plugins>
        <plugin>
          <groupId>org.apache.maven.plugins</groupId>
          <artifactId>maven-shade-plugin</artifactId>
          <executions>
            <execution>
              <phase>package</phase>
              <goals><goal>shade</goal></goals>
              <configuration>
                <finalName>benchmarks</finalName>
                <transformers>
                  <transformer implementation="org.apache.maven.plugins.shade.
                                         resource.ManifestResourceTransformer">
                    <mainClass>org.openjdk.jmh.Main</mainClass>
                  </transformer>
                </transformers>
              </configuration>
            </execution>
          </executions>
        </plugin>
        </plugins>
    </build>  
    
    
    
    @BenchmarkMode(Mode.AverageTime)
    @OutputTimeUnit(TimeUnit.MILLISECONDS)
    @Fork(2, jvmArgs={"-Xms4G", "-Xmx4G"})
    public class ParallelStreamBenchmark {
        private static final long N= 10_000_000L;
        
        @Benchmark
        public long sequentialSum() {
          return Stream.iterate(1L, i -> i + 1).limit(N)
                     .reduce( 0L, Long::sum);
        }

        @TearDown(Level.Invocation)
        public void tearDown() {
            System.gc();
        }
    }
    
When you compile this class, the Maven plugin configured before generates a second JAR file named benchmarks.jar that you can run as follows:
    
    java -jar ./target/benchmarks.jar ParallelStreamBenchmark

1. Use an oversized heap to avoid any influence of the garbage collector as much as possible.
2. We tried to enforce the garbage collector to run after each iteration of our benchmark.


Executing it on a computer equipped with an Intel i7-4600U 2.1 GHz quad-core, it prints the following result:

    Benchmark                                Mode  Cnt    Score    Error  Units
    ParallelStreamBenchmark.sequentialSum    avgt   40  121.843 ±  3.062  ms/op


    @Benchmark
    public long iterativeSum() {
        long result = 0;
        for (long i = 1L; i <= N; i++) {
    result += i; }
        return result;
    }


    Benchmark                                Mode  Cnt    Score    Error  Units
    ParallelStreamBenchmark.iterativeSum     avgt   40    3.278 ±  0.192  ms/op

    Benchmark                                Mode  Cnt    Score    Error  Units
    ParallelStreamBenchmark.parallelSum      avgt   40  604.059 ± 55.288  ms/op
    
The parallel version of the summing method isn’t taking any advantage of our quad-core CPU and is around five times slower than the sequential one. 


Two issues are mixed together:
    1.a. iterate generates boxed objects, which have to be unboxed to numbers before they can be added.
    1.b. iterate is difficult to divide into independent chunks to execute in parallel.
    2.   iterate operation is hard to split into chunks that can be executed independently, because the input of one function application always depends on 
         the result of the previous application.

The whole list of numbers isn’t available at the beginning of the reduction process, making it impossible to efficiently partition the stream in chunks to be 
processed in parallel. By flagging the stream as parallel, you’re adding the overhead of allocating each sum operation on a different thread to the sequential processing.    
  


So how can you use your multicore processors and use the stream to perform a parallel sum in an effective way?    

   1.a LongStream.rangeClosed works on primitive long numbers directly so there’s no boxing and unboxing overhead.
   1.b LongStream.rangeClosed produces ranges of numbers, which can be easily split into independent chunks.
   
       @Benchmark
       public long rangedSum() {
           return LongStream.rangeClosed(1, N)
                            .reduce(0L, Long::sum);
       }
   
   Benchmark                                Mode  Cnt    Score    Error  Units
   ParallelStreamBenchmark.rangedSum        avgt   40    5.315 ±  0.285  ms/op

       @Benchmark
       public long parallelRangedSum() {
           return LongStream.rangeClosed(1, N)
                            .parallel()
                            .reduce(0L, Long::sum);
       }
       
   Benchmark                                  Mode  Cnt  Score    Error  Units
   ParallelStreamBenchmark.parallelRangedSum  avgt   40  2.677 ±  0.214  ms/op


Finally, we got a parallel reduction that’s faster than its sequential counterpart.

Note that this latest version is also around 20% faster than the original iterative one, demonstrating that, when used correctly.

Moving data between multiple cores is also more expensive than you might expect, so it’s important that work to be done in parallel on another core takes 
longer than the time required to transfer the data from one core to another.    


### Using parallel streams correctly

The main cause of errors generated by misuse of parallel streams is the use of algorithms that mutate some shared state.

    public long sideEffectSum(long n) {
        Accumulator accumulator = new Accumulator();
        LongStream.rangeClosed(1, n).forEach(accumulator::add);
        return accumulator.total;
    }
    
    public class Accumulator {
        public long total = 0;
        public void add(long value) { total += value; }
    }

Unfortunately, it’s irretrievably broken because it’s fundamentally sequential. You have a data race on every access of total. 
And if you try to fix that with synchronization, you’ll lose all your parallelism.

    public long sideEffectParallelSum(long n) {
        Accumulator accumulator = new Accumulator();
        LongStream.rangeClosed(1, n).parallel().forEach(accumulator::add);
        return accumulator.total;
    }


### Using parallel streams effectively

Its difficult to advice whether it makes sense to use a parallel stream in a certain situation:

1. Automatic boxing and unboxing operations can dramatically hurt performance. Java 8 includes primitive streams (IntStream, LongStream, and DoubleStream) 
   to avoid such operations, so use them when possible.
2. In particular, operations such as limit and findFirst that rely on the order of the elements are expensive in a parallel stream. 
   For example, findAny will perform better than findFirst because it isn’t constrained to operate in the encounter order. 
   You can always turn an ordered stream into an unordered stream by invoking the method unordered on it. 
3. For a small amount of data, choosing a parallel stream is almost never a winning decision.
4. Take into account how well the data structure underlying the stream decomposes. For instance, an ArrayList can be split much more efficiently than a LinkedList, 
   because the first can be evenly divided without traversing it, as it’s necessary to do with the second. The primitive streams created with the range factory method 
   can be decomposed quickly.
5. A SIZED stream can be divided into two equal parts, and then each part can be processed in parallel more effectively, but a filter operation can throw away an 
   unpredictable number of elements, making the size of the stream itself unknown.
6. Whether a terminal operation has a cheap or expensive merge step. If this is expensive, then the cost caused by the combination of the partial results generated
   by each substream can outweigh the performance benefits of a parallel stream. (for example, the combiner method in a Collector)
   
  
## The fork/join framework

Designed to recursively split a parallelizable task into smaller tasks and then combine the results of each subtask to produce the overall result.  

### Working with RecursiveTask

Create a subclass of RecursiveTask<R>, where R is the type of the result produced by the parallelized task or of RecursiveAction if the task returns no result.
    
    protected abstract R compute();

Responsibility of this method.
1.  Logic of splitting the task at hand into subtasks.
2.  The algorithm to produce the result of a single subtask when it’s no longer possible or convenient to further divide it.


    if (task is small enough or no longer divisible) {
        compute task sequentially
    } else {
        split task in two subtasks
        call this method recursively possibly further splitting each subtask
        wait for the completion of all subtasks
        combine the results of each subtask
    }
    
 
### Best practices for using the fork/join framework

1. Invoking the join method on a task blocks the caller until the result produced by that task is ready. It’s necessary to call it after the computation of both subtasks has been started.
2. The invoke method of a ForkJoinPool shouldn’t be used from within a RecursiveTask. Instead, you should always call the methods compute or fork directly; only 
   sequential code should use invoke to begin parallel computation.
3. Calling the fork method on a subtask is the way to schedule it on the Fork- JoinPool. Doing this allows you to reuse the same thread for one of the two subtasks and avoid the overhead 
   caused by the unnecessary allocation of a further task on the pool.


### Work stealing

This means that the tasks are more or less evenly divided on all the threads in the ForkJoinPool. Each of these threads holds a doubly linked queue of the tasks 
assigned to it, and as soon as it completes a task it pulls another one from the head of the queue and starts executing it.

One thread might complete all the tasks assigned to it much faster than the others, which means its queue will become empty while the other threads are still pretty busy.
In this case, instead of becoming idle, the thread randomly chooses a queue of a different thread and “steals” a task, taking it from the tail of the queue.

## Spliterator

Spliterators are used to traverse the elements of a source, but they’re also designed to do this in parallel.

    The Spliterator interface
    
    public interface Spliterator<T> {
                boolean tryAdvance(Consumer<? super T> action);
                Spliterator<T> trySplit();
                long estimateSize();
                int characteristics();
    }
    
    T is the type of the elements traversed by the Spliterator.
    

tryAdvance : used to sequentially consume the elements of the Spliterator one by one, returning true if there are still other elements to be traversed.
trySplit: used to partition off some of its elements to a second Spliterator (the one returned by the method), allowing the two to be processed in parallel.
estimateSize: provide an estimation of the number of the elements remaining to be traversed.     


### The splitting process            

First step, trySplit is invoked on the first Spliterator and generates a second one.
Step two, it’s called again on these two Spliterators, which results in a total of four. 
Step 3, The framework keeps invoking the method trySplit on a Spliterator until it returns null to signal that the data structure that it’s processing is no longer divisible.
Step 4, This recursive splitting process terminates when all Spliterators have returned null to a trySplit invocation.

<hr>