use std::cmp::Ordering;
use std::collections::{BinaryHeap, HashMap};
use std::hash::Hash;

#[derive(Clone, PartialEq, Eq)]
struct Node<K>
where
    K: Clone + Ord,
{
    k: K,
    g: usize,
    h: usize,
}

impl<K> Node<K>
where
    K: Clone + Ord,
{
    fn f(&self) -> usize {
        self.g + self.h
    }
}

impl<K> PartialOrd for Node<K>
where
    K: Clone + Ord,
{
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl<K> Ord for Node<K>
where
    K: Clone + Ord,
{
    fn cmp(&self, other: &Self) -> Ordering {
        other
            .f()
            .cmp(&self.f())
            .then_with(|| other.g.cmp(&self.g))
            .then_with(|| self.k.cmp(&other.k))
    }
}

pub struct Edge<K> {
    pub key: K,
    pub dg: usize,
}

pub trait Neighborer<K> {
    fn neighbors(&self, k: &K) -> Vec<Edge<K>>;
}

pub fn search<K>(
    start: &[K],
    goal: &K,
    neighborer: &dyn Neighborer<K>,
    heuristic: fn(a: &K, b: &K) -> usize,
) -> Option<(Vec<K>, usize)>
where
    K: Clone + Ord + Hash,
{
    let mut open = BinaryHeap::new();
    let mut gscore = HashMap::new();
    let mut adjacent = HashMap::<K, K>::new();
    for i in start {
        open.push(Node {
            k: i.clone(),
            g: 0,
            h: heuristic(i, goal),
        });
        gscore.insert(i.clone(), 0);
    }
    while let Some(current) = open.pop() {
        if &current.k == goal {
            let mut revpath = vec![current.k.clone()];
            let mut k = current.k;
            while let Some(i) = adjacent.remove(&k) {
                revpath.push(i.clone());
                k = i;
            }
            revpath.reverse();
            return Some((revpath, current.g));
        }

        for i in neighborer.neighbors(&current.k) {
            let ng = current.g + i.dg;
            if let Some(&g) = gscore.get(&i.key) {
                if ng >= g {
                    continue;
                }
            }
            adjacent.insert(i.key.clone(), current.k.clone());
            gscore.insert(i.key.clone(), ng);
            open.push(Node {
                k: i.key.clone(),
                g: ng,
                h: heuristic(&i.key, goal),
            });
        }
    }
    None
}
