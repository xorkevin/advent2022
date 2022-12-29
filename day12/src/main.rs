use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

mod astar;

const PUZZLEINPUT: &str = "input.txt";

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut start = None;
    let mut end = None;
    let mut starts2 = Vec::new();

    let mut grid = Vec::new();

    for line in reader.lines() {
        grid.push(
            line?
                .bytes()
                .enumerate()
                .map(|(x, i)| {
                    let pos = Pos { y: grid.len(), x };
                    match i {
                        b'S' => {
                            start = Some(pos);
                            starts2.push(pos);
                            b'a'
                        }
                        b'E' => {
                            end = Some(pos);
                            b'z'
                        }
                        b'a' => {
                            starts2.push(pos);
                            b'a'
                        }
                        c => c,
                    }
                })
                .collect::<Vec<_>>(),
        );
    }

    let (start, end) = match (start, end) {
        (Some(start), Some(end)) => (start, end),
        _ => return Err("Missing start or end".into()),
    };

    let grid = Grid::new(grid);

    if let Some((_, steps)) = astar::search(&vec![start], &end, &grid, manhattan_distance) {
        println!("Part 1: {}", steps);
    }

    if let Some((_, steps)) = astar::search(&starts2, &end, &grid, manhattan_distance) {
        println!("Part 2: {}", steps);
    }

    Ok(())
}

#[derive(Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Hash)]
struct Pos {
    y: usize,
    x: usize,
}

struct Grid {
    h: usize,
    w: usize,
    grid: Vec<Vec<u8>>,
}

const DIR_DELTAS: [(isize, isize); 4] = [(-1, 0), (0, 1), (1, 0), (0, -1)];

impl Grid {
    fn new(grid: Vec<Vec<u8>>) -> Self {
        Self {
            h: grid.len(),
            w: grid[0].len(),
            grid,
        }
    }

    fn in_bounds(&self, k: &Pos) -> bool {
        k.x < self.w && k.y < self.h
    }

    fn get(&self, k: &Pos) -> Option<u8> {
        if self.in_bounds(k) {
            Some(self.grid[k.y][k.x])
        } else {
            None
        }
    }
}

impl astar::Neighborer<Pos> for Grid {
    fn neighbors(&self, k: &Pos) -> Vec<astar::Edge<Pos>> {
        let limit = self.grid[k.y][k.x] + 1;
        let mut e = Vec::new();
        for (dy, dx) in DIR_DELTAS {
            let k = Pos {
                y: k.y.wrapping_add_signed(dy),
                x: k.x.wrapping_add_signed(dx),
            };
            if let Some(v) = self.get(&k) {
                if v <= limit {
                    e.push(astar::Edge { value: k, dg: 1 })
                }
            }
        }
        e
    }
}

fn manhattan_distance(a: &Pos, b: &Pos) -> usize {
    a.x.abs_diff(b.x) + a.y.abs_diff(b.y)
}
