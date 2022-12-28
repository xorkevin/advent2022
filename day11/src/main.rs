macro_rules! monk {
    ([$($items:expr),* $(,)?]; $var:ident, $op:expr; $test:literal, $jt:literal, $jf:literal) => {
        Monkey {
            items: vec![$($items),*],
            op: |$var: i32| -> i32 { $op },
            test: $test,
            jt: $jt,
            jf: $jf,
        }
    };
    (
        Monkey $_idx:literal:
            Starting items: $($items:expr),*;
            Operation $var:ident: new = $op:expr;
            Test: divisible by $test:literal
                If true: throw to monkey $jt:literal
                If false: throw to monkey $jf:literal
    ) => {
        Monkey {
            items: vec![$($items),*],
            op: |$var: i32| -> i32 { $op },
            test: $test,
            jt: $jt,
            jf: $jf,
        }
    };
}

fn get_monkeys() -> Vec<Monkey> {
    vec![
        monk!(
        Monkey 0:
          Starting items: 84, 66, 62, 69, 88, 91, 91;
          Operation old: new = old * 11;
          Test: divisible by 2
            If true: throw to monkey 4
            If false: throw to monkey 7
        ),
        monk!(
        Monkey 1:
          Starting items: 98, 50, 76, 99;
          Operation old: new = old * old;
          Test: divisible by 7
            If true: throw to monkey 3
            If false: throw to monkey 6
        ),
        monk!(
        Monkey 2:
          Starting items: 72, 56, 94;
          Operation old: new = old + 1;
          Test: divisible by 13
            If true: throw to monkey 4
            If false: throw to monkey 0
        ),
        monk!(
        Monkey 3:
          Starting items: 55, 88, 90, 77, 60, 67;
          Operation old: new = old + 2;
          Test: divisible by 3
            If true: throw to monkey 6
            If false: throw to monkey 5
        ),
        monk!(
        Monkey 4:
          Starting items: 69, 72, 63, 60, 72, 52, 63, 78;
          Operation old: new = old * 13;
          Test: divisible by 19
            If true: throw to monkey 1
            If false: throw to monkey 7
        ),
        monk!(
        Monkey 5:
          Starting items: 89, 73;
          Operation old: new = old + 5;
          Test: divisible by 17
            If true: throw to monkey 2
            If false: throw to monkey 0
        ),
        monk!(
        Monkey 6:
          Starting items: 78, 68, 98, 88, 66;
          Operation old: new = old + 6;
          Test: divisible by 11
            If true: throw to monkey 2
            If false: throw to monkey 5
        ),
        monk!(
        Monkey 7:
          Starting items: 70;
          Operation old: new = old + 7;
          Test: divisible by 5
            If true: throw to monkey 1
            If false: throw to monkey 3
        ),
    ]
}

fn main() {
    {
        let mut monkeys = get_monkeys();
        let mut counts = vec![0; monkeys.len()];
        for _ in 0..20 {
            for n in 0..monkeys.len() {
                let processed = {
                    let i = &monkeys[n];
                    i.items.iter().map(|&j| i.process1(j)).collect::<Vec<_>>()
                };
                for (k, t) in processed {
                    monkeys[t].add(k);
                }
                let i = &mut monkeys[n];
                counts[n] += i.items.len();
                i.discard();
            }
        }
        counts.sort_by(|a, b| b.cmp(a));
        if let [first, second, ..] = counts[..] {
            println!("Part 1: {}", first * second);
        }
    }
    {
        let mut monkeys = get_monkeys();
        let modulus = monkeys.iter().map(|i| i.test).product::<i32>();
        let mut counts = vec![0; monkeys.len()];
        for _ in 0..10000 {
            for n in 0..monkeys.len() {
                let processed = {
                    let i = &monkeys[n];
                    i.items.iter().map(|&j| i.process2(j)).collect::<Vec<_>>()
                };
                for (k, t) in processed {
                    monkeys[t].add(k % modulus);
                }
                let i = &mut monkeys[n];
                counts[n] += i.items.len();
                i.discard();
            }
        }
        counts.sort_by(|a, b| b.cmp(a));
        if let [first, second, ..] = counts[..] {
            println!("Part 2: {}", first * second);
        }
    }
}

struct Monkey {
    items: Vec<i32>,
    op: fn(i32) -> i32,
    test: i32,
    jt: usize,
    jf: usize,
}

impl Monkey {
    fn discard(&mut self) {
        self.items.clear()
    }

    fn add(&mut self, val: i32) {
        self.items.push(val)
    }

    fn process1(&self, val: i32) -> (i32, usize) {
        let f = self.op;
        let k = f(val) / 3;
        if k % self.test == 0 {
            (k, self.jt)
        } else {
            (k, self.jf)
        }
    }

    fn process2(&self, val: i32) -> (i32, usize) {
        let f = self.op;
        let k = f(val);
        if k % self.test == 0 {
            (k, self.jt)
        } else {
            (k, self.jf)
        }
    }
}
