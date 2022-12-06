use std::collections::hash_map::Entry;
use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);
    let mut file = Vec::new();
    for b in reader.bytes() {
        file.push(b?);
    }

    let mut startpos1 = None;
    let mut startpos2 = None;

    let mut seen1 = HashMap::new();
    let mut seen2 = HashMap::new();

    for i in 0..file.len() {
        if startpos1.is_some() && startpos2.is_some() {
            break;
        }

        let c = file[i];

        if startpos1.is_none() {
            let v = seen1.entry(c).or_insert(0);
            *v += 1;
            if i > 3 {
                if let Entry::Occupied(mut e) = seen1.entry(file[i - 4]) {
                    let v = e.get_mut();
                    if *v < 2 {
                        e.remove_entry();
                    } else {
                        *v -= 1;
                    }
                }
            }
            if i >= 3 && is_uniq(&seen1) {
                startpos1 = Some(i + 1);
            }
        }

        if startpos2.is_none() {
            let v = seen2.entry(c).or_insert(0);
            *v += 1;
            if i > 13 {
                if let Entry::Occupied(mut e) = seen2.entry(file[i - 14]) {
                    let v = e.get_mut();
                    if *v < 2 {
                        e.remove_entry();
                    } else {
                        *v -= 1;
                    }
                }
            }
            if i >= 13 && is_uniq(&seen2) {
                startpos2 = Some(i + 1);
            }
        }
    }
    println!(
        "Part 1: {}\nPart 2: {}",
        startpos1.ok_or("Failed to find start 1")?,
        startpos2.ok_or("Failed to find start 2")?,
    );
    Ok(())
}

fn is_uniq(seen: &HashMap<u8, usize>) -> bool {
    seen.values().find(|&&v| v > 1).is_none()
}
