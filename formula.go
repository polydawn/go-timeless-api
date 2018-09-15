package api

type (
	Formula struct {
		Inputs  map[AbsPath]WareID
		Action  FormulaAction
		Outputs map[AbsPath]FormulaOutputSpec
	}

	FormulaOutputSpec struct {
		PackType PackType          `refmt:"packtype"`
		Filter   FilesetPackFilter `refmt:"filter",omitempty`
	}

	// FormulaAction defines the action to perform to "evaluate" the formula --
	// after the input filesets have been assembled, these commands will be run
	// in a contained sandbox on with those filesets,
	// and when the commands terminate, the output filesets will be saved.
	//
	// The definition of the Action includes at minimum what commands to run,
	// but also includes the option of specifying other execution parameters:
	// things like environment variables, working directory, hostname...
	// and (though hopefully you rarely get hung up and need to change these)
	// also things like UID, GID, username, homedir, and soforth.
	// All of these additional parameters have "sensible defaults" if unset.
	//
	// The Action also includes the ability to set "Policy" level -- these
	// define simple privilege levels.  (The default policy is of extremely
	// low privileges.)
	FormulaAction struct {
		// An array of strings to hand as args to exec -- creates a single process.
		Exec []string `refmt:",omitempty"`

		// Noop may be set as an alternative to Exec; this allows manipulations of
		// files that can be done from pure path of inputs and outputs alone.
		Noop bool `refmt:",omitempty"`

		// FUTURE: we want to even more options here as alternatives to exec.
		//  For example, some basic options for manipulating files without full exec.

		// How much power to give the process.  Default is quite low.
		Policy FormulaPolicy `refmt:",omitempty"`

		// The working directory to set when invoking the executable.
		// If not set, will be defaulted to "/task".
		Cwd AbsPath `refmt:",omitempty"`

		// Environment variables.
		Env map[string]string `refmt:",omitempty"`

		// User info -- uid, gid, etc.
		Userinfo *FormulaUserinfo `refmt:",omitempty"`

		// Cradle -- enabled by default, enum value for disable.
		Cradle string `refmt:",omitempty"`

		// Hostname to set inside the container (if the executor supports this -- not all do).
		Hostname string `refmt:",omitempty"`
	}

	FormulaUserinfo struct {
		Uid      *int    `refmt:",omitempty"`
		Gid      *int    `refmt:",omitempty"`
		Username string  `refmt:",omitempty"`
		Homedir  AbsPath `refmt:",omitempty"`
	}
)

/*
	FormulaPolicy constants enumerate the privilege levels a contained process
	can be started with.  (They're a shorthand for linux 'capabilities',
	with some sensible safe sets pre-selected.)

	Policies are meant as a rough, relatively approachable, user-facing
	shorthand for privilege levels.
	In practice they typically map onto linux 'capabilities', but this is
	considered an implementation detail, not guaranteed, and may be executor
	engine specific (for example, the 'chroot' executor cannot provide
	fine-grained capabilities at all).
*/
type FormulaPolicy string

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
	FormulaPolicy_Routine FormulaPolicy = "routine"

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
	FormulaPolicy_Governor FormulaPolicy = "governor"

	/*
		Operate with *ALL CAPABILITIES*.

		This is absolutely not secure against untrusted code -- it is
		completely equivalent in power to root on your host.  Please
		try to use any of the lower power levels first.

		Among the things a system administrator may do is rebooting
		the machine and updating the kernel.  Seriously, *only* use
		with trusted code.
	*/
	FormulaPolicy_Sysad FormulaPolicy = "sysad"
)

type FormulaRunRecord struct {
	Guid      string             `refmt:"guid"`      // random number, presumed globally unique.
	FormulaID FormulaSetupHash   `refmt:"formulaID"` // HID of formula ran.
	Time      int64              `refmt:"time"`      // time at start of build.
	ExitCode  int                `refmt:"exitCode"`  // exit code of the contained process.
	Results   map[AbsPath]WareID `refmt:"results"`   // wares produced by the run!

	Hostname string            `refmt:",omitempty"` // Optional: hostname.  not a trusted field, but useful for debugging.
	Metadata map[string]string `refmt:",omitempty"` // Optional: escape valve.  you can attach freetext here.
}
