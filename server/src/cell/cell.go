package cell

type Cell struct {
	Kind int
	/*
	>0 heal
	=0 noting or seen
	<0 damage
	 */

	Hidden int // =1(hidden) || =0(seen)
}
