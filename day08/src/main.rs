use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut cells = Vec::new();
    for line in reader.lines() {
        let line = line?;
        let mut row = Vec::with_capacity(line.len());
        for &b in line.as_bytes() {
            if b < b'0' || b > b'9' {
                return Err("Invalid grid cell".into());
            }
            row.push((b - b'0') as i32);
        }
        cells.push(row);
    }

    let mut grid = Grid::new(cells)?;
    grid.compute_visible_set();

    println!(
        "Part 1: {}\nPart 2: {}",
        grid.visible.len(),
        grid.max_power()
    );

    Ok(())
}

struct Grid {
    w: usize,
    h: usize,
    grid: Vec<Vec<i32>>,
    visible: HashSet<Tuple2>,
    power: HashMap<Tuple2, Tuple5>,
}

impl Grid {
    fn new(grid: Vec<Vec<i32>>) -> BoxResult<Self> {
        let h = grid.len();
        if h == 0 {
            return Err("Empty grid".into());
        }
        let w = grid[0].len();
        for i in &grid {
            if i.len() != w {
                return Err("Not rectangular grid".into());
            }
        }
        Ok(Self {
            w,
            h,
            grid,
            visible: HashSet::new(),
            power: HashMap::new(),
        })
    }

    fn max_power(&self) -> i32 {
        self.power.values().fold(0, |acc, i| {
            let k = i.get_power();
            if k > acc {
                k
            } else {
                acc
            }
        })
    }

    fn compute_visible_set(&mut self) {
        for y in 0..self.h {
            let mut tallest = -1;
            for x in 0..self.w {
                let pos = Tuple2 {
                    x: x as i32,
                    y: y as i32,
                };
                let k = self.grid[y][x];

                if k > tallest {
                    tallest = k;
                    self.visible.insert(pos);
                }

                self.power.insert(
                    pos,
                    Tuple5 {
                        h: k,
                        t: 0,
                        r: 0,
                        b: 0,
                        l: 0,
                    },
                );

                if x == 0 {
                    continue;
                }

                let mut visible = 1;
                let mut prev = pos.delta(-1, 0);
                let mut prev_power = self.power.get(&prev).unwrap();
                while k > prev_power.h && prev_power.l > 0 {
                    visible += prev_power.l;
                    prev = prev.delta(-prev_power.l, 0);
                    prev_power = self.power.get(&prev).unwrap();
                }
                self.power.get_mut(&pos).unwrap().l = visible;
            }

            tallest = -1;
            for x in (0..self.w).rev() {
                let pos = Tuple2 {
                    x: x as i32,
                    y: y as i32,
                };
                let k = self.grid[y][x];

                if k > tallest {
                    tallest = k;
                    self.visible.insert(pos);
                }

                if x == self.w - 1 {
                    continue;
                }

                let mut visible = 1;
                let mut prev = pos.delta(1, 0);
                let mut prev_power = self.power.get(&prev).unwrap();
                while k > prev_power.h && prev_power.r > 0 {
                    visible += prev_power.r;
                    prev = prev.delta(prev_power.r, 0);
                    prev_power = self.power.get(&prev).unwrap();
                }
                self.power.get_mut(&pos).unwrap().r = visible;
            }
        }

        for x in 0..self.w {
            let mut tallest = -1;
            for y in 0..self.h {
                let pos = Tuple2 {
                    x: x as i32,
                    y: y as i32,
                };
                let k = self.grid[y][x];

                if k > tallest {
                    tallest = k;
                    self.visible.insert(pos);
                }

                if y == 0 {
                    continue;
                }

                let mut visible = 1;
                let mut prev = pos.delta(0, -1);
                let mut prev_power = self.power.get(&prev).unwrap();
                while k > prev_power.h && prev_power.t > 0 {
                    visible += prev_power.t;
                    prev = prev.delta(0, -prev_power.t);
                    prev_power = self.power.get(&prev).unwrap();
                }
                self.power.get_mut(&pos).unwrap().t = visible;
            }

            tallest = -1;
            for y in (0..self.h).rev() {
                let pos = Tuple2 {
                    x: x as i32,
                    y: y as i32,
                };
                let k = self.grid[y][x];

                if k > tallest {
                    tallest = k;
                    self.visible.insert(pos);
                }

                if y == self.h - 1 {
                    continue;
                }

                let mut visible = 1;
                let mut prev = pos.delta(0, 1);
                let mut prev_power = self.power.get(&prev).unwrap();
                while k > prev_power.h && prev_power.b > 0 {
                    visible += prev_power.b;
                    prev = prev.delta(0, prev_power.b);
                    prev_power = self.power.get(&prev).unwrap();
                }
                self.power.get_mut(&pos).unwrap().b = visible;
            }
        }
    }
}

struct Tuple5 {
    h: i32,
    t: i32,
    r: i32,
    b: i32,
    l: i32,
}

impl Tuple5 {
    fn get_power(&self) -> i32 {
        self.t * self.r * self.b * self.l
    }
}

#[derive(Clone, Copy, Eq, Hash, PartialEq)]
struct Tuple2 {
    y: i32,
    x: i32,
}

impl Tuple2 {
    fn delta(&self, x: i32, y: i32) -> Self {
        Self {
            x: self.x + x,
            y: self.y + y,
        }
    }
}
