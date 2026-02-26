// btree.rs
use std::fmt;

// B-Tree (CLRS) with minimum degree t.
// - Max keys per node: 2t-1
// - Max children per node: 2t
#[derive(Debug)]
struct BTree {
    t: usize,
    root: Box<Node>,
}

#[derive(Debug)]
struct Node {
    leaf: bool,
    keys: Vec<i32>,
    children: Vec<Box<Node>>,
}

impl Node {
    fn new(leaf: bool, t: usize) -> Self {
        Node {
            leaf,
            keys: Vec::with_capacity(2 * t - 1),
            children: Vec::with_capacity(2 * t),
        }
    }
}

impl BTree {
    fn new(t: usize) -> Self {
        assert!(t >= 2, "t must be >= 2");
        BTree {
            t,
            root: Box::new(Node::new(true, t)),
        }
    }

    fn search(&self, k: i32) -> bool {
        Self::search_node(&self.root, k)
    }

    fn search_node(x: &Node, k: i32) -> bool {
        let mut i = 0usize;
        while i < x.keys.len() && k > x.keys[i] {
            i += 1;
        }
        if i < x.keys.len() && x.keys[i] == k {
            return true;
        }
        if x.leaf {
            return false;
        }
        Self::search_node(&x.children[i], k)
    }

    fn insert(&mut self, k: i32) {
        let t = self.t;
        if self.root.keys.len() == 2 * t - 1 {
            // root full -> grow height
            let mut new_root = Box::new(Node::new(false, t));
            let old_root = std::mem::replace(&mut self.root, Box::new(Node::new(true, t)));
            new_root.children.push(old_root);

            self.split_child(&mut new_root, 0);
            self.root = new_root;

            self.insert_non_full(&mut self.root, k);
        } else {
            self.insert_non_full(&mut self.root, k);
        }
    }

    fn insert_non_full(&self, x: &mut Box<Node>, k: i32) {
        let t = self.t;

        if x.leaf {
            // insert k into sorted keys
            let mut i = x.keys.len();
            x.keys.push(0);
            while i > 0 && k < x.keys[i - 1] {
                x.keys[i] = x.keys[i - 1];
                i -= 1;
            }
            x.keys[i] = k;
            return;
        }

        // find child index
        let mut i = x.keys.len();
        while i > 0 && k < x.keys[i - 1] {
            i -= 1;
        }

        // if child is full, split first
        if x.children[i].keys.len() == 2 * t - 1 {
            self.split_child(x, i);
            if k > x.keys[i] {
                i += 1;
            }
        }

        self.insert_non_full(&mut x.children[i], k);
    }

    // split x.children[i] (full) into two nodes; median key moves up into x
    fn split_child(&self, x: &mut Box<Node>, i: usize) {
        let t = self.t;

        // Borrow the full child y
        let y_leaf;
        let median_key;
        let mut z = Box::new(Node::new(false, t)); // leaf flag set after we read y

        {
            let y = &mut x.children[i];
            y_leaf = y.leaf;
            z.leaf = y_leaf;

            // y.keys: [0 .. 2t-2], median index = t-1
            median_key = y.keys[t - 1];

            // z gets keys [t ..]
            let right_keys = y.keys.split_off(t); // y keeps [0..t-1]
            // y now has keys [0..t-1], we must remove median at t-1
            y.keys.pop(); // remove median
            z.keys = right_keys;

            if !y_leaf {
                // y.children: [0 .. 2t-1], z gets [t ..]
                let right_children = y.children.split_off(t);
                z.children = right_children;
            }
        }

        // Insert new child z after y
        x.children.insert(i + 1, z);
        // Insert median into parent keys at position i
        x.keys.insert(i, median_key);
    }

    fn traverse(&self) -> Vec<i32> {
        let mut out = Vec::new();
        Self::traverse_node(&self.root, &mut out);
        out
    }

    fn traverse_node(x: &Node, out: &mut Vec<i32>) {
        for i in 0..x.keys.len() {
            if !x.leaf {
                Self::traverse_node(&x.children[i], out);
            }
            out.push(x.keys[i]);
        }
        if !x.leaf {
            Self::traverse_node(&x.children[x.keys.len()], out);
        }
    }
}

fn main() {
    let mut bt = BTree::new(3); // t=3 => max keys/node = 5
    let nums = [10, 20, 5, 6, 12, 30, 7, 17, 3, 2, 4, 25, 26, 27, 28];

    for &v in &nums {
        bt.insert(v);
    }

    let sorted = bt.traverse();
    println!("{}", DisplayVec(&sorted));
    println!("Search 12: {}", bt.search(12));
    println!("Search 99: {}", bt.search(99));
}

// Small helper just for printing nicely
struct DisplayVec<'a>(&'a [i32]);
impl<'a> fmt::Display for DisplayVec<'a> {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        for (i, v) in self.0.iter().enumerate() {
            if i > 0 {
                write!(f, " ")?;
            }
            write!(f, "{v}")?;
        }
        Ok(())
    }
}