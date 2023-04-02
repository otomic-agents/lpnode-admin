package types

type DeploymentStruct struct {
	ApiVersion string
	Kind       string
	Spec       struct {
		Template struct {
			Spec struct {
				Containers []struct {
					Name string
					Env  []struct {
						Name  string
						Value string
					}
				}
			}
		}
	}
}
