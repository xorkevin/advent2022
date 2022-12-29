use std::cmp::Ordering;
use std::collections::{BinaryHeap, HashMap};
use std::hash::Hash;

#[derive(Clone, PartialEq, Eq)]
struct Node<T>
where
    T: Clone + Ord,
{
    v: T,
    g: usize,
    h: usize,
}

impl<T> Node<T>
where
    T: Clone + Ord,
{
    fn f(&self) -> usize {
        self.g + self.h
    }
}

impl<T> PartialOrd for Node<T>
where
    T: Clone + Ord,
{
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl<T> Ord for Node<T>
where
    T: Clone + Ord,
{
    fn cmp(&self, other: &Self) -> Ordering {
        other
            .f()
            .cmp(&self.f())
            .then_with(|| other.g.cmp(&self.g))
            .then_with(|| self.v.cmp(&other.v))
    }
}

pub struct Edge<T> {
    pub value: T,
    pub dg: usize,
}

pub trait Neighborer<T> {
    fn neighbors(&self, k: &T) -> Vec<Edge<T>>;
}

pub fn search<T>(
    start: &[T],
    goal: &T,
    neighborer: &dyn Neighborer<T>,
    heuristic: fn(a: &T, b: &T) -> usize,
) -> Option<(Vec<T>, usize)>
where
    T: Clone + Ord + Hash,
{
    let mut open = BinaryHeap::new();
    let mut gscore = HashMap::new();
    let mut adjacent = HashMap::<T, T>::new();
    for i in start {
        open.push(Node {
            v: i.clone(),
            g: 0,
            h: heuristic(i, goal),
        });
        gscore.insert(i.clone(), 0);
    }
    while let Some(current) = open.pop() {
        if &current.v == goal {
            let mut revpath = vec![current.v.clone()];
            let mut k = current.v;
            while let Some(i) = adjacent.remove(&k) {
                revpath.push(i.clone());
                k = i;
            }
            revpath.reverse();
            return Some((revpath, current.g));
        }

        for i in neighborer.neighbors(&current.v) {
            let ng = current.g + i.dg;
            if let Some(&g) = gscore.get(&i.value) {
                if ng >= g {
                    continue;
                }
            }
            adjacent.insert(i.value.clone(), current.v.clone());
            gscore.insert(i.value.clone(), ng);
            open.push(Node {
                v: i.value.clone(),
                g: ng,
                h: heuristic(&i.value, goal),
            });
        }
    }
    None
}
