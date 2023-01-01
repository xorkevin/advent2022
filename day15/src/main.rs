use regex::Regex;
use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";
const PUZZLE_ROW: i32 = 2000000;
const PUZZLE_BOUND: i32 = 4000000;

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let line_regex = Regex::new(r"x=(-?\d+).*y=(-?\d+).*x=(-?\d+).*y=(-?\d+)").unwrap();

    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut sensors = Vec::new();
    let mut beacons = HashSet::new();
    let mut bounds = None;

    for line in reader.lines() {
        let line = line?;
        let captures = line_regex.captures(&line).ok_or("Invalid line")?;
        let pos = Pos::new(
            captures.get(1).ok_or("Invalid line")?.as_str().parse()?,
            captures.get(2).ok_or("Invalid line")?.as_str().parse()?,
        );
        let beacon = Pos::new(
            captures.get(3).ok_or("Invalid line")?.as_str().parse()?,
            captures.get(4).ok_or("Invalid line")?.as_str().parse()?,
        );
        beacons.insert(beacon);
        let radius = pos.manhattan_distance(&beacon);
        let sensor = Sensor::new(pos, radius);
        if let Some((x1, x2)) = sensor.bounds_x(PUZZLE_ROW) {
            match &mut bounds {
                None => bounds = Some((x1, x2)),
                Some((v1, v2)) => {
                    if x1 < *v1 {
                        *v1 = x1
                    }
                    if x2 > *v2 {
                        *v2 = x2
                    }
                }
            }
        }
        sensors.push(sensor);
    }

    {
        let mut count = 0;
        if let Some((left_bound, right_bound)) = bounds {
            let mut x = left_bound;
            while x <= right_bound {
                let pos = Pos::new(x, PUZZLE_ROW);
                if beacons.contains(&pos) {
                    x += 1;
                    continue;
                }
                for i in &sensors {
                    if !i.in_range(&pos) {
                        continue;
                    }
                    let (_, x2) = i.bounds_x(PUZZLE_ROW).ok_or("Invariant violated")?;
                    if x2 == x {
                        count += 1;
                    } else {
                        count += x2 - x;
                        x = x2 - 1;
                    }
                    break;
                }
                x += 1;
            }
        }
        println!("Part 1: {}", count);
    }
    {
        let mut y = 0;
        'outer2: while y <= PUZZLE_BOUND {
            let mut x = 0;
            'outer: while x <= PUZZLE_BOUND {
                let pos = Pos::new(x, y);
                for i in &sensors {
                    if !i.in_range(&pos) {
                        continue;
                    }
                    let (_, x2) = i.bounds_x(y).ok_or("Invariant violated")?;
                    x = x2 + 1;
                    continue 'outer;
                }
                println!("Part 2: {}", x as i64 * PUZZLE_BOUND as i64 + y as i64);
                break 'outer2;
            }
            y += 1
        }
    }

    Ok(())
}

#[derive(Clone, Copy, PartialEq, Eq, Hash)]
struct Pos {
    y: i32,
    x: i32,
}

impl Pos {
    fn new(x: i32, y: i32) -> Self {
        Self { y, x }
    }

    fn manhattan_distance(&self, other: &Self) -> i32 {
        (self.x - other.x).abs() + (self.y - other.y).abs()
    }
}

struct Sensor {
    pos: Pos,
    radius: i32,
}

impl Sensor {
    fn new(pos: Pos, radius: i32) -> Self {
        Self { pos, radius }
    }

    fn in_range(&self, pos: &Pos) -> bool {
        self.pos.manhattan_distance(&pos) <= self.radius
    }

    fn bounds_x(&self, y: i32) -> Option<(i32, i32)> {
        let vdelta = (self.pos.y - y).abs();
        if vdelta > self.radius {
            None
        } else {
            let delta = self.radius - vdelta;
            Some((self.pos.x - delta, self.pos.x + delta))
        }
    }
}
