package api

/*
	`Policy` constants enumerate the privilege levels a contained process
	can be started with.  (They're a shorthand for linux 'capabilities',
	with some sensible safe sets pre-selected.)

	Policies are meant as a rough, relatively approachable, user-facing
	shorthand for privilege levels.
	In practice they typically map onto linux 'capabilities', but this is
	considered an implementation detail, not guaranteed, and may be executor
	engine specific (for example, the 'chroot' executor cannot provide
	fine-grained capabilities at all).
*/
type Policy string

const (
	/*
		Operate with a low privilege, as if you were a regular user on a
		regular system.  No special permissions will be granted
		(and in systems with capabilities support, special permissions
		will not be available even if processes do manage to
		change uid, e.g. through suid binaries; most capabilities
		are dropped).

		This is the safest mode to run as.  And, naturally, the default.

		Note that you may still (separately) set the Userinfo to values like
		uid=0 and gid=0, even while at 'routine' policy privileges.
		This is fine; an executor engine that supports capabilities dropping
		will still result in operations that the "root" user would normally
		be able to perform (like chown any file) will still result in
		permission denied.
	*/
	Policy_Routine Policy = "routine"

	/*
		Operate with escalated but still relatively safe privilege.
		Dangerous capabilities (e.g. "muck with devices") are dropped,
		but the most commonly used of root's powers (like chown any file)
		are available.

		This may be slightly safer than enabling full 'sysad' mode,
		but you should still prefer to use any of the lower power levels
		if possible.

		This mode is the most similar to what you would experience with
		docker defaults.

		This mode should not be assumed secure when combined with host mounts.
		(For example, one can trivially make an executable file in the
		host mount, set it to owner=0, set it setuid, and thus have a
		beachhead ready for a later phase in an attack.)
	*/
	Policy_Governor Policy = "governor"

	/*
		Operate with *ALL CAPABILITIES*.

		This is absolutely not secure against untrusted code -- it is
		completely equivalent in power to root on your host.  Please
		try to use any of the lower power levels first.

		Among the things a system administrator may do is rebooting
		the machine and updating the kernel.  Seriously, *only* use
		with trusted code.
	*/
	Policy_Sysad Policy = "sysad"
)
