use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut rope1 = Rope::new(1);
    let mut rope2 = Rope::new(9);

    for line in reader.lines() {
        let line = line?;
        let (dir, countstr) = line
            .split_once(' ')
            .ok_or::<BoxError>("Invalid line".into())?;
        let count = countstr.parse::<i32>()?;
        for _ in 0..count {
            rope1.apply_dir(&dir)?;
            rope2.apply_dir(&dir)?;
        }
    }

    println!(
        "Part 1: {}\nPart 2: {}",
        rope1.history.len(),
        rope2.history.len()
    );

    Ok(())
}

#[derive(Clone, Copy, Eq, Hash, PartialEq)]
struct Tuple2 {
    x: i32,
    y: i32,
}

impl Tuple2 {
    fn apply_dir(&mut self, dir: &str) -> BoxResult<()> {
        match dir {
            "U" => self.y -= 1,
            "R" => self.x += 1,
            "D" => self.y += 1,
            "L" => self.x -= 1,
            _ => return Err("Invalid direction".into()),
        }
        Ok(())
    }

    fn delta(&mut self, p: &Self) {
        self.x += p.x;
        self.y += p.y;
    }

    fn dist(&self, p: &Self) -> Self {
        Self {
            x: p.x - self.x,
            y: p.y - self.y,
        }
    }

    fn dir(&self) -> Self {
        Self {
            x: unit_dir(self.x),
            y: unit_dir(self.y),
        }
    }

    fn max_mag(&self) -> i32 {
        self.x.abs().max(self.y.abs())
    }
}

fn unit_dir(a: i32) -> i32 {
    if a == 0 {
        0
    } else if a > 0 {
        1
    } else {
        -1
    }
}

struct Rope {
    h: Tuple2,
    t: Vec<Tuple2>,
    history: HashSet<Tuple2>,
}

impl Rope {
    fn new(size: usize) -> Self {
        let mut history = HashSet::new();
        history.insert(Tuple2 { x: 0, y: 0 });
        Self {
            h: Tuple2 { x: 0, y: 0 },
            t: vec![Tuple2 { x: 0, y: 0 }; size],
            history,
        }
    }

    fn apply_dir(&mut self, dir: &str) -> BoxResult<()> {
        self.h.apply_dir(dir)?;
        let mut next = &self.h;
        let last = self.t.len() - 1;
        for i in 0..self.t.len() {
            let k = self.t[i].dist(&next);
            if k.max_mag() > 1 {
                self.t[i].delta(&k.dir());
                if i == last {
                    self.history.insert(self.t[i]);
                }
            }
            next = &self.t[i];
        }
        Ok(())
    }
}
