use std::collections::HashMap;
use std::hash::Hash;

pub struct Edge<T> {
    pub a: T,
    pub b: T,
    pub c: i32,
}

pub struct PairwiseDistances<T> {
    arr: HashMap<T, HashMap<T, i32>>,
}

impl<T> PairwiseDistances<T>
where
    T: Eq + Hash,
{
    fn new() -> Self {
        Self {
            arr: HashMap::new(),
        }
    }

    pub fn edge_cost(&self, a: &T, b: &T) -> Option<i32> {
        self.arr
            .get(a)
            .and_then(|m| m.get(b))
            .and_then(|&k| Some(k))
    }

    fn set(&mut self, a: T, b: T, c: i32) {
        self.arr.entry(a).or_insert(HashMap::new()).insert(b, c);
    }
}

pub fn compute<T>(nodes: &[T], edges: &[Edge<T>]) -> PairwiseDistances<T>
where
    T: Clone + Eq + Hash,
{
    let mut dist = PairwiseDistances::new();
    for i in edges {
        dist.set(i.a.clone(), i.b.clone(), i.c);
    }
    for i in nodes {
        dist.set(i.clone(), i.clone(), 0);
    }
    for k in nodes {
        for i in nodes {
            for j in nodes {
                if let Some(ck) = dist
                    .edge_cost(&i, &k)
                    .and_then(|cik| dist.edge_cost(&k, &j).and_then(|ckj| Some(cik + ckj)))
                {
                    if match dist.edge_cost(&i, &j) {
                        Some(cij) => cij > ck,
                        None => true,
                    } {
                        dist.set(i.clone(), j.clone(), ck);
                    }
                }
            }
        }
    }
    dist
}
