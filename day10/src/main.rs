use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut vm = VM::new();

    for line in reader.lines() {
        let line = line?;
        let (instr, arg) = match line.split_ascii_whitespace().collect::<Vec<_>>()[..] {
            [instr] => (instr, 0),
            [instr, arg] => (instr, arg.parse::<i32>()?),
            _ => return Err("Invalid instruction line".into()),
        };
        vm.exec(instr, arg)?;
    }

    println!("Part 1: {}\nPart 2:", vm.strength);
    for line in vm.grid {
        println!("{}", String::from_utf8(line)?);
    }

    Ok(())
}

const SCREEN_WIDTH: usize = 40;
const SCREEN_HEIGHT: usize = 6;
const SCREEN_SIZE: usize = SCREEN_WIDTH * SCREEN_HEIGHT;

struct VM {
    cycle: i32,
    rx: i32,
    strength: i32,
    cycle_target: i32,
    cycle_incr: i32,
    grid: Vec<Vec<u8>>,
    scanline: usize,
}

impl VM {
    fn new() -> Self {
        Self {
            cycle: 0,
            rx: 1,
            strength: 0,
            cycle_target: 20,
            cycle_incr: 40,
            grid: vec![vec![0; SCREEN_WIDTH]; SCREEN_HEIGHT],
            scanline: 0,
        }
    }

    fn exec(&mut self, instr: &str, arg: i32) -> BoxResult<()> {
        match instr {
            "noop" => self.exec_noop(),
            "addx" => self.exec_addx(arg),
            _ => return Err("Invalid instruction".into()),
        }
        Ok(())
    }

    fn intersect(&self, x: i32) -> bool {
        (x - self.rx).abs() <= 1
    }

    fn check_cycle(&mut self, cycle: i32, set: bool, next: i32) {
        self.cycle += cycle;
        if self.cycle >= self.cycle_target {
            self.strength += self.cycle_target * self.rx;
            self.cycle_target += self.cycle_incr;
        }
        let y = self.scanline / SCREEN_WIDTH;
        let x = self.scanline % SCREEN_WIDTH;
        self.grid[y][x] = if self.intersect(x as i32) { b'#' } else { b'.' };
        self.scanline = (self.scanline + 1) % SCREEN_SIZE;
        if set {
            self.rx = next;
        }
    }

    fn exec_noop(&mut self) {
        self.check_cycle(1, false, 0)
    }

    fn exec_addx(&mut self, arg: i32) {
        self.check_cycle(1, false, 0);
        self.check_cycle(1, true, self.rx + arg)
    }
}
