use lazy_static::lazy_static;
use regex::Regex;
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

    let mut term = Term::new();

    for line in reader.lines() {
        term.read_input(&line?)?;
    }

    let (total_size, small_dirs) = calc_small_dir_size(&mut term.root);

    let mut freed = 0;
    let target = MIN_UNUSED + total_size - TOTAL_DISK;
    if target > 0 {
        if let Some(k) = find_dir_size_target(&term.root, target) {
            freed = k;
        }
    }

    println!("Part 1: {}\nPart 2: {}", small_dirs, freed);

    Ok(())
}

const TOTAL_DISK: i32 = 70000000;
const MIN_UNUSED: i32 = 30000000;

fn find_dir_size_target(n: &Node, target: i32) -> Option<i32> {
    if n.size < target {
        return None;
    }
    let mut at_most = n.size;
    for v in n.children.values() {
        if !v.is_dir {
            continue;
        }
        if let Some(k) = find_dir_size_target(v, target) {
            if k < at_most {
                at_most = k;
            }
        }
    }
    Some(at_most)
}

const SMALL_DIR_LIMIT: i32 = 100000;

fn calc_small_dir_size(n: &mut Node) -> (i32, i32) {
    let mut total = 0;
    let mut cummulative = 0;
    for v in n.children.values_mut() {
        if v.is_dir {
            let (t, c) = calc_small_dir_size(v);
            total += t;
            cummulative += c;
        } else {
            total += v.size;
        }
    }
    n.size = total;
    if total <= SMALL_DIR_LIMIT {
        cummulative += total;
    }
    (total, cummulative)
}

struct Node {
    children: HashMap<String, Box<Node>>,
    is_dir: bool,
    size: i32,
}

struct Term {
    pwd: Vec<String>,
    root: Box<Node>,
    running: Option<String>,
}

impl Node {
    fn new_dir() -> Box<Self> {
        Box::new(Self {
            children: HashMap::new(),
            is_dir: true,
            size: 0,
        })
    }

    fn new_file(size: i32) -> Box<Self> {
        Box::new(Self {
            children: HashMap::new(),
            is_dir: false,
            size,
        })
    }
}

impl Term {
    fn new() -> Self {
        Self {
            pwd: Vec::new(),
            root: Node::new_dir(),
            running: None,
        }
    }

    fn read_input(&mut self, inp: &str) -> BoxResult<()> {
        lazy_static! {
            static ref RE: Regex = Regex::new(r"^\$ ").unwrap();
        }
        if RE.is_match(inp) {
            self.running = None;
            self.exec(&inp[2..].split_ascii_whitespace().collect::<Vec<_>>()[..])
        } else {
            match &self.running {
                Some(running) => match running.as_str() {
                    "ls" => self.read_output_ls(inp),
                    _ => Err("Invalid input".into()),
                },
                None => Err("Invalid input".into()),
            }
        }
    }

    fn exec(&mut self, cmd: &[&str]) -> BoxResult<()> {
        match cmd {
            &["cd", dir] => self.cd(dir),
            &["ls"] => {
                self.running = Some("ls".into());
                Ok(())
            }
            _ => Err("Invalid cmd".into()),
        }
    }

    fn cd(&mut self, dir: &str) -> BoxResult<()> {
        match dir {
            "" => return Err("No cd dir".into()),
            ".." => {
                self.pwd
                    .pop()
                    .ok_or::<BoxError>("No parent directory from root".into())?;
            }
            "/" => self.pwd.clear(),
            d => self.pwd.push(d.into()),
        }
        Ok(())
    }

    fn read_output_ls(&mut self, inp: &str) -> BoxResult<()> {
        let (kind, name) = if let Some((kind, name)) = inp.split_once(' ') {
            (kind, name)
        } else {
            return Err("Invalid ls output".into());
        };
        if name == "" {
            return Err("Invalid ls file name".into());
        }
        if kind == "dir" {
            self.mkdir(name)
        } else {
            self.touch(name, kind.parse()?)
        }
    }

    fn mkdir_path(&mut self) -> BoxResult<&mut Box<Node>> {
        let mut node = &mut self.root;
        for i in self.pwd.iter() {
            match node.children.entry(i.into()) {
                Entry::Occupied(e) => {
                    let v = e.get();
                    if !v.is_dir {
                        return Err("Mkdir invalid path".into());
                    }
                    node = e.into_mut();
                }
                Entry::Vacant(e) => {
                    node = e.insert(Node::new_dir());
                }
            }
        }
        Ok(node)
    }

    fn mkdir(&mut self, name: &str) -> BoxResult<()> {
        let node = self.mkdir_path()?;
        match node.children.entry(name.into()) {
            Entry::Occupied(e) => {
                if !e.get().is_dir {
                    return Err("Mkdir on non-dir".into());
                }
            }
            Entry::Vacant(e) => {
                e.insert(Node::new_dir());
            }
        }
        Ok(())
    }

    fn touch(&mut self, name: &str, size: i32) -> BoxResult<()> {
        let node = self.mkdir_path()?;
        match node.children.entry(name.into()) {
            Entry::Occupied(mut e) => {
                let v = e.get_mut();
                if v.is_dir {
                    return Err("Touch file on dir".into());
                }
                v.size = size;
            }
            Entry::Vacant(e) => {
                e.insert(Node::new_file(size));
            }
        }
        Ok(())
    }
}
