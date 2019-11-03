# Solution Considerations

Not being too familiar with the Node.JS ecosystem, I wanted to find a solution that delivers good usability while not requiring me to make too many assumptions.

In particular, separating library code from user code is difficult without that background.

Since flamegraphs are visual and interactive, they offer a good user experience for developer who can then draw their own conclusions from the data.

Because of the difficulty encountered in getting complete stack traces from nsolid, I included the option to eliminate any stack frame from within `node_modules` from
the data used to generate the graph.
While this is not 100% accurate, it makes the resulting graphs more useful, provided that the assumption of common ancestry holds.


## Technology Choices

The technologies choices for this proof-of-concept are largely based on familiarity, ease of deployment on Kubernetes as well as time constraints:

* Kubernetes offers an easy way to deploy and expose decoupled applications as used here
* Go is the language I am currently most immediately familiar with
* For Redis there exists an easy to use operator which can provision Sentinel-based failover clusters
* `d3-flame-graph` offers a quick way to visualize the gathered data without needing to provide a query or custom visualisation interface


## Dropped Requirements

Based on the time available, some requirements from the original brief could not be satisfied.

In particular, it is not possible to query the data only for a particular time interval.
This could easily be included by choosing a different data store (e.g. MongoDB) which would support storing the raw traces and their associated time stamps.

While the overall result of the flame graph is not as useful as I had hoped, it is still possible to extract some useful information on the most common blocking calls from it.
