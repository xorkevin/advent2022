use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let (score1, score2) = {
        let mut score1 = 0;
        let mut score2 = 0;
        for line in reader.lines() {
            let line = line?;
            let (a, b) = if let [a, b] = line.split_ascii_whitespace().collect::<Vec<_>>()[..] {
                (to_move(a)?, to_move(b)?)
            } else {
                return Err("Invalid line".into());
            };
            score1 += winner_score(a, b) + b + 1;
            score2 += b * 3 + pick_move(a, b) + 1;
        }
        (score1, score2)
    };
    println!("Part 1: {}\nPart 2: {}", score1, score2);
    Ok(())
}

fn pick_move(a: i32, b: i32) -> i32 {
    (a + b + 2) % 3
}

fn winner_score(a: i32, b: i32) -> i32 {
    ((b - a + 3 + 1) % 3) * 3
}

fn to_move(a: &str) -> Result<i32, Box<dyn std::error::Error>> {
    match a {
        "A" | "X" => Ok(0),
        "B" | "Y" => Ok(1),
        "C" | "Z" => Ok(2),
        _ => Err("Invalid move".into()),
    }
}
