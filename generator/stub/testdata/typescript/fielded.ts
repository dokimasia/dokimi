/** Lists the shards of each node. The member onNode that configures node under the default template would have the name of the method onNode. */
//dokimi:stub field={{.Method}}Stub
export interface Cell {
  /** Reports whether the cell has the node n. */
  node(n: string): boolean;
  /** Returns the shards on the node n. */
  onNode(n: string): string[];
}
