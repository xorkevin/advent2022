use std::cmp::Reverse;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let nums = {
        let mut nums = Vec::new();
        let mut current = 0;
        for line in reader.lines() {
            let line = line?;
            if line == "" {
                nums.push(current);
                current = 0;
                continue;
            }
            current += line.parse::<i32>()?;
        }
        nums.push(current);
        nums.sort_by_key(|&k| Reverse(k));
        nums
    };

    if let [a, ..] = nums[..] {
        println!("Part 1: {}", a);
    }
    if let [a, b, c, ..] = nums[..] {
        println!("Part 2: {}", a + b + c);
    }
    Ok(())
}
