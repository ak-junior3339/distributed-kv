package btree

type Entry struct {
	Key   string
	Value string
}

// The Entry reepresents each Key-Value pair that will be stored in the Node
//{Key:"name",Value:"Rahul"}

type Node struct {
	Leaf     bool    // Whether the Node is A Leaf Node of Not
	Entries  []Entry // it is a Slice (Dynamic sized array) **Sorted** of Key-Value pair's'
	Children []*Node //it is a pointer to the child nodes if any
}

type BTree struct {
	Root  *Node // it holds pointer to the entry point or the root of the tree
	Order int   // Specifies the order / degree of the tree
}

func NewBTree(order int) *BTree { //instantiates a brand new instance in an empty state.
	return &BTree{
		Root: &Node{
			Leaf:    true,
			Entries: make([]Entry, 0),
		},
		Order: order,
	}
}

func (t *BTree) Get(key string) (string, bool) {
	if t.Root == nil {
		return "", false
	}
	return t.Root.search(key)
}

func (n *Node) search(key string) (string, bool) {
	i := 0

	for i < len(n.Entries) && key > n.Entries[i].Key {
		i++
	}
	if i < len(n.Entries) && key == n.Entries[i].Key {
		return n.Entries[i].Value, true
	}
	if n.Leaf == true {
		return "", false
	}
	return n.Children[i].search(key)
}

// splitChild splits the full child node y of parent node x.
// i is the index of y in x.Children.
func (t *BTree) splitChild(x *Node, i int, y *Node) {
	// The number of entries that will move to the new node
	midIdx := t.Order / 2
	medianEntry := y.Entries[midIdx]

	// Create a new node to hold the upper half of y's entries
	z := &Node{
		Leaf:    y.Leaf,
		Entries: append([]Entry(nil), y.Entries[midIdx+1:]...),
	}

	// If y is not a leaf, its upper child pointers must move to z
	if !y.Leaf {
		z.Children = append([]*Node(nil), y.Children[midIdx+1:]...)
		y.Children = y.Children[:midIdx+1]
	}

	// Truncate y to hold only the lower half of its entries
	y.Entries = y.Entries[:midIdx]

	// Insert z into x's child pointers right after y
	x.Children = append(x.Children, nil)
	copy(x.Children[i+2:], x.Children[i+1:])
	x.Children[i+1] = z

	// Insert y's median entry into x's entries
	x.Entries = append(x.Entries, Entry{})
	copy(x.Entries[i+1:], x.Entries[i:])
	x.Entries[i] = medianEntry
}

// Set inserts or updates a key-value pair in the B-Tree.
func (t *BTree) Set(key string, value string) {
	root := t.Root

	// If the root node is completely full, the tree must grow upwards
	if len(root.Entries) == t.Order-1 {
		newRoot := &Node{
			Leaf:     false,
			Children: []*Node{root},
		}
		t.Root = newRoot
		t.splitChild(newRoot, 0, root)
		t.insertNonFull(newRoot, key, value)
	} else {
		t.insertNonFull(root, key, value)
	}
}

// insertNonFull is a recursive helper that finds the right spot in a non-full node.
func (t *BTree) insertNonFull(n *Node, key string, value string) {
	i := len(n.Entries) - 1

	// Case 1: If this is a leaf node, find the sorted position and insert it
	if n.Leaf {
		// Create room by appending an empty entry
		n.Entries = append(n.Entries, Entry{})
		// Shift larger entries to the right
		for i >= 0 && key < n.Entries[i].Key {
			n.Entries[i+1] = n.Entries[i] //classic insertion sort
			i--
		}
		// Place the new entry or update it if the key already exists
		if i >= 0 && n.Entries[i].Key == key {
			n.Entries[i].Value = value               // Update case
			n.Entries = n.Entries[:len(n.Entries)-1] // Remove the unused appended slot
		} else {
			n.Entries[i+1] = Entry{Key: key, Value: value}
		}
		return
	}

	// Case 2: If this is an internal node, find the correct child pathway
	for i >= 0 && key < n.Entries[i].Key {
		i--
	}
	i++

	// Check if the target child node is full

	if len(n.Children[i].Entries) == t.Order-1 {
		t.splitChild(n, i, n.Children[i])
		if i < len(n.Entries) && key > n.Entries[i].Key {
			i++
		}
	}

	// Recurse down
	t.insertNonFull(n.Children[i], key, value)
}
