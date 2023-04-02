package clustertype

type LpClientPodItem struct {
	Name   string
	Status struct {
		AvailableReplicas int32
	}
}
type LpClientServiceItem struct {
	Name          string
	ConditionsLen int
	Ports         []struct {
		Protocol string
		Port     int
	}
}
