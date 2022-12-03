use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let (sum1, sum2) = {
        let mut sum1 = 0;
        let mut sum2 = 0;
        let mut group_size = 0;
        let mut group_common = Vec::new();
        for line in reader.lines() {
            let line = line?.into_bytes();
            if line.len() % 2 != 0 {
                return Err("Invalid line format".into());
            }
            let halflen = line.len() / 2;
            if let Some(c) = find_common(&line[..halflen], &line[halflen..]) {
                sum1 += prio(c)?;
            } else {
                return Err("None in common".into());
            }

            if group_size < 2 {
                if group_size == 0 {
                    group_common = line;
                } else {
                    group_common = find_common_all(&group_common, &line);
                }
                group_size += 1;
                continue;
            }
            if let Some(c) = find_common(&group_common, &line) {
                sum2 += prio(c)?;
            } else {
                return Err("None in common".into());
            }
            group_common.clear();
            group_size = 0;
        }
        (sum1, sum2)
    };
    println!("Part 1: {}\nPart 2: {}", sum1, sum2);
    Ok(())
}

fn find_common_all(a: &[u8], b: &[u8]) -> Vec<u8> {
    let first = HashSet::<_>::from_iter(a.iter());
    let mut common = Vec::new();
    for &i in b {
        if first.contains(&i) {
            common.push(i);
        }
    }
    common
}

fn find_common(a: &[u8], b: &[u8]) -> Option<u8> {
    let first = HashSet::<_>::from_iter(a.iter());
    for &i in b {
        if first.contains(&i) {
            return Some(i);
        }
    }
    None
}

fn prio(c: u8) -> Result<i32, Box<dyn std::error::Error>> {
    if c >= b'a' && c <= b'z' {
        return Ok(c as i32 - 'a' as i32 + 1);
    }
    if c >= b'A' && c <= b'Z' {
        return Ok(c as i32 - 'A' as i32 + 27);
    }
    Err("Invalid prio".into())
}
