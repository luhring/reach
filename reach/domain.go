package reach

// Domain is what groups Reach resources into high-level categories of kinds of resources.
//
// Examples of possible Domains include: AWS, Azure, GCP, etc. There's also a "generic" domain that houses resource kinds like IPAddress and Hostname, for which there is no further infrastructure information available.
type Domain string
