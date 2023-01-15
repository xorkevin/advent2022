use std::collections::HashSet;
use std::fs::File;
use std::hash::Hash;
use std::io::prelude::*;
use std::io::BufReader;
use std::ops::Add;

const PUZZLEINPUT: &str = "input.txt";

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut cloud = HashSet::new();
    let mut start = None;

    for line in reader.lines() {
        let line = line?;
        let p = if let [x, y, z] = line
            .split(",")
            .flat_map(|s| s.parse::<i32>())
            .collect::<Vec<_>>()[..]
        {
            Point::new(x, y, z)
        } else {
            return Err("Invalid line".into());
        };
        cloud.insert(p);
        if start == None {
            start = Some(p);
        }
    }

    let mut start = if let Some(v) = start {
        v
    } else {
        return Err("No points".into());
    };

    let mut border = HashSet::new();

    let mut surface_area = 0;
    for &k in &cloud {
        surface_area += 6 - CARDINAL_DIRS
            .iter()
            .filter(|&&i| cloud.contains(&(k + i)))
            .count();
        for &i in INTER_CARDINAL_DIRS {
            let p = k + i;
            if !cloud.contains(&p) {
                border.insert(p);
                if p.x > start.x {
                    start = p;
                }
            }
        }
    }
    println!("Part 1: {}", surface_area);

    let mut ext_surface_area = 0;
    let mut open_set = vec![start];
    let mut closed_set = HashSet::new();
    closed_set.insert(start);
    while let Some(p) = open_set.pop() {
        for &i in CARDINAL_DIRS {
            let k = p.add(i);
            if cloud.contains(&k) {
                ext_surface_area += 1;
            } else if border.contains(&k) && !closed_set.contains(&k) {
                open_set.push(k);
                closed_set.insert(k);
            }
        }
    }
    println!("Part 2: {}", ext_surface_area);

    Ok(())
}

#[derive(Clone, Copy, PartialEq, Eq, Hash)]
struct Point {
    x: i32,
    y: i32,
    z: i32,
}

impl Point {
    const fn new(x: i32, y: i32, z: i32) -> Self {
        Point { x, y, z }
    }
}

impl Add for Point {
    type Output = Self;

    fn add(self, other: Self) -> Self {
        Self {
            x: self.x + other.x,
            y: self.y + other.y,
            z: self.z + other.z,
        }
    }
}

const CARDINAL_DIRS: &[Point] = &[
    Point::new(1, 0, 0),
    Point::new(-1, 0, 0),
    Point::new(0, 1, 0),
    Point::new(0, -1, 0),
    Point::new(0, 0, 1),
    Point::new(0, 0, -1),
];

const INTER_CARDINAL_DIRS: &[Point] = &[
    Point::new(-1, -1, -1),
    Point::new(-1, -1, 0),
    Point::new(-1, -1, 1),
    Point::new(-1, 0, -1),
    Point::new(-1, 0, 0),
    Point::new(-1, 0, 1),
    Point::new(-1, 1, -1),
    Point::new(-1, 1, 0),
    Point::new(-1, 1, 1),
    Point::new(0, -1, -1),
    Point::new(0, -1, 0),
    Point::new(0, -1, 1),
    Point::new(0, 0, -1),
    Point::new(0, 0, 1),
    Point::new(0, 1, -1),
    Point::new(0, 1, 0),
    Point::new(0, 1, 1),
    Point::new(1, -1, -1),
    Point::new(1, -1, 0),
    Point::new(1, -1, 1),
    Point::new(1, 0, -1),
    Point::new(1, 0, 0),
    Point::new(1, 0, 1),
    Point::new(1, 1, -1),
    Point::new(1, 1, 0),
    Point::new(1, 1, 1),
];
