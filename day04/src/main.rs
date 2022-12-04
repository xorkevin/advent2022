use lazy_static::lazy_static;
use regex::Regex;
use std::cmp::Ordering;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BoxResult<T> = Result<T, Box<dyn std::error::Error>>;

fn main() -> BoxResult<()> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let (count1, count2) = {
        let mut count1 = 0;
        let mut count2 = 0;
        for line in reader.lines() {
            let (a, b) = parse_line(&line?)?;
            if is_fully_contained(&a, &b) {
                count1 += 1;
            } else if is_overlap(&a, &b) {
                count2 += 1;
            }
        }
        (count1, count2)
    };
    println!("Part 1: {}\nPart 2: {}", count1, count1 + count2);
    Ok(())
}

fn parse_line(line: &str) -> BoxResult<(Pair, Pair)> {
    lazy_static! {
        static ref RE: Regex = Regex::new(r"^([0-9]+)-([0-9]+),([0-9]+)-([0-9]+)$").unwrap();
    }
    let captures = RE.captures(line).ok_or("Invalid line")?;
    Ok((
        Pair::from_str(
            captures.get(1).ok_or("Invalid line")?.as_str(),
            captures.get(2).ok_or("Invalid line")?.as_str(),
        )?,
        Pair::from_str(
            captures.get(3).ok_or("Invalid line")?.as_str(),
            captures.get(4).ok_or("Invalid line")?.as_str(),
        )?,
    ))
}

struct Pair(i32, i32);

impl Pair {
    fn new(x: i32, y: i32) -> Self {
        if x < y {
            Self(x, y)
        } else {
            Self(y, x)
        }
    }

    fn from_str(x: &str, y: &str) -> BoxResult<Self> {
        Ok(Self::new(x.parse::<i32>()?, y.parse::<i32>()?))
    }
}

fn is_fully_contained(a: &Pair, b: &Pair) -> bool {
    is_inclusive(a, b) || is_inclusive(b, a)
}

fn is_inclusive(a: &Pair, b: &Pair) -> bool {
    b.0 >= a.0 && b.1 <= a.1
}

fn is_overlap(a: &Pair, b: &Pair) -> bool {
    let mut intervals = vec![
        Pair(a.0, 0),
        Pair(a.1 + 1, 1),
        Pair(b.0, 0),
        Pair(b.1 + 1, 1),
    ];
    intervals.sort_unstable_by(|a, b| {
        if a.0 < b.0 {
            Ordering::Less
        } else if a.0 > b.0 {
            Ordering::Greater
        } else {
            if a.1 != b.1 {
                if a.1 == 1 {
                    Ordering::Less
                } else {
                    Ordering::Greater
                }
            } else {
                Ordering::Equal
            }
        }
    });

    let mut count = 0;
    for Pair(_, y) in intervals {
        if y == 0 {
            count += 1;
            if count > 1 {
                return true;
            }
        } else {
            count -= 1;
        }
    }
    false
}
