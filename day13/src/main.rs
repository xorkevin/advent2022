use std::cmp::Ordering;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::iter::Peekable;
use std::vec::IntoIter;

const PUZZLEINPUT: &str = "input.txt";

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let div1 = Signal::List(vec![Signal::Num(2)]);
    let div2 = Signal::List(vec![Signal::Num(6)]);
    let mut signals = vec![
        Signal::List(vec![Signal::Num(2)]),
        Signal::List(vec![Signal::Num(6)]),
    ];

    let mut left = None;

    let mut count = 0;
    let mut num_pairs = 0;
    for line in reader.lines() {
        let line = line?;
        if line.len() == 0 {
            left = None;
            continue;
        }
        let tokens = tokenize(line.into_bytes().into_iter().peekable())?;
        let (sig, _) = parse_tokens(tokens.into_iter().peekable())?;
        match left {
            None => left = Some(sig),
            Some(l) => {
                num_pairs += 1;
                if let Some(v) = compare_sigs(&l, &sig) {
                    if v {
                        count += num_pairs;
                    }
                }
                left = None;
                signals.push(l);
                signals.push(sig);
            }
        }
    }

    println!("Part 1: {}", count);

    signals.sort_unstable_by(cmp_sigs);

    let div1_idx = signals
        .binary_search_by(|i| cmp_sigs(i, &div1))
        .or_else::<BoxError, _>(|_| Err("Divider missing".into()))?;

    let div2_idx = signals
        .binary_search_by(|i| cmp_sigs(i, &div2))
        .or_else::<BoxError, _>(|_| Err("Divider missing".into()))?;

    println!("Part 2: {}", (div1_idx + 1) * (div2_idx + 1));

    Ok(())
}

fn cmp_sigs(left: &Signal, right: &Signal) -> Ordering {
    if eq_sigs(left, right) {
        Ordering::Equal
    } else {
        match compare_sigs(left, right) {
            Some(v) => {
                if v {
                    Ordering::Less
                } else {
                    Ordering::Greater
                }
            }
            None => Ordering::Equal,
        }
    }
}

fn eq_sigs(left: &Signal, right: &Signal) -> bool {
    match (left, right) {
        (Signal::Num(a), Signal::Num(b)) => a == b,
        (Signal::List(a), Signal::List(b)) => {
            if a.len() != b.len() {
                return false;
            }
            for (l, r) in a.iter().zip(b.iter()) {
                if !eq_sigs(l, r) {
                    return false;
                }
            }
            true
        }
        _ => false,
    }
}

fn compare_sigs(left: &Signal, right: &Signal) -> Option<bool> {
    match (left, right) {
        (Signal::Num(a), Signal::Num(b)) => {
            if a == b {
                None
            } else {
                Some(a < b)
            }
        }
        (Signal::List(a), Signal::List(b)) => {
            for (l, r) in a.iter().zip(b.iter()) {
                if let Some(v) = compare_sigs(l, r) {
                    return Some(v);
                }
            }
            if a.len() == b.len() {
                None
            } else {
                Some(a.len() < b.len())
            }
        }
        (Signal::Num(a), b @ Signal::List(_)) => {
            compare_sigs(&Signal::List(vec![Signal::Num(*a)]), b)
        }
        (a @ Signal::List(_), Signal::Num(b)) => {
            compare_sigs(a, &Signal::List(vec![Signal::Num(*b)]))
        }
    }
}

#[derive(PartialEq, Eq)]
enum Token {
    Lparen,
    Rparen,
    Num(i32),
}

fn tokenize(mut b: Peekable<IntoIter<u8>>) -> BoxResult<Vec<Token>> {
    let mut tokens = Vec::new();
    while let Some(c) = b.next() {
        match c {
            b'[' => tokens.push(Token::Lparen),
            b']' => tokens.push(Token::Rparen),
            b'0'..=b'9' => {
                let mut buf = vec![c];
                while let Some(c2) = b.next_if(|&i| i >= b'0' && i <= b'9') {
                    buf.push(c2);
                }
                tokens.push(Token::Num(String::from_utf8(buf)?.parse::<i32>()?));
            }
            _ => (),
        }
    }
    Ok(tokens)
}

#[derive(Debug)]
enum Signal {
    Num(i32),
    List(Vec<Signal>),
}

fn parse_tokens(
    mut tokens: Peekable<IntoIter<Token>>,
) -> BoxResult<(Signal, Peekable<IntoIter<Token>>)> {
    let head = match tokens.next() {
        Some(v) => v,
        None => return Err("No tokens".into()),
    };
    match head {
        Token::Num(val) => Ok((Signal::Num(val), tokens)),
        Token::Lparen => {
            let mut signals = Vec::new();
            loop {
                match tokens.next_if(|t| t == &Token::Rparen) {
                    Some(_) => break,
                    None => (),
                };
                let (sig, rest) = parse_tokens(tokens)?;
                signals.push(sig);
                tokens = rest;
            }
            Ok((Signal::List(signals), tokens))
        }
        Token::Rparen => Err("Unexpected token".into()),
    }
}
