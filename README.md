# Terminationd

Daemon and toolset for working with AWS EC2 Spot Instances.

`terminationd` will continuously poll the EC2 instance metadata and exit when a query for instance's termination time returns a Time-like value matching the regex `.*T.*Z` ([documentation][termination-time-metadata]). HTTP 404 errors, request timeouts, and non-Time-like response contents are assumed to be evidence that the host is _not_ about to be terminated.

## References

- [Spot Instance `termination-time` Metadata documentation][termination-time-metadata]

[termination-time-metadata]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-interruptions.html#termination-time-metadata
