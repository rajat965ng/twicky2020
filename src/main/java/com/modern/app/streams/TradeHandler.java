package com.modern.app.streams;

import java.util.Arrays;
import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

public class TradeHandler {

    public static void main (String args[]) {
        Trader raoul = new Trader("Raoul", "Cambridge");
        Trader mario = new Trader("Mario","Milan");
        Trader alan = new Trader("Alan","Cambridge");
        Trader brian = new Trader("Brian","Cambridge");
        List<Transaction> transactions = Arrays.asList(
                new Transaction(brian, 2011, 300),
                new Transaction(raoul, 2012, 1000),
                new Transaction(raoul, 2011, 400),
                new Transaction(mario, 2012, 710),
                new Transaction(mario, 2012, 700),
                new Transaction(alan, 2012, 950)
        );

        List<Transaction> trx2011 = transactions.stream()
                .filter(trxn -> trxn.getYear()==2011)
                .sorted(Comparator.comparing(Transaction::getValue))
                .collect(Collectors.toList());
        System.out.println(trx2011);

        List<String> uniqueCities = transactions.stream()
                .map(Transaction::getTrader)
                .distinct()
                .map(Trader::getCity)
                .collect(Collectors.toList());
        System.out.println(uniqueCities);

        List<Trader> traders = transactions.stream()
                .filter(transaction -> transaction.getTrader().getCity().equals("Cambridge"))
                .map(Transaction::getTrader)
                .sorted(Comparator.comparing(Trader::getName))
                .distinct()
                .collect(Collectors.toList());
        System.out.println(traders);

        String tradersName = transactions.stream()
                .map(transaction -> transaction.getTrader().getName())
                .distinct()
                .sorted()
                .collect(Collectors.joining());
        System.out.println(tradersName);

        System.out.println(transactions.stream().map(Transaction::getValue).reduce(Integer::max));
        System.out.println(transactions.stream().reduce((t1, t2) -> Integer.compare(t1.getValue(),t2.getValue())<0?t1:t2));
        System.out.println(transactions.stream().min(Comparator.comparing(Transaction::getValue)));

        IntStream.rangeClosed(1,100).boxed()
                .flatMap(a -> IntStream.rangeClosed(1,100)
                        .filter(b->Math.sqrt(a*a + b*b)%1==0)
                        .mapToObj(b -> new int[]{a,b, (int) Math.sqrt(a*a + b*b)}))
                .forEach(ints -> System.out.println(Arrays.toString(ints)));

    }

}
