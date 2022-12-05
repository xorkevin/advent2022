use lazy_static::lazy_static;
use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut grid1 = Grid::empty();
    let mut grid2 = Grid::empty();

    let mut rows = Vec::new();

    let mut mode_grid = true;

    for line in reader.lines() {
        let line = line?;
        if mode_grid {
            if line.len() == 0 {
                grid1 = Grid::from_rows(&rows)?;
                grid2 = grid1.clone();
                mode_grid = false;
                continue;
            }
            let row = parse_grid_row(line.as_bytes())?;
            if row.len() != 0 {
                rows.push(row);
            }
            continue;
        }
        let instr = Instr::from_str(&line)?;
        grid1.process_instr_1(&instr)?;
        grid2.process_instr_2(&instr)?;
    }
    println!("Part 1: {}\nPart 2: {}", grid1.tops()?, grid2.tops()?);
    Ok(())
}

#[derive(Clone)]
struct Grid {
    grid: Vec<Vec<u8>>,
}

impl Grid {
    fn empty() -> Self {
        Self { grid: Vec::new() }
    }

    fn from_rows(rows: &[Vec<u8>]) -> BoxResult<Self> {
        let h = rows.len();
        if h == 0 {
            return Err("No rows".into());
        }
        let w = rows[0].len();
        for i in rows {
            if i.len() != w {
                return Err("Mismatched rows".into());
            }
        }
        let mut grid = Vec::with_capacity(w);
        for i in 0..w {
            let mut col = Vec::with_capacity(h);
            for j in 0..h {
                let c = rows[h - j - 1][i];
                if c == b'.' {
                    break;
                }
                col.push(c);
            }
            grid.push(col);
        }
        Ok(Self { grid })
    }

    fn process_instr_1(&mut self, &Instr(a, b, c): &Instr) -> BoxResult<()> {
        for _ in 0..a {
            let k = self.pop(b)?;
            self.push(c, k)
        }
        Ok(())
    }

    fn process_instr_2(&mut self, &Instr(a, b, c): &Instr) -> BoxResult<()> {
        let mut stack = Vec::with_capacity(a);
        for _ in 0..a {
            stack.push(self.pop(b)?);
        }
        for &i in stack.iter().rev() {
            self.push(c, i)
        }
        Ok(())
    }

    fn pop(&mut self, col: usize) -> BoxResult<u8> {
        self.grid[col].pop().ok_or("No more items".into())
    }

    fn push(&mut self, col: usize, b: u8) {
        self.grid[col].push(b)
    }

    fn peek(&self, col: usize) -> BoxResult<u8> {
        Ok(self.grid[col]
            .last()
            .ok_or::<BoxError>("No more items".into())?
            .clone())
    }

    fn tops(&self) -> BoxResult<String> {
        let mut s = String::with_capacity(self.grid.len());
        for i in 0..self.grid.len() {
            s.push(self.peek(i)? as char);
        }
        Ok(s)
    }
}

struct Instr(usize, usize, usize);

impl Instr {
    fn from_str(line: &str) -> BoxResult<Self> {
        lazy_static! {
            static ref RE: Regex = Regex::new(r"^move (\d+) from (\d+) to (\d+)$").unwrap();
        }
        let captures = RE.captures(line).ok_or("Invalid line")?;
        Ok(Self(
            captures.get(1).ok_or("Invalid line")?.as_str().parse()?,
            captures
                .get(2)
                .ok_or("Invalid line")?
                .as_str()
                .parse::<usize>()?
                - 1,
            captures
                .get(3)
                .ok_or("Invalid line")?
                .as_str()
                .parse::<usize>()?
                - 1,
        ))
    }
}

fn parse_grid_row(line: &[u8]) -> BoxResult<Vec<u8>> {
    let mut row = Vec::with_capacity((line.len() / 4) + 1);
    for part in line.chunks(4) {
        if part.len() < 3 {
            break;
        }
        match part {
            &[b' ', v, ..] => {
                if v != b' ' {
                    return Ok(Vec::new());
                } else {
                    row.push(b'.');
                }
            }
            &[b'[', v, ..] => row.push(v),
            _ => return Err("Invalid line".into()),
        }
    }
    Ok(row)
}
