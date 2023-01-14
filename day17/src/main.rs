use std::collections::hash_map::Entry;
use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;

const PUZZLEINPUT: &str = "input.txt";
const PUZZLE_PART1: usize = 2022;
const PUZZLE_PART2: usize = 1_000_000_000_000;

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let mut file = File::open(PUZZLEINPUT)?;
    let mut puzzle_bytes = Vec::new();
    file.read_to_end(&mut puzzle_bytes)?;
    if let Some(v) = puzzle_bytes.last() {
        if v.is_ascii_whitespace() {
            puzzle_bytes.pop();
        }
    }

    let mut sim = Sim::new();
    let mut shape_count = 0;
    let mut next_shape = 0;
    let mut next_byte = 0;

    let mut seen_states = HashMap::new();

    let mut found_cycle = false;
    let mut skipped = false;

    let mut d_shape = 0;
    let mut d_height = 0;

    let mut skipped_shapes = 0;
    let mut skipped_height = 0;

    loop {
        if sim.kind == None {
            if shape_count == PUZZLE_PART1 {
                println!("Part 1: {}", sim.top);
            }

            if !found_cycle {
                let state = next_byte * SHAPES.len() + next_shape;
                match seen_states.entry(state) {
                    Entry::Occupied(mut e) => {
                        let v: &mut SimState = e.get_mut();
                        if v.count == 2 {
                            d_shape = shape_count - v.shape_count;
                            d_height = sim.top - v.height;
                            found_cycle = true;
                        }
                        *v = SimState {
                            count: v.count + 1,
                            shape_count,
                            height: sim.top,
                        }
                    }
                    Entry::Vacant(e) => {
                        e.insert(SimState {
                            count: 1,
                            shape_count,
                            height: sim.top,
                        });
                    }
                }
            }

            if !skipped && found_cycle && shape_count >= PUZZLE_PART1 {
                let cycles = (PUZZLE_PART2 - shape_count) / d_shape;
                skipped_shapes = cycles * d_shape;
                skipped_height = cycles * d_height;
                skipped = true;
            }

            if shape_count + skipped_shapes == PUZZLE_PART2 {
                println!("Part 2: {}", sim.top + skipped_height);
                break;
            }

            sim.add_shape(next_shape);
            shape_count += 1;
            next_shape = (next_shape + 1) % SHAPES.len();
        }

        let b = puzzle_bytes[next_byte];
        next_byte = (next_byte + 1) % puzzle_bytes.len();
        let push_right = match b {
            b'<' => false,
            b'>' => true,
            _ => return Err("Invalid dir".into()),
        };
        sim.push_dir(push_right);
        if !sim.fall() {
            sim.commit_shape();
        }
    }

    Ok(())
}

struct SimState {
    count: u8,
    shape_count: usize,
    height: usize,
}

struct Pos {
    y: usize,
    x: usize,
}

struct Sim {
    pos: Pos,
    kind: Option<usize>,
    grid: Vec<Vec<u8>>,
    top: usize,
}

const SHAPES: &[&[&[u8]]] = &[
    &[b"####"],
    &[b".#.", b"###", b".#."],
    &[b"..#", b"..#", b"###"],
    &[b"#", b"#", b"#", b"#"],
    &[b"##", b"##"],
];

impl Sim {
    fn new() -> Self {
        let mut s = Self {
            pos: Pos { y: 0, x: 0 },
            kind: None,
            grid: Vec::new(),
            top: 0,
        };
        s.add_rows();
        s
    }

    fn add_rows(&mut self) {
        let t = self.top + 7;
        while self.grid.len() < t {
            self.grid.push(".......".into());
        }
    }

    fn add_shape(&mut self, kind: usize) {
        self.pos = Pos {
            y: self.top + 3,
            x: 2,
        };
        self.kind = Some(kind);
    }

    fn commit_shape(&mut self) {
        let kind = match self.kind {
            Some(v) => v,
            None => return,
        };
        let shape = SHAPES[kind];
        let h = shape.len();
        for (yp, i) in shape.iter().enumerate() {
            for (x, &j) in i.iter().enumerate() {
                if j == b'#' {
                    let y = h - yp - 1;
                    let ny = self.pos.y + y;
                    self.grid[ny][self.pos.x + x] = b'#';
                    let t = ny + 1;
                    if t > self.top {
                        self.top = t;
                    }
                }
            }
        }
        self.kind = None;
        self.add_rows();
    }

    fn push_dir(&mut self, right: bool) {
        let dir = if right { 1 } else { -1 };
        if self.check_shape_collision(0, dir) {
            return;
        }
        self.pos = Pos {
            y: self.pos.y,
            x: self.pos.x.wrapping_add_signed(dir),
        };
    }

    fn fall(&mut self) -> bool {
        if self.check_shape_collision(-1, 0) {
            return false;
        }
        self.pos = Pos {
            y: self.pos.y - 1,
            x: self.pos.x,
        };
        true
    }

    fn check_shape_collision(&self, dy: isize, dx: isize) -> bool {
        let kind = match self.kind {
            Some(v) => v,
            None => return false,
        };
        let shape = SHAPES[kind];
        let h = shape.len();
        for (yp, i) in shape.iter().enumerate() {
            for (x, &j) in i.iter().enumerate() {
                if j == b'#' {
                    let y = h - yp - 1;
                    if self.is_grid_block(
                        (self.pos.y + y).wrapping_add_signed(dy),
                        (self.pos.x + x).wrapping_add_signed(dx),
                    ) {
                        return true;
                    }
                }
            }
        }
        false
    }

    fn _is_shape_block(&self, kind: usize, y: usize, x: usize) -> bool {
        let shape = SHAPES[kind];
        if y >= shape.len() {
            return false;
        }
        let row = shape[shape.len() - y - 1];
        if x >= row.len() {
            return false;
        }
        row[x] == b'#'
    }

    fn is_grid_block(&self, y: usize, x: usize) -> bool {
        if y >= self.grid.len() {
            return true;
        }
        let row = &self.grid[y];
        if x >= row.len() {
            return true;
        }
        row[x] == b'#'
    }

    fn _render(&self) {
        let mut i = self.grid.len() - 1;
        while i < self.grid.len() {
            if self.kind != None && i >= self.pos.y && i < self.pos.y + 4 {
                for j in 0..7 {
                    if self._is_shape_block(self.kind.unwrap(), i - self.pos.y, j - self.pos.x) {
                        print!("@");
                    } else if self.is_grid_block(i, j) {
                        print!("#");
                    } else {
                        print!(".");
                    }
                }
                println!();
            } else {
                println!("{}", String::from_utf8_lossy(&self.grid[i]));
            }
            i = i.wrapping_sub(1);
        }
    }
}
