## TODO

- [x] **Implement disk caching**

  - Cache indexed file paths to disk (e.g., using JSON or BoltDB)
  - Add option to refresh the cache manually or on a schedule

- [ ] **Optimize search performance**

  - Improve fuzzy matching speed (use a trie or radix tree)
  - Use Go routines to parallelize indexing across directories

- [ ] **Add configuration support**

  - Allow user-defined include/exclude paths
  - Configurable cache location and expiration

- [ ] **UI improvements**

  - Option to run the fuzzy finder in a Wayland/X11 floating window (GUI popup)
  - Display full file paths with horizontal scrolling or truncation strategy
