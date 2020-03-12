package com.modern.app;


import java.util.stream.IntStream;

public class ArrayDemo {

    public static void main(String args[]){

        IntStream.rangeClosed(1,10).map(operand ->  (int) ((Math.random()*100)+operand))
                .filter(value -> value>50 && value<90)
                .forEach(System.out::println);

    }

}
