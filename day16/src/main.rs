use regex::Regex;
use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

mod floydwarshall;

const PUZZLEINPUT: &str = "input.txt";

type BoxError = Box<dyn std::error::Error>;
type BoxResult<T> = Result<T, BoxError>;

fn main() -> BoxResult<()> {
    let line_regex =
        Regex::new(r"Valve (\w+) has flow rate=(\d+); tunnels? leads? to valves? (.*)").unwrap();

    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut valves = Vec::new();
    let mut nodes = Vec::new();
    let mut edges = Vec::new();

    for line in reader.lines() {
        let line = line?;
        let captures = line_regex.captures(&line).ok_or("Invalid line")?;
        let name = captures.get(1).ok_or("Invalid line")?.as_str().to_string();
        let rate = captures
            .get(2)
            .ok_or("Invalid line")?
            .as_str()
            .parse::<i32>()?;
        nodes.push(name.clone());
        if rate > 0 {
            valves.push(Valve {
                name: name.clone(),
                id: 0,
                rate,
            })
        }
        for i in captures.get(3).ok_or("Invalid line")?.as_str().split(", ") {
            edges.push(floydwarshall::Edge {
                a: name.clone(),
                b: i.to_string(),
                c: 1,
            });
        }
    }

    valves.sort_by(|a, b| b.rate.cmp(&a.rate));
    for (n, i) in valves.iter_mut().enumerate() {
        i.id = 1 << n;
    }

    let dist = floydwarshall::compute(&nodes[..], &edges[..]);
    println!(
        "Part 1: {}",
        search_max(0, 30, "AA", &mut HashSet::new(), &valves, &dist, 0)
    );

    let mut path_map = HashMap::new();
    search_paths(
        0,
        0,
        26,
        "AA",
        &mut HashSet::new(),
        &mut path_map,
        &valves,
        &dist,
    );
    let mut all_paths = path_map
        .into_iter()
        .map(|(id, flow)| Path { id, flow })
        .collect::<Vec<_>>();
    all_paths.sort_by(|a, b| b.flow.cmp(&a.flow));

    let mut max_flow = 0;
    for (n, i) in all_paths.iter().enumerate() {
        if i.flow * 2 <= max_flow {
            break;
        }
        for j in &all_paths[n + 1..] {
            let flow = i.flow + j.flow;
            if flow <= max_flow {
                break;
            }
            if i.id & j.id != 0 {
                continue;
            }
            if flow > max_flow {
                max_flow = flow;
            }
        }
    }
    println!("Part 2: {}", max_flow);

    Ok(())
}

struct Path {
    id: u32,
    flow: i32,
}

struct Valve {
    name: String,
    id: u32,
    rate: i32,
}

fn search_max(
    acc: i32,
    remaining: i32,
    pos: &str,
    toggled: &mut HashSet<String>,
    valves: &[Valve],
    dist: &floydwarshall::PairwiseDistances<String>,
    mut candidate: i32,
) -> i32 {
    if remaining <= 0 {
        return acc;
    }
    {
        let mut bound = 0;
        let mut t = remaining;
        for i in valves {
            if t <= 2 {
                break;
            }
            if toggled.contains(&i.name) {
                continue;
            }
            t -= 2;
            bound += t * i.rate;
        }
        if acc + bound < candidate {
            return 0;
        }
    }

    let mut max_flow = acc;
    if max_flow > candidate {
        candidate = max_flow;
    }

    for i in valves {
        if toggled.contains(&i.name) {
            continue;
        }
        let cost = match dist.edge_cost(&pos.to_string(), &i.name) {
            Some(c) => c,
            None => continue,
        };
        let next_remaining = remaining - cost - 1;
        if next_remaining <= 0 {
            continue;
        }
        toggled.insert(i.name.clone());
        let flow = search_max(
            acc + i.rate * next_remaining,
            next_remaining,
            &i.name,
            toggled,
            valves,
            dist,
            candidate,
        );
        toggled.remove(&i.name);
        if flow > max_flow {
            max_flow = flow;
            if max_flow > candidate {
                candidate = max_flow;
            }
        }
    }

    max_flow
}

fn search_paths(
    acc: i32,
    cur_path: u32,
    remaining: i32,
    pos: &str,
    toggled: &mut HashSet<String>,
    all_paths: &mut HashMap<u32, i32>,
    valves: &[Valve],
    dist: &floydwarshall::PairwiseDistances<String>,
) {
    if remaining <= 0 {
        return;
    }
    if match all_paths.get(&cur_path) {
        Some(&v) => acc > v,
        None => true,
    } {
        all_paths.insert(cur_path, acc);
    }

    for i in valves {
        if toggled.contains(&i.name) {
            continue;
        }
        let cost = match dist.edge_cost(&pos.to_string(), &i.name) {
            Some(c) => c,
            None => continue,
        };
        let next_remaining = remaining - cost - 1;
        if next_remaining <= 0 {
            continue;
        }
        toggled.insert(i.name.clone());
        search_paths(
            acc + i.rate * next_remaining,
            cur_path | i.id,
            next_remaining,
            &i.name,
            toggled,
            all_paths,
            valves,
            dist,
        );
        toggled.remove(&i.name);
    }
}
