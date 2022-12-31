use std::collections::hash_map::Entry;
use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::ops::{Add, AddAssign};

const PUZZLEINPUT: &str = "input.txt";

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut grid1 = HashSet::new();
    let mut grid2 = HashSet::new();
    let mut lowest = HashMap::new();
    let mut floor = 0;

    for line in reader.lines() {
        let mut first = true;
        let mut last = Pos::new(0, 0);
        for i in line?.split(" -> ") {
            let pos = if let Some((x, y)) = i.split_once(",") {
                Pos::new(x.parse::<i32>()?, y.parse::<i32>()?)
            } else {
                return Err("Invalid point".into());
            };
            if first {
                first = false;
                last = pos;
                grid1.insert(last);
                grid2.insert(last);
                match lowest.entry(last.x) {
                    Entry::Occupied(mut e) => {
                        let v = e.get_mut();
                        if last.y > *v {
                            *v = last.y;
                            if last.y > floor {
                                floor = last.y;
                            }
                        }
                    }
                    Entry::Vacant(e) => {
                        e.insert(last.y);
                        if last.y > floor {
                            floor = last.y;
                        }
                    }
                }
            } else {
                let delta = last.unit_delta(pos);
                while last != pos {
                    last += delta;
                    grid1.insert(last);
                    grid2.insert(last);
                    match lowest.entry(last.x) {
                        Entry::Occupied(mut e) => {
                            let v = e.get_mut();
                            if last.y > *v {
                                *v = last.y;
                                if last.y > floor {
                                    floor = last.y;
                                }
                            }
                        }
                        Entry::Vacant(e) => {
                            e.insert(last.y);
                            if last.y > floor {
                                floor = last.y;
                            }
                        }
                    }
                }
            }
        }
    }

    floor += 2;

    {
        let mut count = 0;
        while drop_particle1(Pos::new(500, 0), &mut grid1, &lowest) {
            count += 1;
        }
        println!("Part 1: {}", count);
    }

    {
        let mut count = 0;
        while drop_particle2(Pos::new(500, 0), &mut grid2, floor) {
            count += 1;
        }
        println!("Part 2: {}", count);
    }

    Ok(())
}

const DIRS: [Pos; 3] = [Pos::new(0, 1), Pos::new(-1, 1), Pos::new(1, 1)];

fn drop_particle1(mut a: Pos, grid: &mut HashSet<Pos>, lowest: &HashMap<i32, i32>) -> bool {
    if grid.contains(&a) {
        return false;
    }
    'outer: loop {
        match lowest.get(&a.x) {
            Some(&v) => {
                if a.y > v {
                    return false;
                }
            }
            None => return false,
        }
        for i in DIRS {
            let next = a + i;
            if !grid.contains(&next) {
                a = next;
                continue 'outer;
            }
        }
        grid.insert(a);
        return true;
    }
}

fn drop_particle2(mut a: Pos, grid: &mut HashSet<Pos>, floor: i32) -> bool {
    if grid.contains(&a) {
        return false;
    }
    'outer: loop {
        if a.y + 1 < floor {
            for i in DIRS {
                let next = a + i;
                if !grid.contains(&next) {
                    a = next;
                    continue 'outer;
                }
            }
        }
        grid.insert(a);
        return true;
    }
}

#[derive(Clone, Copy, Hash, PartialEq, Eq)]
struct Pos {
    y: i32,
    x: i32,
}

fn unit_dir(a: i32, b: i32) -> i32 {
    if a < b {
        1
    } else if a > b {
        -1
    } else {
        0
    }
}

impl Pos {
    const fn new(x: i32, y: i32) -> Self {
        Self { y, x }
    }

    fn unit_delta(&self, other: Self) -> Self {
        Self {
            y: unit_dir(self.y, other.y),
            x: unit_dir(self.x, other.x),
        }
    }
}

impl Add for Pos {
    type Output = Self;

    fn add(self, other: Self) -> Self {
        Self {
            y: self.y + other.y,
            x: self.x + other.x,
        }
    }
}

impl AddAssign for Pos {
    fn add_assign(&mut self, other: Self) {
        *self = Self {
            y: self.y + other.y,
            x: self.x + other.x,
        }
    }
}
